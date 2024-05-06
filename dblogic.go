package main

import (
	"database/sql"
	"time"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./work.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_data (
			date DATE,
			total_time_worked TEXT,
			total_lines_typed INT
		);
		CREATE TABLE IF NOT EXISTS file_data (
			filename TEXT,
			file_type TEXT,
			directory_path TEXT,
			date_made DATE,
			lines_wrote INT,
			time_spent_in_file TEXT
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertUserData(db *sql.DB, date time.Time, totalTimeWorked time.Duration, totalLinesTyped int) error {
	_, err := db.Exec(`
        INSERT INTO user_data (date, total_time_worked, total_lines_typed)
        VALUES (?, ?, ?)
    `, date, totalTimeWorked.String(), totalLinesTyped)
	return err
}

func InsertFileData(db *sql.DB, filename, fileType, directoryPath string, dateMade time.Time, linesWrote int, timeSpentInFile time.Duration) error {
	_, err := db.Exec(`
        INSERT INTO file_data (filename, file_type, directory_path, date_made, lines_wrote, time_spent_in_file)
        VALUES (?, ?, ?, ?, ?, ?)
    `, filename, fileType, directoryPath, dateMade, linesWrote, timeSpentInFile.String())
	return err
}
