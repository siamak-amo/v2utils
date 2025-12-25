// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"net"
	"bufio"

	"crypto/md5"
	"encoding/hex"
	"path/filepath"

	log "github.com/siamak-amo/v2utils/log"
	utils "github.com/siamak-amo/v2utils/utils"	
)

const (
	CMD_CONVERT int = iota // URL
	CMD_CONVERT_CFG // json
	CMD_TEST
	CMD_TEST_CFG
	CMD_RUN
	CMD_RUN_CFG
) // commands

type Opt struct {
	// User options
	Cmd int                 // CMD_xxx
	url string
	configs string			// file or dir for testing
	output_dir string		// output file(s) dir
	template_file string	// template file path
	in_file string			// input URLs file path
	rm bool					// remove files if broken or invalid
	reverse bool            // print broken configs, not functionals
	verbose bool

	scanner *bufio.Scanner
	GetInput func() (string, bool)

	V2 utils.V2utils
};

// generates filename based on: hash(url)
func (opt Opt) gen_output_filepath(url []byte) string {
	h := md5.New()
	h.Write(url)
	return filepath.Join (opt.output_dir,
		fmt.Sprintf(
			"config_%s.json",
			hex.EncodeToString(h.Sum(nil))[0:16],
		),
	)
}

// Applies opt.template_path  or  the default template
//         opt.template_path == "-" means to read from stdin
func (opt *Opt) Init_CFG() error {
	var e error
	if "-" == opt.template_file {
		if Isatty(os.Stdin) {
			println ("Reading json config from STDIN until EOF:")
		}
		e = opt.V2.Apply_template_byio (os.Stdin);
		return e;
	} else {
		if "" != opt.template_file {
			return opt.V2.Apply_template (opt.template_file)
		} else {
			return opt.V2.Apply_template_bystr (opt.Get_Default_Template());
		}
	}
	return nil
}

func (opt *Opt) Get_Default_Template() string {
	switch (opt.Cmd) {
	case CMD_RUN:
		return DEF_Run_Template;
	case CMD_TEST:
		return DEF_Test_Template;
	case CMD_CONVERT:
		return DEF_Run_Template;
	}
	return "";
}

func (opt Opt) MK_josn_output(url string) {
	if "" == opt.output_dir {
		opt.V2.CFG_Out(os.Stdout, !Isatty(os.Stdout)); // no indent for tty's
	} else {
		// Write to file
		path := opt.gen_output_filepath([]byte(url))
		of, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644,
		);
		if err != nil {
			log.Errorf("File error - %v\n", err)
			return
		}
		defer of.Close()
		if err := opt.V2.CFG_Out(of, true); nil != err {
			log.Errorf("Write error - %v\n", err)
			return
		} else {
			log.Verbosef("Wrote: %s\n", path)
		}
	}
}

func (opt Opt) Test_URL(url string) bool {
	opt.V2.Apply_template_bystr(DEF_Test_Template);
	return opt.V2.Test_URL(url)
}

// PickPort returns an unused TCP port
// The port returned is highly likely to be unused, but not guaranteed.
func PickPort() int {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if nil != err {
		return -1
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
