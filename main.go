package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	cssfile := "./build/css/min.css"
	jsfile := "./build/js/min.js"
	err, files := Getfiles("./css", ".css")
	check(err)
	Concat(files, "./build/css/min.css")
	err, jsfiles := Getfiles("./js", ".js")
	check(err)
	Concat(jsfiles, "./build/js/min.js")
	Minifycss(cssfile)
	Minifyjs(jsfile)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Minifycss(file string) {

	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	b, err := os.ReadFile(file)
	check(err)

	mb, err := m.Bytes("text/css", b)
	check(err)

	err = os.WriteFile(file, mb, 0644)
	check(err)

	fmt.Println(string(mb))

}

func Minifyjs(file string) {
	m := minify.New()

	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	b, err := os.ReadFile(file)
	check(err)

	mb, err := m.Bytes("application/javascript", b)
	check(err)

	err = os.WriteFile(file, mb, 0644)
	check(err)

	fmt.Println(string(mb))
}

func Concat(files []string, fileout string) {

	var buf bytes.Buffer
	for _, file := range files {
		b, err := os.ReadFile(file)
		check(err)

		buf.Write(b)
	}

	err := os.WriteFile(fileout, buf.Bytes(), 0644)
	check(err)

}

func Getfiles(in, ext string) (error, []string) {
	var filelist []string
	files, err := os.ReadDir(in)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		if strings.Contains(file.Name(), ext) {
			filelist = append(filelist, in+"/"+file.Name())
		}
	}
	return nil, filelist
}
