package main

import (
	"github.com/dchest/authcookie"

	"fmt"
	"time"
)

func headers() {
	fmt.Print(`Content-type: text/html; charset=utf-8
Cache-Control: no-cache
Pragma: no-cache
`)

	if gUserAuth {
		exp := time.Now().AddDate(0, 0, 14).UTC()
		au := authcookie.New(gUserSec+"|"+gUserHash, exp, []byte(cSecret))
		fmt.Printf("Set-Cookie: %s-auth=%s; Expires=%s\n", cCookiePrefix, au, exp.Format(time.RFC1123))
	}

	fmt.Println()
}
