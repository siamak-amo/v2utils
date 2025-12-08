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
func Gen_outbound(args URLmap, template string) (dst []core.OutboundDetourConfig, e error) {
	dst = make ([]core.OutboundDetourConfig, 0)
	if e = unmarshal_H (&dst, template); nil != e {
		// log
		return
	}
	switch args[Protocol] {
	case "vless":
		v, e := Gen_vless (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Vless Error:  %v\n", e)
		}
		break
	case "vmess":
		v, e := Gen_vmess (args)
		if e == nil {
			dst = append (dst, *v)
		} else {
			fmt.Printf("Vmess Error:  %v\n", e)
		}
		break

	default:
		return nil, not_implemented ("protocol " + args[Protocol])
	}
	return
}


