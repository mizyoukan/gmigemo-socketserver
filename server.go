package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/koron/gomigemo/embedict"
	"github.com/koron/gomigemo/migemo"
)

var (
	host      = flag.String("h", "127.0.0.1", "Server host")
	port      = flag.Int("p", 13730, "Server port")
	vim       = flag.Bool("v", false, "Use vim style regexp")
	emacs     = flag.Bool("e", false, "Use emacs style regexp")
	noNewline = flag.Bool("n", false, "Don't use newline match")
)

func adjustMatcher(m migemo.Matcher) {
	o := m.GetOptions()
	o.OpWSpaces = ""
	if *vim {
		o.OpOr = "\\|"
		o.OpGroupIn = "\\%("
		o.OpGroupOut = "\\)"
		if *noNewline {
			o.OpWSpaces = "\\_s*"
		}
	} else if *emacs {
		o.OpOr = "\\|"
		o.OpGroupIn = "\\("
		o.OpGroupOut = "\\)"
		if *noNewline {
			o.OpWSpaces = "\\s-*"
		}
	} else if *noNewline {
		o.OpWSpaces = "\\s*"
	}
	m.SetOptions(o)
}

func query(d migemo.Dict, s string) (string, error) {
	m, err := d.Matcher(s)
	if err != nil {
		return "", nil
	}
	adjustMatcher(m)
	return m.Pattern()
}

func proc(c *net.TCPConn, d migemo.Dict) (err error) {
	if err = c.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		err = fmt.Errorf("set deadline error: %v", err)
		return
	}

	b := make([]byte, 1024)
	n, err := c.Read(b)
	if err != nil {
		err = fmt.Errorf("read request error: %v", err)
		return
	}
	b = b[:n]

	var p string
	if p, err = query(d, string(b)); err != nil {
		err = fmt.Errorf("migemo query error: %v", err)
		return
	}

	if _, err = c.Write([]byte(p)); err != nil {
		err = fmt.Errorf("write response error: %v", err)
		return
	}

	return
}

func main() {
	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Resolve TCP address error:", err)
		os.Exit(1)
	}

	d, err := embedict.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Migemo dictionary load error:", err)
		os.Exit(1)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Listen tcp error:", err)
		os.Exit(1)
	}
	defer func() {
		if err := l.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "Listener close error:", err)
		}
	}()

	for {
		if c, err := l.AcceptTCP(); err != nil {
			fmt.Fprintln(os.Stderr, "Accept tcp error:", err)
			os.Exit(1)
		} else {
			go func() {
				defer func() {
					if err := c.Close(); err != nil {
						fmt.Fprintf(os.Stderr, "[%v] Connection close error: %v\n",
							c.RemoteAddr(), err)
					}
				}()
				if err := proc(c, d); err != nil {
					fmt.Fprintf(os.Stderr, "[%v] %v\n", c.RemoteAddr(), err)
				}
			}()
		}
	}
}
