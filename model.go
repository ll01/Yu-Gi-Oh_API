package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
)

var tables = []string{"archtype_table", "attribute_table", "effect_keyword_table", "foreign_name_table", "link_arrow_table"}

type CardNames struct {
	NameFR sql.NullString `json:"name_fr"`
	NameDE sql.NullString `json:"name_de"`
	NameIT sql.NullString `json:"name_it"`
	NameKR sql.NullString `json:"name_kr"`
	NamePT sql.NullString `json:"name_pt"`
	NameES sql.NullString `json:"name_es"`
}
type card struct {
	ID             int            `json:"id"`
	Passcode       int            `json:"passcode"`
	NameEN         string         `json:"name_en"`
	NameJP         sql.NullString `json:"name_jp"`
	Cardtype       sql.NullString `json:"card_type"`
	Attribute      sql.NullString `json:"attribute"`
	LevelOrRank    sql.NullInt64  `json:"level/rank/link"`
	Scale          sql.NullInt64  `json:"scale"`
	Attack         sql.NullInt64  `json:"attack"`
	Defence        sql.NullInt64  `json:"defence"`
	Material       sql.NullString `json:"material"`
	CardNames      []string       `json:"cardnames"`
	Attributes     []string       `json:"attributes"`
	EffectKeyWords []string       `json:"effectkeywords"`
	LinkArrows     []string       `json:"linkarrows"`
	Archtypes      []string       `json:"archtypes"`
}

func (currentCard *card) getCardFromID(cardDatabase *sql.DB, cardIDToSearch int) error {
	err := setMainCardData("id", strconv.Itoa(cardIDToSearch), cardDatabase, currentCard)

	setAuxiliaryData(GetTableNameInstance().Archtype(), currentCard.Archtypes, cardDatabase)
	setAuxiliaryData(GetTableNameInstance().LinkArrow(), currentCard.LinkArrows, cardDatabase)
	setAuxiliaryData(GetTableNameInstance().EffectKeyword(), currentCard.EffectKeyWords, cardDatabase)
	setAuxiliaryData(GetTableNameInstance().Attribute(), currentCard.Attributes, cardDatabase)

	return err
}

func (currentCard *card) getCardFromPasscode(cardDatabase *sql.DB, cardPasscodeToSearch int) error {
	err := setMainCardData("passcode", strconv.Itoa(cardPasscodeToSearch), cardDatabase, currentCard)
	return err
}

func setMainCardData(columnName, dataToSearchFor string, cardDatabase *sql.DB, currentCard *card) error {
	rows, err := cardDatabase.Query("SELECT * FROM main_card_data WHERE main_card_data." + columnName + " = " + dataToSearchFor)
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&currentCard.ID, &currentCard.Passcode, &currentCard.NameEN,
			&currentCard.NameJP, &currentCard.Cardtype, &currentCard.Attribute, &currentCard.LevelOrRank,
			&currentCard.Scale, &currentCard.Attack, &currentCard.Defence, &currentCard.Material)
		checkErr(err)
	}
	return err
}
func setAuxiliaryData(tableName string, currentCardData []string, cardDatabase *sql.DB) {
	rows, err := cardDatabase.Query("SELECT name FROM " + tableName +
		" LEFT JOIN main_card_data ON " + tableName + ".passcode=main_card_data.passcode")
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		var temp sql.NullString
		err = rows.Scan(&temp)
		if temp.Valid {
			currentCardData = append(currentCardData, temp.String)
		}
	}
	checkErr(err)
}

func (currentCard *card) setNames(currentNameData []CardNames, cardDatabase *sql.DB) {
	rows, err := cardDatabase.Query("SELECT name, contry_code FROM foreign_name_table LEFT JOIN main_card_data " +
		"ON foreign_name_table.passcode=main_card_data.passcode WHERE " + strconv.Itoa(currentCard.Passcode) +
		"=main_card_data.passcode")
	checkErr(err)
	if rows.Next() {
		var name sql.NullString
		var contryCode string
		err = rows.Scan(&name, &contryCode)
	}
}

func dfapfaos(name, contryCode string) {
	switch contryCode {
	case "FR":

	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
