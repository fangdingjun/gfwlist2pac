gfwlist2pac
===========

gfwlist2pac implement by go.

Install
=======

    go get github.com/fangdingjun/gfwlist2pac

Usage
=====

    Usage: gfwlist2pac [OPTION]...

    -p, --proxy      special proxy address, ex: 127.0.0.1:8080
    -t, --proxy_type special proxy type, ex: HTTP,SOCKS5, HTTPS, etc
    -o, --out        save the output to the file, default stdout
    -u, --url        the autoproxy gfwlist url, defaullt: https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt

    -h, --help       show help and exit
    -v, --version    show version and exit

License
=======

   Copyright (C) 2014-2015 Dingjun Fang

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
