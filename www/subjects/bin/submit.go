package main

import (
	"fmt"
	"strconv"
	"strings"
)

func submit() {

	// CONFIG answer: what subjects
	var what, subjects string

	var qid, skip int
	var err error
	qid, err = strconv.Atoi(strings.TrimSpace(gReq.FormValue("qid")))
	if xx(err) {
		return
	}
	if strings.TrimSpace(gReq.FormValue("skip")) != "" {
		skip = 1
	} else {

		// BEGIN CONFIG answer: what, subjects

		what = strings.TrimSpace(gReq.FormValue("what"))
		subjects = strings.TrimSpace(gReq.FormValue("subjects"))

		// TODO: tests

		/*
			if animal == "" {
				x(fmt.Errorf("Missing answer for animal"))
				return
			}
		*/
		// END CONFIG answer: what subjects

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

	// CONFIG answer: what subjects
	// NOTE: The number of question marks must match the number of fields and arguments
	_, err = tx.Exec("INSERT INTO answers(qid, uid, skip, what, subjects) VALUES (?, ?, ?, ?, ?);",
		qid,
		gUserID,
		skip,
		// CONFIG answer: what subjects
		what,
		subjects)
	if xx(err) {
		tx.Rollback()
		return

	}

	tx.Commit()

	question()

}
