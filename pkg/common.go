// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"io"
	"errors"
	"strings"

	"github.com/siamak-amo/v2utils/internal"

	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/infra/conf/serial"
	"github.com/xtls/xray-core/main/confloader"
	_ "github.com/xtls/xray-core/main/confloader/external"

)

type V2utils struct {
	CFG *conf.Config
	set_template bool
	Xray_instance *core.Instance // xray-core client instance
};


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

func (v2 *V2utils) UnsetTemplate() {
	if v2.set_template {
		v2.CFG = nil;
		v2.set_template = false;
	}
}

func (v2 *V2utils) HasTemplate() bool {
	return v2.set_template;
}

func (v2 *V2utils) HasInboundConfig() bool {
	return (nil != v2.CFG.InboundConfigs && 0 != len(v2.CFG.InboundConfigs));
}

func (v2 *V2utils) SetDefaultInboundConfig() {
	if c, e := internal.Gen_main(DEF_Run_Template); nil == e {
		v2.CFG.InboundConfigs = c.InboundConfigs;
	} else {
		panic(e); // it's ours, broken default template
	}
}

// Applies the template v2.template_path to v2.CFG
func (v2 *V2utils) Apply_template(file_path string) error {
	t := core.ConfigSource{
		Name: file_path,
		Format: GetFormatByExtension(file_path),
	}
	if "" == t.Format {
		return errors.New("invalid config file extension")
	}
	r, err := confloader.LoadConfig(t.Name)
	if nil != err {
		return err
	} else {
		c, err := serial.ReaderDecoderByFormat[t.Format](r)
		if nil != err {
			return err
		}
		v2.CFG = c;
	}
	v2.set_template = true
	return nil
}

func (v2 *V2utils) Apply_template_bystr(template string) error {
	var e error
	if v2.CFG, e = internal.Gen_main(template); nil != e {
		return e;
	}
	v2.set_template = true
	return nil
}

func (v2 *V2utils) Apply_template_byio(rio io.Reader) error {
	var err error
	v2.CFG, err = internal.Gen_main_io(rio);
	if nil == err {
		v2.set_template = true
	}
	return err;
}
