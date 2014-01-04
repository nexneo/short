// short is a utility to shorten urls using google url shortener
//
// INSTALL
//
// go get github.com/nexneo/short
//
// UPDATE
//
// go get -u github.com/nexneo/short
//
// USAGE
//
// `short http://github.com` [print short url]
//
// `short -c http://github.com` [print and copy short url to clipboard (Mac OS X only)]
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nexneo/easyreq"
)

type request struct {
	Url string `json:"longUrl"`
}

type response struct {
	ShortUrl string `json:"id"`
}

func main() {
	req := easyreq.NewJson(&request{longUrl})
	data, err := req.Do("POST", "https://www.googleapis.com/urlshortener/v1/url")
	die(err)
	res := &response{}
	die(json.NewDecoder(data.Body).Decode(res))
	data.Body.Close()

	ret := res.ShortUrl
	fmt.Println(ret)

	if *doCopy {
		cmd := exec.Command("pbcopy")
		in, err := cmd.StdinPipe()
		die(err)
		die(cmd.Start())

		_, err = fmt.Fprint(in, ret)
		die(err)

		in.Close()
		die(cmd.Wait())

		log.Println("[Copied to clipboard]")
	}
}

var doCopy = flag.Bool("c", false, "Copy to clipboard")
var longUrl string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n\n", "short [-c] URL")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *doCopy && len(os.Args) == 3 {
		longUrl = os.Args[2]
	} else if len(os.Args) == 2 {
		longUrl = os.Args[1]
	}

	if longUrl == "" {
		flag.Usage()
		os.Exit(0)
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
