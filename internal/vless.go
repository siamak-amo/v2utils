// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"strconv"
	"net/url"
	"encoding/json"
	"github.com/xtls/xray-core/infra/conf"
)

// internal
func Gen_vless(args URLmap) (dst *conf.OutboundDetourConfig, e error) {
    map_normal (args, Vless_ENC, "none")
    map_normal (args, ServerPort, "443")
    map_normal (args, Vless_Level, "0")
    dst = &conf.OutboundDetourConfig{}
    if e = unmarshal_H (dst,
        fmt.Sprintf (
            `{
                "protocol": "%s",
                "settings": {"vnext": [
                  {
                    "address": "%s",
                    "port": %s,
                    "users": [{
                      "encryption": "%s",  "level": %s,  "id": "%s"
                    }]
                  }
                ]},
                "tag": "proxy"
             }`,
            args[Protocol], args[ServerAddress], args[ServerPort],
            args[Vless_ENC], args[Vless_Level], args[Vxess_ID],
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

// URL generator
func Gen_vless_URL(src *conf.OutboundDetourConfig) *url.URL {
	var vless VLessVnext
	u := &url.URL{ Scheme: "vless" };
	if e := json.Unmarshal (*src.Settings, &vless); nil != e {
		return nil
	}
	if len(vless.Vnext) == 0 || len(vless.Vnext[0].Users) == 0 {
		return nil
	}
	vnext := vless.Vnext[0]
	u.User = url.User(vnext.Users[0].ID);
	u.Host = fmt.Sprintf("%s:%d", vnext.Address, vnext.Port)

	q := u.Query()
	AddQuery (q, "level", strconv.Itoa(vnext.Users[0].Level))
	AddQuery (q, "encryption", vnext.Users[0].Encryption)
	Init_vlessURL_stream(src.StreamSetting, q)

	u.RawQuery = q.Encode()
	return u
}
