// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"strings"

	"crypto/md5"
	"encoding/hex"
	"path/filepath"

	log "github.com/siamak-amo/v2utils/log"
	pkg "github.com/siamak-amo/v2utils/pkg"
)

var (
	Supported_CFG_Formats = []string{".json", ".toml", ".yaml"}

	Stdin_is_tty = log.Isatty(os.Stdin)
	Stdout_is_tty = log.Isatty(os.Stdout)
)

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
	if "-" == opt.cfg {
		if Stdin_is_tty {
			println ("Reading json config from STDIN until EOF:")
		}
		return opt.v2.Apply_template_byio (os.Stdin);
	} else {
		if "" != opt.cfg {
			return opt.v2.Apply_template (opt.cfg)
		} else {
			opt.Apply_Default_Template();
		}
	}
	return nil
}

func result2string(result *pkg.TestResult) string {
	if "" != result.IP {
		return fmt.Sprintf("[IP: %s] (%dms)", result.IP, result.Duration);
	} else {
		return fmt.Sprintf("(%dms)", result.Duration);
	}
}

// makes long xray error messages shorter
func shortError(err string) string {
	idx := 0
	max_colon_count := 2
	for i := len(err)-1; i >= 0 && max_colon_count != 0; i -= 1 {
		if ':' == err[i] {
			idx = i+1
			max_colon_count -= 1
		}
	}
	return strings.TrimSpace(err[idx:]);
}

func (opt *Opt) get_contester() pkg.ConnectivityTester_I {
	if opt.verbose {
		// To also get the IP address of VPN
		return pkg.AdvancedTester;
	} else {
		// Simply test connectivity of VPN
		return pkg.SimpleTester;
	}
}

func (opt *Opt) Test_CFG() (bool) {
	err, result := opt.v2.Test_CFG(opt.cfg, opt.get_contester());

	if nil == err && opt.verbose {
		log.Infof("File '%s':  %s OK.\n", opt.cfg, result2string(result));
	}
	if ! opt.reverse {
		if nil == err {
			fmt.Println(opt.cfg);
		} else {
			log.Warnf("File '%s' is broken - %s\n", opt.cfg, shortError(err.Error()))
		}
	} else if nil != err { // Only print broken configs
		fmt.Println(opt.cfg);
	}
	return (err == nil);
}

func (opt *Opt) Test_URL() (bool) {
	err, result := opt.v2.Test_URL(opt.url, opt.get_contester());

	if nil == err && opt.verbose {
		log.Infof("URL '%s':  %s OK.\n", opt.url, result2string(result));
	}
	if ! opt.reverse {
		if nil == err {
			fmt.Println(opt.url)
		} else {
			log.Warnf("URL '%s' is broken - %s\n", opt.url, shortError(err.Error()));
		}
	} else if nil != err { // Only print broken urls
		fmt.Println(opt.url);
	}
	return (nil == err);
}

func (opt *Opt) Apply_URL() error {
	return opt.v2.Apply_URL(opt.url);
}
func (opt *Opt) Init_Outbound_byURL() error {
	return opt.v2.Init_Outbound_byURL(opt.url);
}
func (opt *Opt) Apply_template() error {
	return opt.v2.Apply_template(opt.cfg);
}

func (opt *Opt) Apply_Default_Template() {
	e := opt.v2.Apply_template_bystr( opt.Get_Default_Template() );
	if nil != e {
		panic(e); // it's ours, the default template is broken.
	}
}

func (opt Opt) Get_Default_Template() string {
	switch (opt.cmd) {
	case CMD_RUN_URL, CMD_RUN_CFG:
		return pkg.DEF_Run_Template;
	case CMD_TEST_URL, CMD_TEST_CFG:
		return pkg.DEF_Test_Template;
	case CMD_CONVERT_URL, CMD_CONVERT_CFG:
		return pkg.DEF_Run_Template;
	}
	return "";
}

func (opt Opt) MK_josn_output(url string) error {
	if "" == opt.output_dir {
		if err := opt.v2.CFG_Out(os.Stdout, !Stdout_is_tty); nil != err {
			return err;
		}
	} else {
		// Write to file
		path := opt.gen_output_filepath([]byte(url))
		of, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644,
		);
		if err != nil {
			return err
		}
		defer of.Close()
		if err := opt.v2.CFG_Out(of, true); nil != err {
			return err
		} else {
			log.Verbosef("Wrote: %s\n", path)
		}
	}
	return nil;
}
