package main

import (
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"fmt"
	"net/http/cgi"
	"runtime/debug"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Print("Content-type: text/plain\n\n")
			fmt.Println(r)
			fmt.Println(string(debug.Stack()))
		}
	}()

	var err error
	gReq, err = cgi.Request()
	if xx(err) {
		return
	}
	if xx(gReq.ParseForm()) {
		return
	}

	gDB, err = sql.Open("sqlite3", "../db/data.sqlite")
	if xx(err) {
		return
	}
	defer gDB.Close()

	if ok := getLogin(); !ok {
		return
	}

	action := gReq.FormValue("action")
	if gReq.Method == "GET" {
		if action == "dl" {
			download()
		} else if action == "login" {
			login()
		} else if gUserAuth {
			question()
		} else {
			loginForm()
		}
	} else if gReq.Method == "POST" {
		if action == "login" {
			loginRequest()
		} else if !gUserAuth {
			loginForm()
		} else if action == "submit" {
			submit()
		} else if action == "unskip" {
			unskip()
		} else {
			x(fmt.Errorf("Missing or invalid POST action %q", action))
		}
	} else {
		x(fmt.Errorf("Method not allowed: %s", gReq.Method))
	}
}
