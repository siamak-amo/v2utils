// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"

	core "github.com/xtls/xray-core/infra/conf"
)


func Gen_main(input string) (dst *core.Config, e error) {
	dst = &core.Config{}
	if e = unmarshal_H (dst, input); nil != e {
		// log
	}
	return
}

func Gen_inbound(input string) (dst []core.InboundDetourConfig, e error) {
	dst = make ([]core.InboundDetourConfig, 0)
	if e = unmarshal_H (&dst, input); nil != e {
		// log
	}
	return
}

// internal
func Gen_outbound(args URLmap) (dst []core.OutboundDetourConfig, e error) {
	dst = make ([]core.OutboundDetourConfig, 0)
	switch args[Protocol] {
	case "vless":
		v, e := Gen_vless (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Vless Error:  %v\n", e)
			return nil, e
		}
		break
	case "vmess":
		v, e := Gen_vmess (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Vmess Error:  %v\n", e)
			return nil, e
		}
		break

	case "ss","shadowsocks":
		v, e := Gen_ss (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Shadowsocks Error:  %v\n", e)
			return nil, e
		}
		break

	case "trojan":
		v, e := Gen_trojan (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Vmess Error:  %v\n", e)
			return nil, e
		}
		break

	default:
		return nil, not_implemented ("protocol " + args[Protocol])
	}
	return
}


