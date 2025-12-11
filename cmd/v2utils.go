// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"bufio"
	"golang.org/x/term"

	flag "github.com/spf13/pflag"
	log "github.com/siamak-amo/v2utils/log"
	"github.com/xtls/xray-core/core"
)

func (opt *Opt) RegisterFlag() {
	opt.url = flag.String(
		"url", "",
		"proxy URL e.g. vless://xxx");
	opt.in_file = flag.String(
		"input", "",
		"path to proxy URLs file");
	opt.template_file = flag.String(
		"template", "",
		"path to json template file");
	opt.output_dir = flag.String(
		"output", "", "output directory path")

	flag.Parse();
}

func (opt *Opt) Set_rd_url() {
	opt.GetInput = func() (string, bool) {
		return *opt.url, true
	}
}

func (opt *Opt) Set_rd_file() error {
	f, err := os.Open(*opt.in_file)
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

// returns negative on fatal failures
func (opt *Opt) ParseFlags() int {
	opt.RegisterFlag();
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
		if "" == *opt.url {
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
	// Initializing file paths
	if "" != *opt.template_file {
		opt.Template = core.ConfigSource{
			*opt.template_file,
			GetFormatByExtension(*opt.template_file),
		};
	}

	switch (opt.Cmd) {
	case CMD_RUN:
		opt.Set_rd_url()
		break;

	case CMD_TEST, CMD_CONVERT:
		if "" != *opt.output_dir {
			if err := os.MkdirAll(*opt.output_dir, 0o755); nil != err {
				log.Errorf ("Could not create dir - %v\n", err);
				return -1
			}
		}
		if "" != *opt.url {
			opt.Set_rd_url();
		} else if "" != *opt.in_file {
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
			if e := opt.Convert_url2json(ln); nil != e {
				return
			}
			break;

		case CMD_RUN:
			if "" != opt.Template.Name {
				opt.Apply_template (&opt.CFG)
			} else {
				// TODO: Use a default template instead of dying
				log.Errorf("no template provided! dying...\n");
				return;
			}
			if e := opt.Init_Outbound_byURL(ln); nil != e {
				break;
			}
			if e := opt.Run_Xray(); nil != e {
				break;
			} else {
				for ;; {};
				defer opt.Kill_Xray();
			}
			break;

		case CMD_TEST:
			if opt.Test_URL(ln) {
				// TODO: On verbose, print additional info
				println(ln)
			} else {
				log.Infof("Broken URL '%s'\n", ln);
			}
			break;
		}
	}
}

func main() {
	opt := Opt{};

	if ret := opt.ParseFlags(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}
	opt.Do(); // The main loop
}
