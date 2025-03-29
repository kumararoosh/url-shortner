package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connect_to_db() *sql.DB {
	connStr := "host=localhost port=5432 user=shorty password=secret dbname=shortener sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	fmt.Println("Connected to Postgres!")

	return db
}

func make_table_with_schema() {
	db := connect_to_db()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS shortened_urls (
		original_url VARCHAR(255),
		shortened_url VARCHAR(255)
	)`)

	if err != nil {
		log.Fatalf("Cannot create table in psql: %v", err)
	}

	db.Close()

}

func insert_row_into_table(original_url string, new_url string) {
	db := connect_to_db()

	sqlStatement := `
	INSERT INTO shortened_urls (original_url, shortened_url) VALUES ($1, $2)
	`

	_, err := db.Exec(sqlStatement, original_url, new_url)

	if err != nil {
		log.Fatalf("cannot insert to psql: %v", err)
	}
	
	db.Close();
}