package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type card struct {
	ID          int            `json:"id"`
	Passcode    int            `json:"passcode"`
	NameEN      string         `json:"name_en"`
	NameJP      sql.NullString `json:"name_jp"`
	Cardtype    sql.NullString `json:"card_type"`
	Attribute   sql.NullString `json:"attribute"`
	LevelOrRank sql.NullInt64  `json:"level/rank/link"`
	Scale       sql.NullInt64  `json:"scale"`
	Attack      sql.NullInt64  `json:"attack"`
	Defence     sql.NullInt64  `json:"defence"`
	Material    sql.NullString `json:"material"`
}

func (current_card *card) getCardFromID(cardDatabase *sql.DB, cardIDToSearch int) error {
	rows, err := cardDatabase.Query("SELECT * FROM ygo_main WHERE  ygo_main.id = " + strconv.Itoa(cardIDToSearch))
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&current_card.ID, &current_card.Passcode, &current_card.NameEN,
			&current_card.NameJP, &current_card.Cardtype, &current_card.Attribute, &current_card.LevelOrRank,
			&current_card.Scale, &current_card.Attack, &current_card.Defence, &current_card.Material)
		checkErr(err)
	}
	return err
}

func (current_card *card) getCardFromPasscode(cardDatabase *sql.DB, cardPasscodeToSearch int) {
	rows, err := cardDatabase.Query("SELECT * FROM ygo_main WHERE  ygo_main.passcode = " + strconv.Itoa(cardPasscodeToSearch))
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&current_card.ID, &current_card.Passcode, &current_card.NameEN,
			&current_card.NameJP, &current_card.Cardtype, &current_card.Attribute, &current_card.LevelOrRank,
			&current_card.Scale, &current_card.Attack, &current_card.Defence, &current_card.Material)
		checkErr(err)

	}
}

func (current_card *card) getCardFromNameEN(cardDatabase *sql.DB, cardNameEnToSearch string) {
	rows, err := cardDatabase.Query("SELECT * FROM ygo_main WHERE  ygo_main.name_en = " + cardNameEnToSearch)
	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&current_card.ID, &current_card.Passcode, &current_card.NameEN,
			&current_card.NameJP, &current_card.Cardtype, &current_card.Attribute, &current_card.LevelOrRank,
			&current_card.Scale, &current_card.Attack, &current_card.Defence, &current_card.Material)
		checkErr(err)

	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
