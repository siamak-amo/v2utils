// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"

	"hash/fnv"
	"encoding/hex"
	"path/filepath"

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
	CFG conf.Config
	Template core.ConfigSource

	url *string
	in_file *string // input URLs file path
	output_dir *string // output file(s) dir
	template_file *string // template file path

	scanner *bufio.Scanner
	GetInput func() (string, bool)
};

func (opt Opt) Out(buff []byte) (error) {
	if "" == *opt.output_dir {
		fmt.Println(string(buff));
	} else {
		path := opt.GetOutput_filepath(buff)
		of, err := os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644,
		);
		if err != nil {
			log.Errorf("Out failed to open file: %v\n", err)
			return err // fatal
		}
		defer of.Close()
		if _, err = of.Write(buff); nil != err {
			log.Errorf("Out failed to write: %v\n", err)
			return err // fatal
		}
		log.Infof("wrote %s", path)
	}
	return nil
}

// generates filename based on: hash(file_content)
func (opt Opt) GetOutput_filepath(file_content []byte) string {
	h := fnv.New64a()
	h.Write(file_content)
	return filepath.Join (*opt.output_dir,
		fmt.Sprintf("config_%s.json", hex.EncodeToString(h.Sum(nil))),
	)
}

func GetFormatByExtension(filename string) string {
	idx := strings.LastIndexByte(filename, '.')
	if idx == -1 {
		return ""
	}
	switch strings.ToLower(filename[idx+1:]) {
	case "pb", "protobuf":
		return "protobuf"
	case "yaml", "yml":
		return "yaml"
	case "toml":
		return "toml"
	case "json", "jsonc":
		return "json"
	default:
		return ""
	}
}

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
