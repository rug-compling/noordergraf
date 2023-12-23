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

	// What algorithm to use for presenting questions to participants
	// 1 = random order
	// 2 = questions with the least answers so far first
	// 3 = questions with the least IDENTICAL answers so far first
	// Method 1 is suitable for small datasets, where it is expected that
	// all participants will answer all questions.
	// Methods 2 and 3 are suitable for large datasets. Participants may
	// answer only a small subset of questions. You want to make sure the
	// overall coverage is as large as possible.
	// Method 3 selects questions where participants disagree first, to get
	// as many answers for questions where the right answer is less certain.
	// Method 3 is not suited for questions where differences are expected,
	// such as when the particpants are asked to enter a longer piece of text,
	// where variation in interpunction, capitalisation, newlines, choice of
	// words, etc is expected.
	cAlgo = 3

	// END CONFIG: configurable global variables
)
