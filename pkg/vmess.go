// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	core "github.com/xtls/xray-core/infra/conf"
)

// internal
func Gen_vmess(args URLmap) (dst *core.OutboundDetourConfig, e error) {
    map_normal (args, Vmess_Sec, "none")
    map_normal (args, ServerPort, "443")
    map_normal (args, Vmess_AlterID, "0")
    dst = &core.OutboundDetourConfig{}
    if e = unmarshal_H (dst,
        fmt.Sprintf (
            `{
                "protocol": "%s",
                "settings": {"vnext": [
                  {
                    "address": "%s",
                    "port": %s,
                    "users": [{
                      "security": "%s",  "alterId": %s,  "id": "%s"
                    }]
                  }
                ]},
                "tag": "proxy"
             }`,
            args[Protocol], args[ServerAddress], args[ServerPort],
            args[Vmess_Sec], args[Vmess_AlterID], args[Vxess_ID],
        )); nil != e {
        // log
		return
    }
    if dst.StreamSetting, e = Gen_streamSettings (args); nil != e {
        // log
		return
    }
    return
}
