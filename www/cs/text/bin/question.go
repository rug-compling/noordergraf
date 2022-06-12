package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func question() {

	total, skipped, done, ok := getDone()
	if !ok {
		return
	}

	if done+skipped >= total {
		var src string
		if skipped == 0 {
			src = "finished.html"
		} else {
			src = "skipped.html"
		}
		b, err := ioutil.ReadFile("../templates/" + src)
		if xx(err) {
			return
		}
		t, err := template.New("foo").Parse(string(b))
		if xx(err) {
			return
		}
		headers()
		xx(t.Execute(os.Stdout, questionType{
			Done:    done,
			Skipped: skipped,
			Todo:    total - done - skipped,
		}))
		return
	}

	// CONFIG question: image qtext
	rows, err := gDB.Query(fmt.Sprintf(`
			SELECT qid,
				image,
                qtext
			FROM questions
			WHERE qid NOT IN ( SELECT qid FROM answers WHERE uid = %d )
			ORDER BY RANDOM()
			LIMIT 1`,
		gUserID))
	if xx(err) {
		return
	}
	var qid int
	// CONFIG question: image qtext
	var image, qtext string
	for rows.Next() {
		// CONFIG question: item image qtext
		if xx(rows.Scan(&qid, &image, &qtext)) {
			return
		}
		ok = true
	}
	if !ok {
		xx(fmt.Errorf("Missing row"))
	}

	b, err := ioutil.ReadFile("../templates/question.html")
	if xx(err) {
		return
	}

	t, err := template.New("foo").Parse(string(b))
	if xx(err) {
		return
	}
	headers()
	xx(t.Execute(os.Stdout, questionType{
		Done:    done,
		Skipped: skipped,
		Todo:    total - done - skipped,
		Qid:     qid,

		// CONFIG question: image qtext
		Image: image,
		Qtext: qtext,
	}))
}

func getDone() (total, skipped, done int, ok bool) {
	if !gUserAuth {
		xx(fmt.Errorf("Not logged in"))
		return
	}
	rows, err := gDB.Query(fmt.Sprintf(`
			SELECT * FROM
			( SELECT COUNT(*) FROM questions),
			( SELECT COUNT(*) FROM answers WHERE uid = %d AND skip > 0),
			( SELECT COUNT(*) FROM answers WHERE uid = %d AND skip = 0)`,
		gUserID, gUserID))
	if xx(err) {
		return
	}
	for rows.Next() {
		if xx(rows.Scan(&total, &skipped, &done)) {
			return
		}
		ok = true
	}
	if !ok {
		xx(fmt.Errorf("Missing row"))
		return
	}
	return
}
