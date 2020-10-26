package migrations

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// CsvToDB loads the data from provided titanic.csv to the database
func CsvToDB(filename string, databaseURL string) {
	// Open the file
	recordFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	// Read the records
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	//Connect to the DB
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Check if people table is empty
	_, tableCheck := db.Query("select * from people;")
	log.Println("Creating table people...")
	if tableCheck == nil {
		log.Println("People table is presented in the DB")
		os.Exit(0)
	}

	// Create people table if it doesn't exist
	_, err = db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE people (
			uuid uuid DEFAULT uuid_generate_v4 (),
			survived VARCHAR NOT NULL,
			pclass SMALLSERIAL NOT NULL,
			name VARCHAR NOT NULL,
			sex VARCHAR,
			age FLOAT,
			siblingsOrSpousesAboard SMALLSERIAL,
			parentsOrChildrenAboard SMALLSERIAL,
			fare FLOAT,
			PRIMARY KEY (uuid)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Loading all records
	for i := 1; i < len(allRecords); i++ {
		_, err = db.Exec(
			"INSERT INTO people (survived, pclass, name, sex, age, siblingsOrSpousesAboard, parentsOrChildrenAboard, fare) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
			allRecords[i][0], allRecords[i][1], allRecords[i][2], allRecords[i][3],
			allRecords[i][4], allRecords[i][5], allRecords[i][6], allRecords[i][7])
		if err != nil {
			log.Fatal(err)
		}
	}

	// Closing file and DB

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = recordFile.Close()
	if err != nil {
		log.Fatal(err)
	}

}
