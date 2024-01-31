// convert between png and jpg
// no gap under here due to godoc.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
)

// ToPng converts an an image from png to jpeg or from jpeg to png
// change function name to indicate what it actually does
// filepath.basename validate stupid user input
// go doc filepath to strip extension
func ToPng(
	scannedfilename string, imageBytes []byte,
) ([]byte, error) {
	// do super cool string split stuff
	var splitit []string
	var fronthalf string

	splitit = strings.Split(scannedfilename, ".")
	fronthalf = splitit[0]
	// some awesome person implemented an algorithm to check content
	// types in the http package
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/png":
		jpgcombo := fronthalf + ".jpg"
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("unable to decode png: %w", err)
		}
		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return nil, fmt.Errorf("unable to encode jpeg: %w", err)
		}
		os.WriteFile(jpgcombo, buf.Bytes(), 0644)
	case "image/jpeg":
		// declare type when defining variables so you know what
		// type it is.
		// define err/e up top
		pngcombo := fronthalf + ".png"
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("unable to decode jpeg: %w", err)
		}

		buf := &bytes.Buffer{}
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("unable to encode png: %w", err)
		}
		// look into os.create to create and defer closing
		// and finally f(file).write
		os.WriteFile(pngcombo, buf.Bytes(), 0644)
		fmt.Print("[+] Conversion complete\n")
	}

	return nil, fmt.Errorf(
		"unable to convert %#v to png",
		contentType,
	)
}

func main() {
	var bytearraytocnv []byte
	var scannedfilename string
	var e error
	var filenamearg string
	// put flag methods in init function
	flag.StringVar(
		&filenamearg,
		"f",
		"",
		"Specify filename to png or jpeg",
	)
	flag.Parse()
	//  look at flag.Narg
	if filenamearg != "" {
		// if len(os.Args) > 2 {
		fmt.Printf(
			"[+] Found filename specified in argument: %s\n",
			filenamearg,
		)
		bytearraytocnv, e = os.ReadFile(filenamearg)
		if e != nil {
			fmt.Println(e.Error())
			os.Exit(10)
		}
		ToPng(filenamearg, bytearraytocnv)
		os.Exit(0)
	}
	fmt.Print("Type complete path of png or jpg\n")
	_, e = fmt.Scanln(&scannedfilename)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(9)
	}
	bytearraytocnv, e = os.ReadFile(scannedfilename)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(10)
	}
	ToPng(scannedfilename, bytearraytocnv)
}
