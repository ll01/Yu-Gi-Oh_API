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

type card struct {
	ID             int         `json:"id"`
	Passcode       int         `json:"passcode"`
	NameEN         string      `json:"name_en"`
	NameJP         null.String `json:"name_jp"`
	NameFR         null.String `json:"name_fr"`
	NameDE         null.String `json:"name_de"`
	NameIT         null.String `json:"name_it"`
	NameKR         null.String `json:"name_kr"`
	NamePT         null.String `json:"name_pt"`
	NameES         null.String `json:"name_es"`
	Cardtype       null.String `json:"card_type"`
	Attribute      null.String `json:"attribute"`
	LevelOrRank    null.Int    `json:"level/rank/link"`
	Scale          null.Int    `json:"scale"`
	Attack         null.Int    `json:"attack"`
	Defence        null.Int    `json:"defence"`
	Material       null.String `json:"material"`
	Attributes     []string    `json:"attributes"`
	EffectKeyWords []string    `json:"effectkeywords"`
	LinkArrows     []string    `json:"linkarrows"`
	Archtypes      []string    `json:"archtypes"`
}

func (currentCard *card) getCardFromID(cardDatabase *sql.DB, cardIDToSearch int) {
	setMainCardData("id", strconv.Itoa(cardIDToSearch), cardDatabase, currentCard)
	BuildCard(currentCard, cardDatabase)

}

func (currentCard *card) getCardFromPasscode(cardDatabase *sql.DB, cardPasscodeToSearch int) {
	setMainCardData("passcode", strconv.Itoa(cardPasscodeToSearch), cardDatabase, currentCard)
	BuildCard(currentCard, cardDatabase)
}

func (currentCard *card) getCardFromName(cardDatabase *sql.DB, cardName, contryCode string) {
	passcode, err := GetPasscodeFromName(cardName, contryCode, cardDatabase)
	if err == nil {
		currentCard.getCardFromPasscode(cardDatabase, passcode)
	}
}

func setMainCardData(columnName, dataToSearchFor string, cardDatabase *sql.DB, currentCard *card) {
	rows, err := cardDatabase.Query("SELECT * FROM main_card_data WHERE main_card_data." + columnName + " = " + dataToSearchFor)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&currentCard.ID, &currentCard.Passcode, &currentCard.NameEN,
			&currentCard.NameJP, &currentCard.Cardtype, &currentCard.Attribute, &currentCard.LevelOrRank,
			&currentCard.Scale, &currentCard.Attack, &currentCard.Defence, &currentCard.Material)
		checkErr(err)
	}
	HandleCardSearchError(err)
}
func (currentCard *card) setAuxiliaryData(tableName string, currentCardData []string, cardDatabase *sql.DB) []string {
	rows, err := cardDatabase.Query("SELECT name FROM " + tableName +
		" LEFT JOIN main_card_data ON " + tableName + ".passcode=main_card_data.passcode WHERE main_card_data.passcode = " +
		strconv.Itoa(currentCard.Passcode))
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var temp sql.NullString
		err = rows.Scan(&temp)
		if temp.Valid {
			currentCardData = append(currentCardData, temp.String)
		}
	}
	checkErr(err)
	return currentCardData
}

func (currentCard *card) setGlobalCardNames(cardDatabase *sql.DB) {
	rows, err := cardDatabase.Query("SELECT name, contry_code FROM foreign_name_table LEFT JOIN main_card_data ON " +
		"foreign_name_table.passcode=main_card_data.passcode WHERE  main_card_data.passcode = " + strconv.Itoa(currentCard.Passcode))
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var name null.String
		var contryCode string
		err = rows.Scan(&name, &contryCode)
		checkErr(err)
		switch contryCode {
		case "FR":
			currentCard.NameFR = name
		case "DE":
			currentCard.NameDE = name
		case "IT":
			currentCard.NameIT = name
		case "KR":
			currentCard.NameKR = name
		case "PT":
			currentCard.NamePT = name
		case "ES":
			currentCard.NameES = name

		}
	}
}

func GetPasscodeFromName(cardName, contryCode string, cardDatabase *sql.DB) (int, error) {
	var passcode = 0
	var query string
	cardName = "\"" + cardName + "\""
	switch contryCode {
	case "EN":
		query = "SELECT passcode FROM main_card_data WHERE name_EN = " + cardName
	case "JP":
		query = "SELECT passcode FROM main_card_data WHERE name_JP = " + cardName
	default:
		query = "SELECT passcode FROM foreign_name_table WHERE name = " + cardName
	}
	rows, err := cardDatabase.Query(query)
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&passcode)
	}
	// what if jp is null

	return passcode, err
}

func BuildCard(currentCard *card, cardDatabase *sql.DB) {
	currentCard.Archtypes = currentCard.setAuxiliaryData(GetTableNameInstance().Archtype(), currentCard.Archtypes, cardDatabase)
	currentCard.LinkArrows = currentCard.setAuxiliaryData(GetTableNameInstance().LinkArrow(), currentCard.LinkArrows, cardDatabase)
	currentCard.EffectKeyWords = currentCard.setAuxiliaryData(GetTableNameInstance().EffectKeyword(), currentCard.EffectKeyWords, cardDatabase)
	currentCard.Attributes = currentCard.setAuxiliaryData(GetTableNameInstance().Attribute(), currentCard.Attributes, cardDatabase)
	currentCard.setGlobalCardNames(cardDatabase)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
