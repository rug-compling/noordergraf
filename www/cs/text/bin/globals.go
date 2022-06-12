package main

import (
	"database/sql"
	"net/http"
)

var (
	gUserAuth bool
	gUserHash string
	gUserSec  string
	gUserID   int

	gDB *sql.DB

	gReq *http.Request
)
