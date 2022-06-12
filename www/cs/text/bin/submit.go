package main

import (
	"fmt"
	"strconv"
	"strings"
)

func submit() {

	// CONFIG answer: atext
	var atext string

	var qid, skip int
	var err error
	qid, err = strconv.Atoi(strings.TrimSpace(gReq.FormValue("qid")))
	if xx(err) {
		return
	}
	if strings.TrimSpace(gReq.FormValue("skip")) != "" {
		skip = 1
	} else {

		// BEGIN CONFIG answer: atext

		atext = strings.TrimSpace(gReq.FormValue("atext"))

		if atext == "" {
			x(fmt.Errorf("Geen antwoord"))
			return
		}

		// END CONFIG answer: atext

	}

	tx, err := gDB.Begin()
	if xx(err) {
		return
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM answers WHERE qid = %d AND uid = %d", qid, gUserID))
	if xx(err) {
		tx.Rollback()
		return

	}

	// CONFIG answer: atext
	// NOTE: The number of question marks must match the number of fields and arguments
	_, err = tx.Exec("INSERT INTO answers(qid, uid, skip, atext) VALUES (?, ?, ?, ?);",
		qid,
		gUserID,
		skip,
		// CONFIG answer: atext
		atext)
	if xx(err) {
		tx.Rollback()
		return

	}

	tx.Commit()

	question()

}
