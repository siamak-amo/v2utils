// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"

	"crypto/md5"
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
	Verbose *bool

	url *string
	in_file *string // input URLs file path
	output_dir *string // output file(s) dir
	template_file *string // template file path

	scanner *bufio.Scanner
	GetInput func() (string, bool)

	Xray_instance *core.Instance // xray-core client instance
};

// generates filename based on: hash(url)
func (opt Opt) GetOutput_filepath(url []byte) string {
	h := md5.New()
	h.Write(url)
	return filepath.Join (*opt.output_dir,
		fmt.Sprintf(
			"config_%s.json",
			hex.EncodeToString(h.Sum(nil))[0:16],
		),
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

// PickPort returns an unused TCP port
// The port returned is highly likely to be unused, but not guaranteed.
func PickPort() int {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if nil != err {
		return -1
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
