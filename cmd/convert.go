// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"fmt"
	"encoding/json"

	log "github.com/siamak-amo/v2utils/log"
)

// Converts proxy URL @url to template @opt.Template
func (opt Opt) Convert_url2json(url string) {
	// Apply the template if provided
	if "" != opt.Template.Name {
		opt.Apply_template(opt.CFG)
	}
	if e := opt.Init_Outbound_byURL(url); nil != e {
		return
	}

	// TODO: skip null and empty strings
	// TODO: flag to print by indent or compact
	b, err := json.Marshal(opt.CFG)
	// b, err := json.MarshalIndent (cf, "", "    ")
	if err != nil {
		log.Errorf ("%v\n", err);
		return
	}
	fmt.Print (string(b));
}
