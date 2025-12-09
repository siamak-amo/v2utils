// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"bufio"
	"golang.org/x/term"
	flag "github.com/spf13/pflag"
)

const (
	CMD_CONVERT int = iota
	CMD_TEST
	CMD_RUN
)

type Opt struct {
	Cmd int
	URL *string
	IN_File *string

	scanner *bufio.Scanner
	GetInput func() (string, bool)
};

func (opt *Opt) New() {
	opt.URL = flag.String(
		"url", "",
		"proxy URL e.g. vless://xxx");
	opt.IN_File = flag.String(
		"input", "",
		"path to proxy URLs file");

	flag.Parse();
}

func (opt *Opt) Set_rd_url() {
	opt.GetInput = func() (string, bool) {
		return *opt.URL, true
	}
}

func (opt *Opt) Set_rd_file() error {
	f, err := os.Open(*opt.IN_File)
	if nil != err {
		return err
	}
	opt.scanner = bufio.NewScanner(f)
	opt.GetInput = func() (string, bool) {
		if opt.scanner.Scan() {
			return opt.scanner.Text(), false
		} else {
			f.Close();
			return "", true
		}
	}
	return nil
}

func (opt *Opt) Set_rd_stdin() {
	if term.IsTerminal (int(os.Stdin.Fd())) {
		println ("Reading from STDIN until EOF:")
	}
	opt.scanner = bufio.NewScanner(os.Stdin)
	opt.GetInput = func() (string, bool) {
		if opt.scanner.Scan() {
			return opt.scanner.Text(), false
		} else {
			return "", true
		}
	}
}

// returns negative for fatal errors
func (opt *Opt) ParseFlags() int {
	argv := flag.Args()
	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "error:  missing COMMAND")
		fmt.Fprintln(os.Stderr, "usage:  v2utils [COMMAND] [OPTIONS]")
		return -1
	}
	switch (argv[0]) {
	case "convert","CONVERT", "conv", "c","C":
		opt.Cmd = CMD_CONVERT;
		break;
	case "test","Test","TEST", "t","T":
		opt.Cmd = CMD_TEST;
		break;
	case "run","Run","RUN", "r","R":
		opt.Cmd = CMD_RUN
		if "" == *opt.URL {
			println ("Run command needs a URL (--url).");
			return -1
		}
		break;

	default:
		println ("Invalid command.");
		return -1
	}
	return 0;
}

// returns negative on fatal failures
func (opt *Opt) Init() int {
	switch (opt.Cmd) {
	case CMD_RUN:
		opt.Set_rd_url()
		break;

	case CMD_TEST, CMD_CONVERT:
		if "" != *opt.URL {
			opt.Set_rd_url();
		} else if "" != *opt.IN_File {
			if e := opt.Set_rd_file(); nil != e {
				fmt.Printf ("%v\n", e);
				return -1
			}
		} else {
			opt.Set_rd_stdin();
		}
	}
	return 0;
}

func (opt Opt) Do() {
	EOF := false
	var ln string
	for true != EOF {
		ln, EOF = opt.GetInput();
		if len(ln) == 0 || ln[0] <= ' ' || ln[0] == '#' {
			continue;
		}

		switch (opt.Cmd) {
		case CMD_CONVERT:
			fmt.Printf("Not Implemented -- Converting `%s` --> json\n", ln);
			break;

		case CMD_RUN:
			fmt.Printf("Not Implemented -- Running xray-core with URL: `%s`\n", ln);
			break;

		case CMD_TEST:
			fmt.Printf ("Not Implemented -- Testing URL: `%s`\n", ln);
			break;
		}
	}
}


func main() {
	opt := Opt{};
	opt.New();
	if ret := opt.ParseFlags(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}
	opt.Do(); // The main loop
}
