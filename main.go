package main

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"regexp"
	"strings"

	compression "github.com/nurlantulemisov/imagecompression"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	//load files
	cssfile := "./build/css/min.css"
	jsfile := "./build/js/min.js"

	//get files into []string
	err, files := Getfiles("./css", ".css")
	check(err)
	err, jsfiles := Getfiles("./js", ".js")
	check(err)
	err, imgfiles := GetImageFiles("./img")
	check(err)

	//concatenate files
	Concat(files, "./build/css/min.css")
	Concat(jsfiles, "./build/js/min.js")

	//minify files
	Minifycss(cssfile)
	Minifyjs(jsfile)

	//optimize images
	for _, str := range imgfiles {
		go func() {
			Optimizer(str)
		}()

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
	check(err)
	for _, file := range files {

		if strings.Contains(file.Name(), ext) {
			filelist = append(filelist, in+"/"+file.Name())
		}
	}
	return nil, filelist
}
func GetImageFiles(in string) (error, []string) {
	var filelist []string
	files, err := os.ReadDir(in)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		if strings.Contains(file.Name(), ".png") || strings.Contains(file.Name(), ".jpg") || strings.Contains(file.Name(), ".webp") {
			filelist = append(filelist, in+"/"+file.Name())
		}
	}
	return nil, filelist
}

func Optimizer(origin string) {
	file, err := os.Open(origin)
	check(err)

	img, err := png.Decode(file)
	check(err)

	fi, err := os.Stat(origin)
	check(err)

	compressing, _ := compression.New(50)
	compressingImage := compressing.Compress(img)

	f, err := os.Create(origin)
	check(err)

	defer func(f *os.File) {
		err := f.Close()
		check(err)
	}(f)

	err = png.Encode(f, compressingImage)
	check(err)
	fif, err := os.Stat(origin)
	check(err)

	fmt.Println(origin, " before: ", fi.Size(), "after: ", fif.Size())
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
