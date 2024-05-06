package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rivo/tview"
)

var (
	startTime       time.Time
	punchedIn       bool
	linesTyped      int
	totalTimeWorked time.Duration
)

func main() {
	app := tview.NewApplication()
	table := setupTable()

	db, err := InitDB()
	if err != nil {
		panic(err)
	}

	menu := tview.NewList().
		AddItem("Clock In", "Start the work session", '1', func() { punchIn(db) }).
		AddItem("Clock Out", "End the work session", '2', func() { punchOut(db) }).
		AddItem("History", "View the work history", '3', func() { showHistory(db, table) }).
		AddItem("Quit", "Exit the application", 'q', func() {
			app.Stop()
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(menu, 7, 1, true).
		AddItem(table, 0, 3, true)

	updateTicker(app, table)

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			insertUserData(db, time.Now(), totalTimeWorked, linesTyped)
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func setupTable() *tview.Table {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(false, false)

	labels := []string{"Current Time", "Duty Status", "Time Worked", "Lines Typed"}
	for i, label := range labels {
		table.SetCell(i, 0, tview.NewTableCell(label).
			SetAlign(tview.AlignRight).
			SetSelectable(false))
		table.SetCell(i, 1, tview.NewTableCell("").
			SetAlign(tview.AlignLeft))
	}

	return table
}

func updateTicker(app *tview.Application, table *tview.Table) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			app.QueueUpdateDraw(func() {
				updateTable(table)
			})
		}
	}()
}

func updateTable(table *tview.Table) {
	currentTime := time.Now().Format("15:04:05")
	dutyStatus := "Off Duty"
	timeWorked := totalTimeWorked.String()
	if punchedIn {
		dutyStatus = "On Duty"
		timeWorked = (totalTimeWorked + time.Since(startTime)).Round(time.Second).String()
	}

	table.SetCell(0, 1, tview.NewTableCell(currentTime))
	table.SetCell(1, 1, tview.NewTableCell(dutyStatus))
	table.SetCell(2, 1, tview.NewTableCell(timeWorked))
	table.SetCell(3, 1, tview.NewTableCell(string(rune(linesTyped))))
}

func insertUserData(db *sql.DB, date time.Time, totalTimeWorked time.Duration, totalLinesTyped int) error {
	_, err := db.Exec(`
		INSERT INTO user_data (date, total_time_worked, total_lines_typed)
		VALUES (?, ?, ?)
	`, date, totalTimeWorked.String(), totalLinesTyped)
	return err
}

func punchOut(db *sql.DB) {
	if punchedIn {
		totalTimeWorked += time.Since(startTime)
		punchedIn = false
	}
	insertUserData(db, time.Now(), totalTimeWorked, linesTyped)
}

func punchIn(db *sql.DB) {
	if !punchedIn {
		startTime = time.Now()
		punchedIn = true
	}
	insertUserData(db, startTime, totalTimeWorked, linesTyped)
}

func showHistory(db *sql.DB, table *tview.Table) {
	rows, err := db.Query("SELECT * FROM user_data")
	if err != nil {
		return
	}
	defer rows.Close()

	for i := 0; i < 10; i++ {
		table.SetCell(i, 0, tview.NewTableCell(""))
		table.SetCell(i, 1, tview.NewTableCell(""))
	}

	rowIndex := 0
	for rows.Next() {
		var date time.Time
		var totalTimeWorked string
		var totalLinesTyped int
		err := rows.Scan(&date, &totalTimeWorked, &totalLinesTyped)
		if err != nil {
			// handle error i guess?
			continue
		}
		table.SetCell(rowIndex, 0, tview.NewTableCell(date.Format("2006-01-02")))
		table.SetCell(rowIndex, 1, tview.NewTableCell(fmt.Sprintf("Time Worked: %s, Lines Typed: %d", totalTimeWorked, totalLinesTyped)))
		rowIndex++
	}

	if err := rows.Err(); err != nil {
		// handle error too
		return
	}
}
