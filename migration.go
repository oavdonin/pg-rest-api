package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// CsvToDB loads the data from provided titanic.csv to the database
func CsvToDB(filename string, db *sql.DB) {
	// Open the file
	recordFile, err := os.Open(filename)
	if err != nil {
		log.Println("An error encountered ::", err)
		return
	}

	// Setup the reader
	reader := csv.NewReader(recordFile)

	// Read the records
	allRecords, err := reader.ReadAll()
	if err != nil {
		log.Println("An error encountered ::", err)
		return
	}

	// Check if people table is empty
	var numOfPeopleInDb uint64
	log.Println("Checking if database contains data inside \"people\" table whether force ingestion has been requested...")
	db.QueryRow("select count(*) from people;").Scan(&numOfPeopleInDb)
	repopulateDb, err := strconv.ParseBool(getEnv("POPULATE_DB", "false"))
	if err != nil {
		log.Fatal(err)
	}
	if numOfPeopleInDb > 5 && repopulateDb != true {
		log.Println("Table \"people\" exists in the DB. Populate flag is not set. Continue w/o data ingestion.")
	} else {
		log.Println("Continue with data ingestion")
		// Create people table if it doesn't exist
		_, err = db.Exec(`
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
			CREATE TABLE IF NOT EXISTS people (
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
			TRUNCATE TABLE people;
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
	}
	// Closing file and DB
	err = recordFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// getEnv implements a fallback value for standard env function
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
