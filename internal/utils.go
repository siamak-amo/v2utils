// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"io"
	"net/url"
	log "github.com/siamak-amo/v2utils/log"
	"github.com/xtls/xray-core/infra/conf"
)


func Gen_main(input string) (dst *conf.Config, e error) {
	dst = &conf.Config{}
	if e = unmarshal_H (dst, input); nil != e {
		// log
	}
	return
}
func Gen_main_io(input io.Reader) (dst *conf.Config, e error) {
	dst = &conf.Config{}
	if e = unmarshal_HIO (dst, input); nil != e {
		// log
	}
	return
}

func Gen_inbound(input string) (dst []conf.InboundDetourConfig, e error) {
	dst = make ([]conf.InboundDetourConfig, 0)
	if e = unmarshal_H (&dst, input); nil != e {
		// log
	}
	return
}

// internal
func Gen_outbound(args URLmap) (dst []conf.OutboundDetourConfig, e error) {
	dst = make ([]conf.OutboundDetourConfig, 0)
	switch args[Protocol] {
	case "vless":
		v, e := Gen_vless (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			log.Errorf("Vless Error:  %v\n", e)
			return nil, e
		}
		break
	case "vmess":
		v, e := Gen_vmess (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			log.Errorf("Vmess Error:  %v\n", e)
			return nil, e
		}
		break

	case "ss","shadowsocks":
		v, e := Gen_ss (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			log.Errorf("Shadowsocks Error:  %v\n", e)
			return nil, e
		}
		break

	case "trojan":
		v, e := Gen_trojan (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			log.Errorf("Vmess Error:  %v\n", e)
			return nil, e
		}
		break

	default:
		return nil, not_implemented ("protocol " + args[Protocol])
	}
	return
}


func Gen_URL(src *conf.OutboundDetourConfig) *url.URL {
	switch (src.Protocol) {
	case "vless":
		return Gen_vless_URL(src);
	case "vmess":
		return Gen_vmess_URL(src);
	case "ss", "shadowsocks":
		return Gen_ss_URL(src);
	case "trojan":
		return Gen_trojan_URL(src);
	}
	return nil
}
