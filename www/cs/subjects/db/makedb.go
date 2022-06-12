package main

import (
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
)

func main() {

	////////////////////////////////////////////////////////////////
	//
	// Open database
	//

	db, err := sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range []string{
		"PRAGMA encoding = 'UTF-8';",
		"PRAGMA foreign_keys = 1;",
	} {
		_, err = db.Exec(p)
		if err != nil {
			log.Fatal(err)
		}
	}

	////////////////////////////////////////////////////////////////
	//
	// Create table: questions
	//

	// CONFIG question: image persons names
	_, err = db.Exec(`CREATE TABLE questions (
                        qid     INTEGER PRIMARY KEY AUTOINCREMENT,
                        label   TEXT UNIQUE,
                        image   TEXT,
                        persons TEXT
                      );`)
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Read data into table: questions
	//

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fp, err := os.Open("questions.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(fp)
	r.Comment = '#'
	r.TrimLeadingSpace = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// CONFIG question: image persons
		// NOTE: number of question marks must match number of fields and arguments
		_, err = tx.Exec("INSERT INTO questions(label, image, persons) VALUES (?, ?, ?);",
			record[0], // label (required)
			record[1], // image
			record[2]) // persons
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Create table: users
	//

	_, err = db.Exec(`CREATE TABLE users (
                        uid     INTEGER PRIMARY KEY AUTOINCREMENT,
                        email   TEXT NOT NULL UNIQUE,
                        sec     TEXT,
                        pw      TEXT,
                        expires TEXT
                      );`)
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Create table: answers
	//

	// CONFIG answer: what subjects
	_, err = db.Exec(`CREATE TABLE answers (
                        aid      INTEGER PRIMARY KEY AUTOINCREMENT,
                        qid      INTEGER,
                        uid      INTEGER,
                        skip     INTEGER DEFAULT 0,
                        what     TEXT DEFAULT "",
                        subjects TEXT DEFAULT "",
                      FOREIGN KEY(qid) REFERENCES questions(qid),
                      FOREIGN KEY(uid) REFERENCES users(uid)
                      );`)
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Create a view for easy browsing of answers
	//

	_, err = db.Exec(`CREATE VIEW results AS
        SELECT * FROM answers LEFT JOIN questions USING(qid);`)
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Make some indexes
	//

	for _, cmd := range []string{
		"CREATE INDEX auid ON answers(uid);",
		"CREATE INDEX askip ON answers(skip);",
		"CREATE UNIQUE INDEX uemail ON users(email);"} {
		_, err = db.Exec(cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	////////////////////////////////////////////////////////////////
	//
	// Close database
	//

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

}
