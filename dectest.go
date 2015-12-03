package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/shopspring/decimal"
	"github.com/faide/decimal"
	"log"
)

func main() {
	var (
		Db     *sqlx.DB
		id1    int             = 1
		id2    int             = 2
		amnt1  decimal.Decimal = decimal.NewFromFloat(55.33)
		amnt2  decimal.Decimal = decimal.New(0, 0)
		schema                 = `
			CREATE TABLE IF NOT EXISTS some_record
				(
				id integer,
				amount numeric
				);`
	)

	Db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln(err)
	}
	err = Db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// create the schema
	tx := Db.MustBegin()
	tx.MustExec(schema)
	tx.Commit()

	// set a value
	tx = Db.MustBegin()
	tx.Exec("INSERT INTO some_record (id, amount) VALUES ($1, $2)", id1, amnt1)
	tx.Exec("INSERT INTO some_record (id, amount) VALUES ($1, $2)", id2, amnt2)
	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	// use previously set value to query and test our scanner
	amount := decimal.Decimal{}
	err = Db.QueryRow(
		"SELECT amount FROM some_record WHERE id==?1", id1).Scan(&amount)

	switch err {

	case sql.ErrNoRows:
		log.Println("No row found...")

	case nil:
		log.Println("Found:", amount)

	default:
		log.Println("Error occured:", err)

	}

	err = Db.QueryRow(
		"SELECT amount FROM some_record WHERE id==?1", id2).Scan(&amount)

	switch err {

	case sql.ErrNoRows:
		log.Println("No row found...")

	case nil:
		log.Println("Found:", amount)

	default:
		log.Println("Error occured:", err)
	}
}
