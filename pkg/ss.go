// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"

	"net/url"
	"encoding/json"
	"encoding/base64"

	"github.com/xtls/xray-core/infra/conf"
)

func Gen_ss(args URLmap) (dst *conf.OutboundDetourConfig, e error) {
    map_normal (args, ServerPort, "443")
	dst = &conf.OutboundDetourConfig{}
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
	if dst.StreamSetting, e = Gen_streamSettings (args); nil != e {
        // log
		return
    }
	return
} 

func Gen_ss_URL(src *conf.OutboundDetourConfig) *url.URL {
	var ss SSCFG;
	if e := json.Unmarshal (*src.Settings, &ss); nil != e {
		return nil
	}
	u := &url.URL{ Scheme: "ss" }
	if 0 == len(ss.Servers) {
		return nil
	}

	server := ss.Servers[0]
	u.User = url.User (base64.StdEncoding.EncodeToString (
		[]byte(fmt.Sprintf ("%s:%s", server.Method, server.Password)),
	))
	u.Host = fmt.Sprintf ("%s:%d", server.Address, server.Port)

	// q := u.Query()
	// Init_ssURL_stream(src.StreamSetting, u.Query());
	// u.RawQuery = q.Encode()
	return u
}
