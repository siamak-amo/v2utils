// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	"strings"
	"strconv"
	"net/url"
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
	if e := unmarshal_H (&vless, string(*src.Settings)); nil != e {
		return nil
	}
	if len(vless.Vnext) != 1  || len(vless.Vnext[0].Users) != 1 {
		return nil
	}
	vnext := vless.Vnext[0]
	u.User = url.User(vnext.Users[0].ID);
	u.Host = fmt.Sprintf("%s:%d", vnext.Address, vnext.Port)

	q := u.Query()
	AddQuery (q, "level", strconv.Itoa(vnext.Users[0].Level))
	AddQuery (q, "encryption", vnext.Users[0].Encryption)

	var stream conf.StreamConfig
	if nil != src.StreamSetting {
		stream = *src.StreamSetting

		if nil != stream.Network {
			net := string(*stream.Network)
			AddQuery (q, "type", net)
			switch (net) {
			case "tcp":
				if v,e := encode_tcp_header(stream.TCPSettings.HeaderConfig); nil == e {
					AddQuery (q, "headerType", v.Type)
					AddQuery (q, "path", v.Request.Path)
					AddQuery (q, "host", v.Request.Headers["Host"])
				}
				break;
			case "grpc":
				AddQuery (q, "serviceName", stream.GRPCSettings.ServiceName)
				AddQuery (q, "multiMode", strconv.FormatBool(stream.GRPCSettings.MultiMode))
				AddQuery (q, "authority", stream.GRPCSettings.Authority)
				break;
			case "ws":
				AddQuery (q, "host", stream.WSSettings.Host)
				AddQuery (q, "path", stream.WSSettings.Path)
				break;
			}
		}

		sec := stream.Security
		AddQuery (q, "security", sec)
		switch (sec) {
		case "tls":
			AddQuery (q, "allowInsecure", strconv.FormatBool(stream.TLSSettings.Insecure))
			AddQuery (q, "sni", stream.TLSSettings.ServerName)
			AddQuery (q, "fp", stream.TLSSettings.Fingerprint)
			AddQuery (q, "alpn", strings.Join(*stream.TLSSettings.ALPN, ","))
			break;
		case "reality":
			AddQuery (q, "fp", stream.REALITYSettings.Fingerprint)
			AddQuery (q, "spx", stream.REALITYSettings.SpiderX)
			AddQuery (q, "pbk", stream.REALITYSettings.PublicKey)
			AddQuery (q, "sid", stream.REALITYSettings.ShortId)
			AddQuery (q, "sni", stream.REALITYSettings.ServerNames[0])
			break;
		}
	}

	u.RawQuery = q.Encode()
	return u
}
