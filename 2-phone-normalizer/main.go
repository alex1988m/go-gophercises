package main

import (
	// "database/sql"
	// "strconv"

	"regexp"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func dbInit() (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Error("error while opening database:", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error("error while pinging database:", err)
		return nil, err
	}
	return db, err
}

func prepareTable() error {
	var err error
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS phone (id INTEGER PRIMARY KEY, phone TEXT)")
	if err != nil {
		log.Error("error while creating table:", err)
	}

	stmt, err := db.Prepare("INSERT INTO phone(phone) VALUES(?)")
	if err != nil {
		log.Error("error while preparing statement:", err)
	}
	defer stmt.Close()
	phones := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, p := range phones {
		_, err = stmt.Exec(p)
		if err != nil {
			log.Error("error while inserting phone:", err)
		}
	}
	return nil
}

func main() {
	var err error
	db, err = dbInit()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// err = prepareTable()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	log.Debug("opening sqlite database")

	tx, err := db.Begin()
	if err != nil {
		log.Error("error while starting transaction:", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Error("error while rolling back transaction:", err)
		} else {
			tx.Commit()
		}
	}()

	rows, err := tx.Query("SELECT * FROM phone")
	if err != nil {
		log.Error("error while querying database:", err)
	}
	defer rows.Close()

	normStmt, err := tx.Prepare("UPDATE phone SET phone = ? WHERE id = ?")
	if err != nil {
		log.Error("error while preparing statement:", err)
	}
	defer normStmt.Close()

	for rows.Next() {
		var id int
		var phone string
		err = rows.Scan(&id, &phone)
		if err != nil {
			log.Error("error while scanning rows:", err)
		}
		_, err := normStmt.Exec(normalizePhone(phone), id)
		if err != nil {
			log.Error("error while updating phone:", err)
		}
		// fmt.Printf("id: %d, phone: %s\n", id, normalizePhone(phone))
	}

	if err = rows.Err(); err != nil {
		log.Error("error while iterating rows:", err)
	}

	log.Debug("closing sqlite database")
}

func normalizePhone(phone string) string {
	// Compile a regular expression to match non-digit characters
	re := regexp.MustCompile(`\D`)

	// Replace all non-digit characters with an empty string
	normalized := re.ReplaceAllString(phone, "")

	return normalized
}
