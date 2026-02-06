// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"io"
	"errors"
	"encoding/json"

	"github.com/siamak-amo/v2utils/internal"
)

// Initializes @v2.CFG.OutboundConfig by the provided proxy URL @url
func (v2 *V2utils) Init_Outbound_byURL(url string) (error) {
	var e error
	var umap internal.URLmap

	// Parse the URL
	umap, e = internal.ParseURL(url);
	if nil != e {
		return e
	}
	// Generate outbound config
	if v2.CFG.OutboundConfigs, e = internal.Gen_outbound(umap); nil != e {
		return e
	}
	return nil
}

func (v2 V2utils) Apply_URL(url string) (error) {
	var err error
	if err = v2.Init_Outbound_byURL(url); nil != err {
		return err
	}
	return nil
}

// Makes json output of v2.CFG, on @w
func (v2 V2utils) CFG_Out(w io.Writer, indention bool) (error) {
	encoder := json.NewEncoder(w)
	if indention {
		encoder.SetIndent("", "    ")
	}
	return encoder.Encode(v2.CFG);
}

// URL generator
func (v2 V2utils) Convert_conf2url() (string, error) {
	if 0 >= len(v2.CFG.OutboundConfigs) {
		return "", errors.New("Empty outbound configs")
	}
	url := internal.Gen_URL(&v2.CFG.OutboundConfigs[0]);
	if nil == url {
		return "", errors.New("Gen URL failed")
	}
	return url.String(), nil
}
