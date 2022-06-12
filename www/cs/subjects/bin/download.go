package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func download() {
	fmt.Print("Content-type: text/plain; charset=utf-8\n\n")

	// CONFIG answer: what subjects
	rows, err := gDB.Query(`
			SELECT label,
				what,
				subjects
			FROM results
			WHERE skip = 0
			ORDER by qid,
				what,
				subjects
			;`)

	if err != nil {
		fmt.Println(err)
		return
	}

	w := csv.NewWriter(os.Stdout)

	seen := false
	for rows.Next() {
		seen = true

		// CONFIG answer: what subjects
		var what, subjects string

		var label string

		// CONFIG answer: what subjects
		err := rows.Scan(&label, &what, &subjects)
		if err != nil {
			rows.Close()
			fmt.Println(err)
			break
		}

		err = w.Write([]string{
			label,
			// CONFIG answer: what subjects
			what,
			subjects,
			// Yes, that comma after the last parameter is needed
		})
		if err != nil {
			rows.Close()
			fmt.Println(err)
			break

		}
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println(err)
	}

	if !seen {
		fmt.Println("Nothing submitted yet")
	}
}
