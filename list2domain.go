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

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

/*
parse autoproxy gfwlist line,
return domain name
*/
func parse(line string) string {

	/* remove space */
	line = strings.Trim(line, " ")

	if line == "" {
		return ""
	}

	/* ignore ip address */
	if net.ParseIP(line) != nil {
		return ""
	}

	/* ignore pattern */
	if strings.Index(line, ".") == -1 {
		return ""
	}

	/* ignore comment, whitelist, regex */
	if line[0] == '[' ||
		line[0] == '!' ||
		line[0] == '/' ||
		line[0] == '@' {
		return ""
	}

	return gethostname(line)
}

func gethostname(line string) string {
	c := line[0]
	ss := line

	/* replace '*' */
	if strings.Index(ss, "/") == -1 {
		if strings.Index(ss, "*") != -1 && ss[:2] != "||" {
			ss = strings.Replace(ss, "*", "/", -1)
		}
	}

	switch c {
	case '.':
		ss = fmt.Sprintf("http://%s", ss[1:])
	case '|':
		switch ss[1] {
		case '|':
			ss = fmt.Sprintf("http://%s", ss[2:])
		default:
			ss = ss[1:]
		}
	default:
		if strings.HasPrefix(ss, "http") {
			ss = ss
		} else {
			ss = fmt.Sprintf("http://%s", ss)
		}
	}

	/* process */
	u, err := url.Parse(ss)
	if err != nil {
		log.Printf("%s: %s\n", line, err)
		return ""
	}
	host := u.Host
	if n := strings.Index(host, "*"); n != -1 {
		for i := n; i < len(host); i++ {
			if host[i] == '.' {
				host = host[i:]
				break
			}
		}
	}
	return strings.TrimLeft(host, ".")
}

/*
read the custom domain list,
one domain per line
*/
func read_custom_list(fn string) map[string]int {
	if fn == "" {
		return map[string]int{}
	}

	fp, err := os.Open(fn)
	if err != nil {
		log.Println(err)
		return map[string]int{}
	}

	defer fp.Close()

	domains := map[string]int{}

	reader := bufio.NewReader(fp)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Printf("%s\n", err)
			}
			break
		}

		ll := string(line)

		if ll != "" {
			if _, ok := domains[ll]; !ok {
				domains[ll] = 1
			}
		}
	}

	return domains
}
