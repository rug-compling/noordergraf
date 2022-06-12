package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

func unskip() {

	total, skipped, done, ok := getDone()
	if !ok {
		return
	}

	reset := strings.TrimSpace(gReq.FormValue("unskip"))
	if reset == "yes" {
		_, err := gDB.Exec(fmt.Sprintf("DELETE FROM answers WHERE uid = %d AND skip = 1", gUserID))
		if xx(err) {
			return
		}
		question()
		return
	}

	if reset == "no" {
		b, err := ioutil.ReadFile("../templates/finished.html")
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

	question()
}
