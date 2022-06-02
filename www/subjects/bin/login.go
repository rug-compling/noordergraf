package main

import (
	"github.com/dchest/authcookie"

	"crypto/sha256"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	expireFormat = "2006-01-02T15:04:05"
)

func loginForm() {
	b, err := ioutil.ReadFile("../templates/login.html")
	if xx(err) {
		return
	}
	headers()
	fmt.Print(string(b))
}

func loginRequest() {
	email := strings.TrimSpace(gReq.FormValue("email"))

	if email == "" {
		x(fmt.Errorf("Missing e-mail address"))
		return
	}

	if !strings.Contains(email, "@") {
		x(fmt.Errorf("Invalid e-mail address"))
		return
	}

	tx, err := gDB.Begin()
	if xx(err) {
		return
	}

	ehash := getHash(email)

	auth := rand16()
	sec := rand16()
	expires := time.Now().Add(time.Hour).Format(expireFormat)

	result, err := gDB.Exec(fmt.Sprintf("UPDATE `users` SET `sec` = %q, `pw` = %q, `expires` = %q WHERE `email` = %q",
		sec, auth, expires, ehash))
	if xx(err) {
		tx.Rollback()
		return
	}
	n, err := result.RowsAffected()
	if xx(err) {
		tx.Rollback()
		return
	}
	if n == 0 {
		_, err = tx.Exec("INSERT INTO `users` (email, sec, pw, expires) VALUES (?,?,?,?);",
			ehash, sec, auth, expires)
		if xx(err) {
			tx.Rollback()
			return
		}
	}

	if xx(tx.Commit()) {
		return
	}

	err = sendmail(
		email,
		"Log in",
		fmt.Sprintf(
			"Go to this URL to log in: %sbin/?action=login&pw=%s\n\nNOTE: This URL will be valid for one hour\n",
			cBaseUrl, url.QueryEscape(auth)))
	if xx(err) {
		return
	}

	b, err := ioutil.ReadFile("../templates/loginmailed.html")
	if xx(err) {
		return
	}

	t, err := template.New("foo").Parse(string(b))
	if xx(err) {
		return
	}
	headers()
	if xx(t.Execute(os.Stdout, email)) {
		return
	}
}

func login() {

	if ok := expire(); !ok {
		return
	}

	pw := strings.TrimSpace(gReq.FormValue("pw"))
	if pw == "" {
		pw = "none" // Or someone without a password could log in
	}

	rows, err := gDB.Query(fmt.Sprintf("SELECT `email`,`uid`,`sec` FROM `users` WHERE `pw` = %q", pw))
	if xx(err) {
		return
	}
	if rows.Next() {
		err := rows.Scan(&gUserHash, &gUserID, &gUserSec)
		rows.Close()
		if xx(err) {
			return
		}
		_, err = gDB.Exec(fmt.Sprintf("UPDATE `users` SET `pw` = '', `expires` = '9999' WHERE `email` = %q", gUserHash))
		if xx(err) {
			return
		}
		gUserAuth = true
		headers()
		fmt.Printf(`<head>
<meta http-equiv="Refresh" content="0; URL=%sbin/">
</head>
`,
			cBaseUrl)
	} else {
		x(fmt.Errorf("Log in failed"))
	}
}

func expire() (ok bool) {

	// Remove expired log-in requests

	now := time.Now().Format(expireFormat)
	tx, err := gDB.Begin()
	if xx(err) {
		return
	}

	// Only remove users that haven't submitted any data yet
	_, err = gDB.Exec(fmt.Sprintf("DELETE FROM users WHERE `expires` < %q AND `uid` NOT IN ( SELECT `uid` FROM answers )", now))
	if xx(err) {
		tx.Rollback()
		return
	}

	// These have already submitted data, don't remove user, only clear log-in tokens
	_, err = gDB.Exec(fmt.Sprintf("UPDATE `users` SET `sec` = \"\", `pw` = \"\", `expires` = \"\" WHERE `expires` < %q", now))
	if xx(err) {
		tx.Rollback()
		return
	}
	if xx(tx.Commit()) {
		return
	}

	ok = true
	return
}

func getLogin() (ok bool) {
	auth, err := gReq.Cookie(cCookiePrefix + "-auth")
	if err != nil {
		ok = true
		return
	}
	s := strings.SplitN(authcookie.Login(auth.Value, []byte(cSecret)), "|", 2)
	if len(s) == 2 {
		gUserSec = s[0]
		gUserHash = s[1]
		rows, err := gDB.Query(fmt.Sprintf("SELECT `uid`,`sec` FROM `users` WHERE `email` = %q", gUserHash))
		if xx(err) {
			return
		}
		for rows.Next() {
			var sec string
			if xx(rows.Scan(&gUserID, &sec)) {
				return
			}
			if sec == gUserSec {
				gUserAuth = true // Yes, this user has logged in
			}
		}
	}
	ok = true
	return
}

func rand16() string {
	a := make([]byte, 16)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		a[i] = byte(97 + rnd.Intn(26))
	}
	return string(a)
}

func getHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(s)))
}

func sendmail(to, subject, body string) (err error) {
	msg := fmt.Sprintf(`From: %q <%s>
To: %s
Subject: %s
Content-type: text/plain; charset=UTF-8

%s
`, cMailName, cMailFrom, to, subject, body)

	if cSmtpUser != "" {
		auth := smtp.PlainAuth("", cSmtpUser, cSmtpPass, strings.Split(cSmtpServ, ":")[0])
		err = smtp.SendMail(cSmtpServ, auth, cMailFrom, []string{to}, []byte(msg))
	} else {
		err = smtp.SendMail(cSmtpServ, nil, cMailFrom, []string{to}, []byte(msg))
	}
	return
}
