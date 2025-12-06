package pkg

import (
	"fmt"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
)

// internal
func Gen_vmess(args URLmap) (dst *v4.OutboundDetourConfig, e error) {
    map_normal (args, Vmess_Sec, "none")
    dst = &v4.OutboundDetourConfig{}
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
