/*
   go gfwlist2pac

   Copyright (C) 2014-2015 Dingjun Fang<fangdingjun@gmail.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
Convert autoproxy gfwlist to PAC file.

Usage: gfwlist2pac [OPTION]...

    -p, --proxy      special proxy address, ex: 127.0.0.1:8080
    -t, --proxy_type special proxy type, ex: HTTP,SOCKS5, HTTPS, etc
    -o, --out        save the output to the file, default stdout
    -u, --url        the autoproxy gfwlist url, defaullt: https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt

    -h, --help       show help and exit
    -v, --version    show version and exit
*/
package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	flag "github.com/ogier/pflag"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	proxy             = flag.StringP("proxy", "p", "127.0.0.1:8080", "")
	proxy_type        = flag.StringP("proxy_type", "t", "HTTP", "")
	output            = flag.StringP("out", "o", "", "")
	autoproxylist_url = flag.StringP("url", "u", "", "")
	show_version      = flag.BoolP("version", "v", false, "")
)

const (
	Help = `Usage: gfwlist2pac [OPTION]...
    -p, --proxy      special proxy address, ex: 127.0.0.1:8080
    -t, --proxy_type special proxy type, ex: HTTP,SOCKS5, HTTPS, etc
    -o, --out        save the output to the file, default stdout
    -u, --url        the autoproxy gfwlist url, defaullt: https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt

    -h, --help       show help and exit
    -v, --version    show version and exit
`

	Version     = "gfwlist2pac 0.0.1"
	gfwlist_url = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
)

type args struct {
	/* proxy address */
	Proxy string

	/* gfwlist domains name */
	Domains map[string]int

	/* custom domains name */
	Custom map[string]int
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s", Help)
		os.Exit(1)
	}
	flag.Parse()

	if *show_version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	uri := gfwlist_url
	if *autoproxylist_url != "" {
		uri = *autoproxylist_url
	}

	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	decoder := base64.NewDecoder(base64.StdEncoding, res.Body)

	var domains map[string]int = map[string]int{}

	reader := bufio.NewReader(decoder)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			break
		}

		s := parse(string(line))
		if s != "" {
			if _, ok := domains[s]; !ok {
				domains[s] = 1
			}
		}
	}

	res.Body.Close()

	/* template arguments */
	aa := args{
		fmt.Sprintf("%s %s", *proxy_type, *proxy),
		domains,
		read_custom_list("custom.txt"),
	}

	/* read template file */
	t, err := template.ParseFiles("pac.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	out := os.Stdout
	if *output != "" {
		/* custom output */
		out, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	/* output */
	err = t.Execute(out, aa)
	if err != nil {
		log.Fatal(err)
	}
}
