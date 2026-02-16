// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"bufio"
	"strings"

	"path/filepath"

	log "github.com/siamak-amo/v2utils/log"
)

const (
	// read URL
	RURL_BUILTIN int = iota
	RURL_FILE
	// read config file
	RCFG_STDIN
	RCFG_FILE
)

var (
	read_method int = RURL_BUILTIN

	global_scanner *bufio.Scanner
	global_index int
	global_cfg_list []string
)

func HasSuffixs(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.HasSuffix (haystack, needle) {
			return true;
		}
	}
	return false;
}

func is_comment(input string) bool {
	return (0 == len(input) || input[0] <= ' ' || '#' == input[0]);
}

// @return:  end of inputs (bool)
func (opt *Opt) GetInput() (bool) {
	switch (read_method) {
	case RURL_BUILTIN:
		if global_index >= len(opt.urls) {
			return true
		} else {
			opt.url = opt.urls[global_index]
			global_index += 1
			return false
		}
		break;
	case RURL_FILE:
	read_url_from_file:
		if global_scanner.Scan() {
			opt.url = global_scanner.Text()
			if is_comment(opt.url) {
				goto read_url_from_file;
			}
			return false
		} else {
			return true
		}
		break;
	case RCFG_FILE:
	read_cfg_from_file:
		if global_index >= len(global_cfg_list) {
			return true;
		}
		opt.cfg = global_cfg_list[global_index];
		global_index += 1
		if is_comment(opt.cfg) {
			goto read_cfg_from_file;
		}
		return false
	case RCFG_STDIN:
		opt.cfg = "-"
		return false;
	}
	return true
}

func (opt *Opt) Set_rd_url() {
	if 0 == len(opt.urls) ||
		(1 == len(opt.urls) && opt.urls[0] == "-") {
		opt.Set_rd_stdin();
	} else {
		read_method = RURL_BUILTIN;
		global_index = 0
	}
}

func (opt *Opt) Set_rd_file() error {
	f, err := os.Open(opt.in_file)
	if nil != err {
		return err
	}
	global_scanner = bufio.NewScanner(f)
	read_method = RURL_FILE
	return nil
}

func (opt *Opt) Set_rd_stdin() {
	if Stdin_is_tty {
		println ("Reading URLs from STDIN until EOF:")
	}
	global_scanner = bufio.NewScanner(os.Stdin)
	read_method = RURL_FILE
}


// Set_rd_cfg_xxx functions
// To initialize opt.GetInput to a function that returns
// file path (.json, .toml, .yaml), or `-` to read from stdin.

func (opt *Opt) Set_rd_cfg_stdin() {
	read_method = RCFG_STDIN
}

func (opt *Opt) Set_rd_cfg() {
	if 0 == len(opt.configs) ||
		(1 == len(opt.configs) && opt.configs[0] == "-") {
		read_method = RCFG_STDIN;
	} else {
		global_index = 0
		read_method = RCFG_FILE;
		for _, path := range opt.configs {
			if path == "-" {
				global_cfg_list = append(global_cfg_list, path)
				continue;
			}
			fileInfo, err := os.Stat(path)
			if nil != err {
				log.Errorf("%v\n", err);
				continue;
			}
			if ! fileInfo.IsDir() {
				if ! HasSuffixs(path, Supported_CFG_Formats) {
					log.Infof("file '%s' was ignored - invalid extension\n", path)
				} else {
					global_cfg_list = append(global_cfg_list, path)
				}
			} else {
				filepath.Walk(path,
					func(file_path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if HasSuffixs(info.Name(), Supported_CFG_Formats) {
							global_cfg_list = append(global_cfg_list, file_path)
						}
						return nil
					},
				);
			}
		}
	}
}
