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
	// Check if database exists
	//

	_, err := os.Stat("data.sqlite")
	if err != nil {
		log.Fatal("Database does not exists")
	}

	////////////////////////////////////////////////////////////////
	//
	// Open database from table: questions
	//

	db, err := sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Read from table: questions
	//

	qids := make(map[string]int)
	rows, err := db.Query("SELECT qid, label FROM questions")
	if err != nil {
		log.Fatal(err)
	}

	var qid int
	var label string
	for rows.Next() {
		err = rows.Scan(&qid, &label)
		if err != nil {
			log.Fatal(err)
		}
		qids[label] = qid
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////
	//
	// Update data into table: questions
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
		qid, ok := qids[record[0]]
		if ok {
			// CONFIG question: image qtext
			// NOTE: number of question marks must match number of arguments
			_, err = tx.Exec(`UPDATE questions
                              SET image = ?,
                                  qtext = ?
                              WHERE qid = ?`,
				record[1], // image
				record[2], // qtext
				qid)
		} else {
			// CONFIG question: image qtext
			// NOTE: number of question marks must match number of fields and arguments
			_, err = tx.Exec("INSERT INTO questions(label, image, qtext) VALUES (?, ?, ?)",
				record[0], // label (required)
				record[1], // image
				record[2]) // qtext
		}
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
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
