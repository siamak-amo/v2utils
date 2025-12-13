// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"

	"path/filepath"
	"golang.org/x/term"

	flag "github.com/spf13/pflag"
	log "github.com/siamak-amo/v2utils/log"
)

func (opt *Opt) RegisterFlag() {
	opt.url = flag.String ("url", "", "proxy URL e.g. vless://xxx");
	opt.in_file = flag.String ("input", "", "path to proxy URLs file");
	opt.template_file = flag.String ("template", "", "path to json template file");
	opt.output_dir = flag.String ("output", "", "output directory path");
	opt.configs = flag.String ("config", "", "path to config file or dir");
	opt.Verbose = flag.Bool ("verbose", false, "verbose");

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

// Gives file path to json config files
func (opt *Opt) Set_rd_cfg_stdin() {
	opt.GetInput = func() (string, bool) { return "-", true; }
}

func (opt *Opt) Set_rd_cfg() {
	fileInfo, err := os.Stat(*opt.configs);
	if nil != err {
		log.Errorf("%v\n", err);
		opt.GetInput = func() (string, bool) { return "", true; }
		return;
	}
	if fileInfo.IsDir() {
		// Find all .json files here
		jsonFiles := []string{}
		filepath.Walk(
			*opt.configs,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(info.Name(), ".json") {
					jsonFiles = append(jsonFiles, path)
				}
				return nil
			},
		);
		opt.GetInput = func() (string, bool) {
			if len(jsonFiles) == 0 {
				return "", true
			} else {
				_f := jsonFiles[0]
				jsonFiles = jsonFiles[1:]
				return _f, false
			}
		}
	} else {
		opt.GetInput = func() (string, bool) {
			if !strings.HasSuffix(*opt.configs, ".json") {
				return "", true
			}
			return *opt.configs , true;
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
		if "" != *opt.configs {
			opt.Cmd = CMD_CONVERT_CFG;
		} else {
			// Default: converting URLs (stdin / --url)
			opt.Cmd = CMD_CONVERT;
		}
		break;

	case "test","Test","TEST", "t","T":
		if "" != *opt.configs {
			// For Testing config (.json) files
			opt.Cmd = CMD_TEST_CFG;
		} else {
			// Default: testing URLs
			opt.Cmd = CMD_TEST;
		}
		break;

	case "run","Run","RUN", "r","R":
		if "" != *opt.configs {
			opt.Cmd = CMD_RUN_CFG;
		} else if "" != *opt.url {
			opt.Cmd = CMD_RUN;
		} else {
			log.Errorf("Run command needs a URL (--url) or config file (--config).");
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
				log.Errorf ("%v\n", e);
				return -1
			}
		} else {
			opt.Set_rd_stdin();
		}
		break;

	case CMD_CONVERT_CFG, CMD_TEST_CFG, CMD_RUN_CFG:
		if "-" == *opt.configs {
			opt.Set_rd_cfg_stdin();
		} else {
			opt.Set_rd_cfg();
		}
		break;
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
			if "" == *opt.template_file {
				log.Errorf("No template provided, using the default template\n");
				fmt.Printf("Default template:%s\n", opt.Get_Default_Template()); // should print this
			}
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Invalid template - %v\n", e)
				return;
			}
			if e := opt.Init_Outbound_byURL(ln); nil != e {
				log.Errorf("Invalid or unsupported URL - %v\n", e)
				break;
			}
			if e := opt.Exec_Xray(); nil != e {
				log.Errorf("Exec xray-core failed - %v\n", e)
				break;
			}
			break;

		case CMD_TEST:
			if opt.Test_URL(ln) {
				if *opt.Verbose {
					log.Infof("`%s` OK.\n", ln)
				} else {
					println(ln)
				}
				if "" != *opt.output_dir {
					opt.CFG_Out(ln); // Also generate json files
				}
			} else {
				log.Infof("Broken URL '%s'\n", ln);
			}
			break;

		case CMD_TEST_CFG:
			log.Errorf("(%s) CMD_TEST_CFG -- Not Implemented.\n", ln);
			break;

		case CMD_RUN_CFG:
			log.Errorf("(%s) CMD_RUN_CFG -- Not Implemented.\n", ln);
			break;

		case CMD_CONVERT_CFG:
			log.Errorf("(%s) CMD_CONVERT_CFG -- Not Implemented.\n", ln);
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
	if *opt.Verbose {
		log.LogLevel = log.Verbose;
	}
	opt.Do(); // The main loop
}
