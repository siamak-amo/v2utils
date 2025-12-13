// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"encoding/json"

	log "github.com/siamak-amo/v2utils/log"
)

// Converts proxy URL @url to template @opt.Template
// returns error Only on fatal failures
func (opt Opt) Convert_url2json(url string) (error) {
	if "" != opt.Template.Name {
		opt.Apply_template(&opt.CFG)
	}
	if e := opt.Init_Outbound_byURL(url); nil != e {
		return nil // not fatal
	}

	if err = opt.CFG_Out(url); nil != err {
		log.Errorf("%v", err)
		return err // IO error is fatal (invalid path / broken pipe)
	}
	return nil
}

func (opt Opt) CFG_Out() (error) {
	// TODO: skip null and empty strings
	// TODO: flag to print by indent or compact
	b, err := json.Marshal(opt.CFG)
	// b, err := json.MarshalIndent (cf, "", "    ")
	if err != nil {
		log.Errorf ("json.Marshal failed - %v\n", err);
		return nil // not fatal
	}

	if "" == *opt.output_dir {
		// Using stdout
		println(b);
	} else {
		// Write to file
		path := opt.GetOutput_filepath(b)
		of, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644,
		);
		if err != nil {
			log.Errorf("Out failed to open file: %v\n", err)
			return err // fatal
		}
		defer of.Close()
		if _, err = of.Write(b); nil != err {
			log.Errorf("Out failed to write: %v\n", err)
			return err // fatal
		}
		log.Verbosef("Wrote: %s", path)
	}
	return nil
}
