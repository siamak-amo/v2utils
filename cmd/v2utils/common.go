// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"

	"crypto/md5"
	"encoding/hex"
	"path/filepath"

	log "github.com/siamak-amo/v2utils/log"
	pkg "github.com/siamak-amo/v2utils/pkg"
)

var (
	Supported_CFG_Formats = []string{".json", ".toml", ".yaml"}
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
		if Isatty(os.Stdin) {
			println ("Reading json config from STDIN until EOF:")
		}
		return opt.V2.Apply_template_byio (os.Stdin);
	} else {
		if "" != opt.cfg {
			return opt.V2.Apply_template (opt.cfg)
		} else {
			opt.Apply_Default_Template();
		}
	}
	return nil
}

func (opt *Opt) Test_CFG() (bool, int64) {
	return opt.V2.Test_CFG(opt.cfg);
}
func (opt *Opt) Test_URL() (bool, int64) {
	return opt.V2.Test_URL(opt.url);
}

func (opt *Opt) Apply_URL() error {
	return opt.V2.Apply_URL(opt.url);
}
func (opt *Opt) Init_Outbound_byURL() error {
	return opt.V2.Init_Outbound_byURL(opt.url);
}
func (opt *Opt) Apply_template() error {
	return opt.V2.Apply_template(opt.cfg);
}

func (opt *Opt) Apply_Default_Template() {
	e := opt.V2.Apply_template_bystr( opt.Get_Default_Template() );
	if nil != e {
		panic(e); // it's ours, the default template is broken.
	}
}

func (opt Opt) Get_Default_Template() string {
	switch (opt.Cmd) {
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
		if err := opt.V2.CFG_Out(os.Stdout, !Isatty(os.Stdout)); nil != err {
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
		if err := opt.V2.CFG_Out(of, true); nil != err {
			return err
		} else {
			log.Verbosef("Wrote: %s\n", path)
		}
	}
	return nil;
}
