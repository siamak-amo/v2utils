// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"bufio"
	pkg "github.com/siamak-amo/v2utils/pkg"
	log "github.com/siamak-amo/v2utils/log"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	"github.com/xtls/xray-core/main/confloader"
	_ "github.com/xtls/xray-core/main/confloader/external"
)

const (
	CMD_CONVERT int = iota
	CMD_TEST
	CMD_RUN
) // commands

type Opt struct {
	Cmd int // CMD_xxx
	CFG *conf.Config
	Template core.ConfigSource

	url *string
	in_file *string  // input URLs file path
	template_file *string // template file path

	scanner *bufio.Scanner
	GetInput func() (string, bool)
};

// Applies the template opt.Template to @dst
func (opt *Opt) Apply_template(dst *conf.Config) {
	r, err := confloader.LoadConfig(opt.Template.Name)
	if nil != err {
		log.Errorf("%v\n", err);
	} else {
		c, err := serial.ReaderDecoderByFormat[opt.Template.Format](r)
		if nil != err {
			log.Errorf("%v\n", err);
		} else {
			*dst = *c  // for the first time

			// TODO: maybe accept template array and merge them via:
			// dst.Override(c, file.Name)
		}
	}
}

// Initializes @opt.Cfg by provided proxy URL @url
func (opt *Opt) Init_Outbound_byURL(url string) (error) {
	var e error
	var umap pkg.URLmap

	// Parse the URL
	umap, e = pkg.ParseURL(url);
	if nil != e {
		log.Errorf ("ParseURL failed - %s\n", e);
		return e
	}
	// Generate outbound config
	if opt.CFG.OutboundConfigs, e = pkg.Gen_outbound(umap); nil != e {
		return e
	}
	return nil
}
