// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"fmt"
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
	pkg "github.com/siamak-amo/v2utils/pkg"
	log "github.com/siamak-amo/v2utils/log"
)

func (opt Opt) Convert_url2json(url string) {
	cf := conf.Config{}
	if "" != opt.Template.Name {
		opt.Apply_template(&cf)
	}

	var e error
	var umap pkg.URLmap
	umap, e = pkg.ParseURL(url);
	if nil != e {
		log.Errorf ("ParseURL failed - %s\n", e);
		return
	}
	cf.OutboundConfigs, e = pkg.Gen_outbound(umap);
	if nil != e {
		return
	}

	// TODO: skip null and empty strings
	// TODO: flag to print by indent or compact
	b, err := json.Marshal(cf)
	// b, err := json.MarshalIndent (cf, "", "    ")
	if err != nil {
		log.Errorf ("%v\n", e);
		return
	}
	fmt.Print (string(b));
}
