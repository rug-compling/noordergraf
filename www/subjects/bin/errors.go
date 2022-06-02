package main

import (
	"bytes"
	"fmt"
	"html"
	"runtime"
)

// Print plain error
func x(err error, msg ...interface{}) bool {
	return xerr(false, err, msg...)
}

// Print error with line number
func xx(err error, msg ...interface{}) bool {
	return xerr(true, err, msg...)
}

func xerr(withLineno bool, err error, msg ...interface{}) bool {
	if err == nil {
		return false
	}

	var b bytes.Buffer
	_, filename, lineno, ok := runtime.Caller(2)
	if withLineno && ok {
		b.WriteString(fmt.Sprintf("%v:%v: %v", filename, lineno, err))
	} else {
		b.WriteString(err.Error())
	}
	if len(msg) > 0 {
		b.WriteString(",")
		for _, m := range msg {
			b.WriteString(fmt.Sprintf(" %v", m))
		}
	}
	headers()
	fmt.Printf(`<html>
<head>
<title>Error</title>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" href="../favicon.ico" type="image/ico">
<link rel="stylesheet" href="../style.css" type="text/css">
</head>
<body class="error">
<div id="header">
&nbsp;
</div>
<div id="main">
<h1>Error</h2>
%s
</div>
<div id="footer">
<a href=".">Continue</a>
</body>
</html>
`, html.EscapeString(b.String()))

	return true

}
