// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	core "github.com/xtls/xray-core/infra/conf"
)

func Gen_ss(args URLmap) (dst *core.OutboundDetourConfig, e error) {
    map_normal (args, ServerPort, "443")
	dst = &core.OutboundDetourConfig{}
    if e = unmarshal_H (dst,
        fmt.Sprintf (
			`{
                "protocol": "shadowsocks",
                "settings": {
                  "servers": [
                      {
                          "address": "%s",
                          "port": %s,
                          "method": "%s",
                          "password": "%s"
                      }
                  ]
                },
                "tag": "proxy"
             }`,
			args[ServerAddress], args[ServerPort],
			args[SS_Method], args[SS_Password],
		),
	); nil != e {
		// log
	}
	return
} 
