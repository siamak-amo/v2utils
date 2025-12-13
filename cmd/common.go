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

	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	"github.com/xtls/xray-core/main/confloader"
	_ "github.com/xtls/xray-core/main/confloader/external"
)

const (
	CMD_CONVERT int = iota // URL
	CMD_CONVERT_CFG // json
	CMD_TEST
	CMD_TEST_CFG
	CMD_RUN
	CMD_RUN_CFG
) // commands

type Opt struct {
	Cmd int // CMD_xxx
	CFG *conf.Config
	Verbose *bool

	url *string
	in_file *string // input URLs file path
	configs *string // file or dir for testing
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

// Applies the template opt.template_path to opt.CFG
func (opt *Opt) Apply_template() error {
	t := core.ConfigSource{
		Name: *opt.template_file,
		Format: GetFormatByExtension(*opt.template_file),
	}
	r, err := confloader.LoadConfig(t.Name)
	if nil != err {
		return err
	} else {
		c, err := serial.ReaderDecoderByFormat[t.Format](r)
		if nil != err {
			return err
		}
		opt.CFG = c;
	}
	return nil
}

// Applies opt.template_path  or  the default template
func (opt *Opt) Init_CFG() error {
	var e error
	if "" != *opt.template_file {
		return opt.Apply_template()
	} else {
		// Apply the default template
		if opt.CFG, e = pkg.Gen_main(opt.Get_Default_Template()); nil != e {
			panic(e) // it's ours. broken default template
		}
	}
	return nil
}

func (opt *Opt) Get_Default_Template() string {
	switch (opt.Cmd) {
	case CMD_RUN:
		return DEF_Run_Template;
	case CMD_TEST:
		return DEF_Test_Template;
	case CMD_CONVERT:
		return DEF_Run_Template;
	}
	return "";
}

// Initializes @opt.CFG.OutboundConfig by the provided proxy URL @url
func (opt *Opt) Init_Outbound_byURL(url string) (error) {
	var e error
	var umap pkg.URLmap

	// Parse the URL
	umap, e = pkg.ParseURL(url);
	if nil != e {
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
