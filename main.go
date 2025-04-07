package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "./mydata.db"

func main() {
	fmt.Println("Starting go-Data app....")

	err := os.Remove(dbFile)
	if err != nil {
		return
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)

		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)

	}
	fmt.Println("Successfully connected to database.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		quantity INTEGER
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table 'items' created or already exists.")

	insertSQL := "INSERT INTO items(name, quantity) VALUES(?, ?)"
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("Error preparing insert statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	itemsToInsert := []struct {
		Name     string
		Quantity int
	}{
		{"Apple", 5},
		{"Banana", 10},
		{"Orange", 3},
	}

	fmt.Println("Inserting data...")
	for _, item := range itemsToInsert {
		_, err := stmt.Exec(item.Name, item.Quantity)
		if err != nil {
			log.Printf("Error inserting item %s: %v\n", item.Name, err)
		} else {
			fmt.Printf("Inserted: %s, Quantity: %d\n", item.Name, item.Quantity)
		}
	}

	fmt.Println("Querying data...")
	querySQL := "SELECT id, name, quantity FROM items"
	rows, err := db.Query(querySQL)
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows:", err)

		}
	}(rows)

	fmt.Println("Items in database:")
	for rows.Next() {
		var id int
		var name string
		var quantity int

		err := rows.Scan(&id, &name, &quantity)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		fmt.Printf("ID: %d, Name: %s, Quantity: %d\n", id, name, quantity)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	fmt.Println("App finished successfully.")

}
