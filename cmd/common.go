// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
	"strings"

	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"golang.org/x/term"

	pkg "github.com/siamak-amo/v2utils/pkg"
	log "github.com/siamak-amo/v2utils/log"

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
	rm *bool // remove files if broken or invalid

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
//         opt.template_path == "-" means to read from stdin
func (opt *Opt) Init_CFG() error {
	var e error
	if "-" == *opt.template_file {
		if term.IsTerminal (int(os.Stdin.Fd())) {
			println ("Reading json config from STDIN until EOF:")
		}
		res, err := ioutil.ReadAll(os.Stdin)
		if nil != err {
			return err
		}
		opt.CFG, e = pkg.Gen_main(string(res))
		return e;
	}

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


//////////////////////////////////////
//  Read URL/json helper functions  //
//////////////////////////////////////
func (opt *Opt) Set_rd_url() {
	opt.GetInput = func() (string, bool) {
		return *opt.url, true
	}
}

func (opt *Opt) Set_rd_file() error {
	f, err := os.Open(*opt.in_file)
	if nil != err {
		return err
	}
	opt.scanner = bufio.NewScanner(f)
	opt.GetInput = func() (string, bool) {
		if opt.scanner.Scan() {
			return opt.scanner.Text(), false
		} else {
			f.Close();
			return "", true
		}
	}
	return nil
}

func (opt *Opt) Set_rd_stdin() {
	if term.IsTerminal (int(os.Stdin.Fd())) {
		println ("Reading URLs from STDIN until EOF:")
	}
	opt.scanner = bufio.NewScanner(os.Stdin)
	opt.GetInput = func() (string, bool) {
		if opt.scanner.Scan() {
			return opt.scanner.Text(), false
		} else {
			return "", true
		}
	}
}

// Gives file path to json config files
func (opt *Opt) Set_rd_cfg_stdin() {
	opt.GetInput = func() (string, bool) { return "-", true; }
}

func (opt *Opt) Set_rd_cfg() {
	fileInfo, err := os.Stat(*opt.configs);
	if nil != err {
		log.Errorf("%v\n", err);
		opt.GetInput = func() (string, bool) { return "", true; }
		return;
	}
	if fileInfo.IsDir() {
		// Find all .json files here
		jsonFiles := []string{}
		filepath.Walk(
			*opt.configs,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasSuffix(info.Name(), ".json") {
					jsonFiles = append(jsonFiles, path)
				}
				return nil
			},
		);
		opt.GetInput = func() (string, bool) {
			if len(jsonFiles) == 0 {
				return "", true
			} else {
				_f := jsonFiles[0]
				jsonFiles = jsonFiles[1:]
				return _f, false
			}
		}
	} else {
		opt.GetInput = func() (string, bool) {
			if !strings.HasSuffix(*opt.configs, ".json") {
				return "", true
			}
			return *opt.configs , true;
		}
	}
}
