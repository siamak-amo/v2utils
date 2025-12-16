// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"encoding/json"

	log "github.com/siamak-amo/v2utils/log"
	pkg "github.com/siamak-amo/v2utils/pkg"
)

// Converts proxy URL @url to template @opt.Template
// Only returns error on fatal failures
func (opt Opt) Convert_url2json(url string) (error) {
	var err error
	if err = opt.Apply_template(); nil != err {
		return err
	}
	if err = opt.Init_Outbound_byURL(url); nil != err {
		log.Errorf("Invalid or unsupported URL - %v\n", err)
		return nil // not fatal
	}

	if err = opt.CFG_Out(url); nil != err {
		log.Errorf("%v", err)
		return err // IO error is fatal (invalid path / broken pipe)
	}
	return nil
}

// Makes output of opt.CFG, either on stdout or file
func (opt Opt) CFG_Out(url string) (error) {
	var err error
	var b []byte

	if "" == *opt.output_dir {
		// Use the compact format for stdout
		b, err = json.Marshal(opt.CFG)
	} else {
		b, err = json.MarshalIndent (opt.CFG, "", "    ")
	}
	if err != nil {
		panic(err) // it's ours. we have generated opt.CFG wrong
	}

	if "" == *opt.output_dir {
		println(b);
	} else {
		// Write to file
		path := opt.GetOutput_filepath([]byte(url))
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

// URL generator
func (opt Opt) Convert_conf2json() string {
	if 0 >= len(opt.CFG.OutboundConfigs) {
		return ""
	}
	if url := pkg.Gen_URL(&opt.CFG.OutboundConfigs[0]); nil != url {
		return url.String()
	}
	return ""
}
