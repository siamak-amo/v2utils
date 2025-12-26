package main

import (
	"os"
	"bufio"
	"strings"

	"path/filepath"
	"golang.org/x/term"

	log "github.com/siamak-amo/v2utils/log"
)

func Isatty(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func (opt *Opt) Set_rd_url() {
	if opt.url[0] == '-' {
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
	if Isatty(os.Stdin) {
		println ("Reading URLs from STDIN until EOF:")
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
		// Find all .json files here
		jsonFiles := []string{}
		filepath.Walk(
			opt.configs,
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
			if !strings.HasSuffix(opt.configs, ".json") {
				return "", true
			}
			return opt.configs , true;
		}
	}
}
