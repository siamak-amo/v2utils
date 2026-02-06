// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"net/url"
	"encoding/json"
	"github.com/xtls/xray-core/infra/conf"
)

func Gen_trojan (args URLmap) (dst *conf.OutboundDetourConfig, e error) {
    map_normal (args, ServerPort, "443")
	dst = &conf.OutboundDetourConfig{}
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

func Gen_trojan_URL(src *conf.OutboundDetourConfig) *url.URL {
	var trojan TrojanCFG;
	if e := json.Unmarshal (*src.Settings, &trojan); nil != e {
		return nil
	}
	u := &url.URL{ Scheme: "trojan" }
	if 0 == len(trojan.Servers) {
		return nil
	}

	server := trojan.Servers[0]
	q := u.Query()
	u.User = url.User(server.Password)
	u.Host = fmt.Sprintf ("%s:%d", server.Address, server.Port)
	Init_trojanURL_stream (src.StreamSetting, q);

	u.RawQuery = q.Encode()
	return u
}
