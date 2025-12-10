// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	log "github.com/siamak-amo/v2utils/log"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	"github.com/xtls/xray-core/main/confloader"
	_ "github.com/xtls/xray-core/main/confloader/external"
)

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
