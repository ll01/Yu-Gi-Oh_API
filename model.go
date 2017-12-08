package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
)

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
	Attributes     []string       `json:"attributes"`
	EffectKeyWords []string       `json:"effectkeywords"`
	LinkArrows     []string       `json:"linkarrows"`
	Archtypes      []string       `json:"archtypes"`
}

func (currentCard *card) getCardFromID(cardDatabase *sql.DB, cardIDToSearch int) error {
	err := setMainCardData("id", strconv.Itoa(cardIDToSearch), cardDatabase, currentCard)
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
func setAuxiliaryData(tableName, columnName string, currentCardData []string, cardDatabase *sql.DB) {
	rows, err := cardDatabase.Query("SELECT " + columnName + " FROM " + tableName +
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

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
