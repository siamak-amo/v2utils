// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"bufio"
	"strings"

	"path/filepath"
	"golang.org/x/term"

	log "github.com/siamak-amo/v2utils/log"
)

var (
	global_scanner *bufio.Scanner
)

func Isatty(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func HasSuffixs(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.HasSuffix (haystack, needle) {
			return true;
		}
	}
	return false;
}


// Set_rd_xxx functions
// To initialize opt.Getinput to a function that returns
// proxy URL or `-` to read from stdin.

func (opt *Opt) Set_rd_url() {
	if (len(opt.url) == 1 && opt.url[0] == '-') || len(opt.url) == 0 {
		opt.Set_rd_stdin();
	} else {
		opt.GetInput = func() (string, bool) {
			return opt.url, true
		}
	}
}

func (opt *Opt) Set_rd_file() error {
	f, err := os.Open(opt.in_file)
	if nil != err {
		return err
	}
	global_scanner = bufio.NewScanner(f)
	opt.GetInput = func() (string, bool) {
		if global_scanner.Scan() {
			return global_scanner.Text(), false
		} else {
			f.Close();
			return "", true
		}
	}
	return nil
}

func (opt *Opt) Set_rd_stdin() {
	if Isatty(os.Stdin) {
		println ("Reading URLs from STDIN until EOF:")
	}
	global_scanner = bufio.NewScanner(os.Stdin)
	opt.GetInput = func() (string, bool) {
		if global_scanner.Scan() {
			return global_scanner.Text(), false
		} else {
			return "", true
		}
	}
}


// Set_rd_cfg_xxx functions
// To initialize opt.GetInput to a function that returns
// file path (.json, .toml, .yaml), or `-` to read from stdin.

func (opt *Opt) Set_rd_cfg_stdin() {
	opt.GetInput = func() (string, bool) { return "-", false; }
}

func (opt *Opt) Set_rd_cfg() {
	fileInfo, err := os.Stat(opt.configs);
	if nil != err {
		log.Errorf("%v\n", err);
		opt.GetInput = func() (string, bool) { return "", true; }
		return;
	}
	if fileInfo.IsDir() {
		jsonFiles := []string{}
		filepath.Walk(
			opt.configs,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if HasSuffixs(info.Name(), Supported_CFG_Formats) {
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
			if ! HasSuffixs(opt.configs, Supported_CFG_Formats) {
				log.Infof("file '%s' was ignored - invalid extension\n", opt.configs)
				return "", true
			}
			return opt.configs , true;
		}
	}
}
