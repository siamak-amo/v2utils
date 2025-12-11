// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	core "github.com/xtls/xray-core/infra/conf"
)

func Gen_trojan (args URLmap) (dst *core.OutboundDetourConfig, e error) {
    map_normal (args, ServerPort, "443")
	dst = &core.OutboundDetourConfig{}
    if e = unmarshal_H (dst,
        fmt.Sprintf (
			`{
                "protocol": "trojan",
                "settings": {
                  "servers": [
                      {
                          "address": "%s",
                          "port": %s,
                          "password": "%s"
                      }
                  ]
                },
                "tag": "proxy"
             }`,
			args[ServerAddress], args[ServerPort], args[Trojan_Password],
		),
	); nil != e {
		// log
	}
	if dst.StreamSetting, e = Gen_streamSettings (args); nil != e {
        // log
		return
    }
	return
} 
