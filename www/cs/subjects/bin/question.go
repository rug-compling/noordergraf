package main

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	/*
		Bij tweede en derde algoritme:
		COUNT(*) geeft minstens 1, ook als er geen antwoord is
		(vanwege de LEFT JOIN).	In het laatste geval is uid=NULL.
		Door uid als tweede sorteerveld toe te voegen krijg
		je eerst de vragen met 0 antwoorden, en dan die met 1.
	*/
	algo = []string{
		// CONFIG question: image name tagline
		`SELECT qid,
			image,
			persons
		FROM questions
		WHERE qid NOT IN ( SELECT qid FROM answers WHERE uid = %d )
		ORDER BY RANDOM()
		LIMIT 1`,
		// CONFIG question: image name tagline
		`SELECT *
		FROM (
			SELECT qid,
				image,
				persons
			FROM questions
			OUTER LEFT JOIN answers
			USING (qid)
			WHERE skip IS NULL OR SKIP = 0
			GROUP BY (qid)
			ORDER BY COUNT(*), uid
		)
		WHERE qid NOT IN ( SELECT qid FROM answers WHERE uid = %d )
		LIMIT 1`,
		// CONFIG question: image name tagline (2 times)
		// CONFIG answer: animal colour, size
		`SELECT qid,
			image,
			persons
		FROM (
			SELECT qid,
				image,
				persons,
				COUNT(*) as c,
				uid
			FROM questions
			OUTER LEFT JOIN answers
			USING (qid)
			WHERE skip IS NULL OR SKIP = 0
			GROUP BY qid, what, subjects
		)
		WHERE qid NOT IN ( SELECT qid FROM answers WHERE uid = %d )
		GROUP BY qid
		ORDER BY MAX(c), uid
		LIMIT 1`,
	}
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

	rows, err := gDB.Query(fmt.Sprintf(algo[cAlgo-1], gUserID))
	if xx(err) {
		return
	}
	var qid int
	// CONFIG question: image persons
	var image, persons string
	for rows.Next() {
		// CONFIG question: item image persons
		if xx(rows.Scan(&qid, &image, &persons)) {
			return
		}
		ok = true
	}
	if !ok {
		xx(fmt.Errorf("Missing row"))
	}

	var plist string
	if strings.TrimSpace(persons) == "" {
		plist = "<input type=\"hidden\" name=\"subjects\" value=\"--NONE--\">\n"
	} else {
		var buf bytes.Buffer
		buf.WriteString("<p>\nWie is/zijn hier begraven?<br> (Of, bij gedenkteken: wie wordt hier herdacht?)</p>\n<p>\n")

		pp := strings.Split(persons, "|")
		for i, p := range pp {
			aa := strings.SplitN(p, "#", 2)
			id := html.EscapeString(aa[0])
			name := html.EscapeString(aa[1])
			checked := ""
			if i == 0 && len(pp) == 1 {
				checked = " checked"
			}
			fmt.Fprintf(&buf, `
<input type="checkbox" id="p%d" name="subjects" value="%s"%s>
<label for="p%d">%s</label><br>
`, i, id, checked, i, name)
		}

		buf.WriteString("</p>\n<hr>\n")
		plist = buf.String()
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

		// CONFIG question: image persons
		Image:   html.EscapeString(image),
		Persons: plist,
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
