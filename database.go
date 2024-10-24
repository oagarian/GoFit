package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func openDatabase(dataSourceName string) (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}

func createTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			basal_metabolism REAL,
			week INTEGER
		);
		CREATE TABLE IF NOT EXISTS calories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			day TEXT,
			calories_consumed REAL,
			week INTEGER
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func isBasalMetabolismSet(currentWeek int) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE week = ?`, currentWeek).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func insertBasalMetabolism(metabolism float64, currentWeek int) {
	_, err := db.Exec(`INSERT INTO users (basal_metabolism, week) VALUES (?, ?)`, metabolism, currentWeek)
	if err != nil {
		log.Fatal(err)
	}
}

func addCalories(dayOfWeek string, calories float64, currentWeek int) {
	_, err := db.Exec(`INSERT INTO calories (day, calories_consumed, week) VALUES (?, ?, ?)`, dayOfWeek, calories, currentWeek)
	if err != nil {
		log.Fatal(err)
	}
}

func getWeeklyCalories(currentWeek int) (map[string]float64, float64) {
	rows, err := db.Query(`SELECT day, SUM(calories_consumed) FROM calories WHERE week = ? GROUP BY day`, currentWeek)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	weeklyCalories := map[string]float64{}
	for rows.Next() {
		var day string
		var calories float64
		if err := rows.Scan(&day, &calories); err != nil {
			log.Fatal(err)
		}
		weeklyCalories[day] = calories
	}

	var basalMetabolism float64
	err = db.QueryRow(`SELECT basal_metabolism FROM users WHERE week = ?`, currentWeek).Scan(&basalMetabolism)
	if err != nil {
		log.Fatal(err)
	}

	return weeklyCalories, basalMetabolism
}
