package main

type questionType struct {

	// All fields must start with an uppercase letter

	// CONFIG question: Image Persons
	Image   string
	Persons string

	Qid     int
	Done    int
	Skipped int
	Todo    int
}

const (
	// BEGIN CONFIG: configurable global variables

	// URL of the application, including trailing slash
	cBaseUrl = "https://noordergraf.rug.nl/cs/subjects/"

	// For sending mail to the participant: name and mail address of sender
	cMailName = "Noordergraf Subjects Survey"
	cMailFrom = "p.c.j.kleiweg@rug.nl"

	// Username and password for the smtp server
	// Leave these empty if the server doesn't need log in
	cSmtpUser = ""
	cSmtpPass = ""

	// Smtp server, including port number
	cSmtpServ = "smtp.rug.nl:25"

	// To prevent name clashes with other cookies on the same site
	cCookiePrefix = "ngcs"

	// Used for encryption
	cSecret = "change this bla bla boo bah"

	// END CONFIG: configurable global variables
)
