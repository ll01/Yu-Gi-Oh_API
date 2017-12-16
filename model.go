package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/guregu/null"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
)

var tables = []string{"archtype_table", "attribute_table", "effect_keyword_table", "foreign_name_table", "link_arrow_table"}

type CardNames struct {
	NameFR null.String `json:"name_fr"`
	NameDE null.String `json:"name_de"`
	NameIT null.String `json:"name_it"`
	NameKR null.String `json:"name_kr"`
	NamePT null.String `json:"name_pt"`
	NameES null.String `json:"name_es"`
}
type card struct {
	ID              int         `json:"id"`
	Passcode        int         `json:"passcode"`
	NameEN          string      `json:"name_en"`
	NameJP          null.String `json:"name_jp"`
	Cardtype        null.String `json:"card_type"`
	Attribute       null.String `json:"attribute"`
	LevelOrRank     null.Int    `json:"level/rank/link"`
	Scale           null.Int    `json:"scale"`
	Attack          null.Int    `json:"attack"`
	Defence         null.Int    `json:"defence"`
	Material        null.String `json:"material"`
	Attributes      []string    `json:"attributes"`
	EffectKeyWords  []string    `json:"effectkeywords"`
	LinkArrows      []string    `json:"linkarrows"`
	Archtypes       []string    `json:"archtypes"`
	globalCardNames CardNames   `json:"cardnames"`
}

func (currentCard *card) getCardFromID(cardDatabase *sql.DB, cardIDToSearch int) error {
	err := setMainCardData("id", strconv.Itoa(cardIDToSearch), cardDatabase, currentCard)

	currentCard.Archtypes = currentCard.setAuxiliaryData(GetTableNameInstance().Archtype(), currentCard.Archtypes, cardDatabase)
	currentCard.LinkArrows = currentCard.setAuxiliaryData(GetTableNameInstance().LinkArrow(), currentCard.LinkArrows, cardDatabase)
	currentCard.EffectKeyWords = currentCard.setAuxiliaryData(GetTableNameInstance().EffectKeyword(), currentCard.EffectKeyWords, cardDatabase)
	currentCard.Attributes = currentCard.setAuxiliaryData(GetTableNameInstance().Attribute(), currentCard.Attributes, cardDatabase)

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
func (currentCard *card) setAuxiliaryData(tableName string, currentCardData []string, cardDatabase *sql.DB) []string {
	rows, err := cardDatabase.Query("SELECT name FROM " + tableName +
		" LEFT JOIN main_card_data ON " + tableName + ".passcode=main_card_data.passcode WHERE main_card_data.passcode = " +
		strconv.Itoa(currentCard.Passcode))
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
	return currentCardData
}

func (currentCard *card) setNames(currentNameData []CardNames, cardDatabase *sql.DB) {
	rows, err := cardDatabase.Query("SELECT name, contry_code FROM foreign_name_table LEFT JOIN main_card_data ON" +
		"foreign_name_table.passcode=main_card_data.passcode WHERE  main_card_data.passcode = " + strconv.Itoa(currentCard.Passcode))

	checkErr(err)
	if rows.Next() {
		var name null.String
		var contryCode string
		err = rows.Scan(&name, &contryCode)
		currentCard.SetCardNames(name, contryCode)
	}
}

func (currentCard *card) SetCardNames(name null.String, contryCode string) {
	switch contryCode {
	case "FR":
		currentCard.globalCardNames.NameFR = name
		break
	case "DE":
		currentCard.globalCardNames.NameDE = name
		break
	case "IT":
		currentCard.globalCardNames.NameIT = name
		break
	case "KR":
		currentCard.globalCardNames.NameKR = name
		break
	case "PT":
		currentCard.globalCardNames.NamePT = name
		break
	case "ES":
		currentCard.globalCardNames.NameES = name
		break

	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
