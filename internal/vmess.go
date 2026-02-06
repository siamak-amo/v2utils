// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"strconv"

	"net/url"
	"encoding/json"
	"encoding/base64"

	"github.com/xtls/xray-core/infra/conf"
)

// internal
func Gen_vmess(args URLmap) (dst *conf.OutboundDetourConfig, e error) {
    map_normal (args, Vmess_Sec, "none")
    map_normal (args, ServerPort, "443")
    map_normal (args, Vmess_AlterID, "0")
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

func Gen_vmess_URL(src *conf.OutboundDetourConfig) *url.URL {
	var vmess VmessVnext
	if e := json.Unmarshal (*src.Settings, &vmess); nil != e {
		return nil
	}
	if len(vmess.Vnext) == 0 || len(vmess.Vnext[0].Users) == 0 {
		return nil
	}
	vnext := vmess.Vnext[0]
	res := make(map[string]string, 0)
	res["add"] = vnext.Address
	res["port"] = strconv.Itoa(vnext.Port)
	res["id"] = vnext.Users[0].ID
	res["aid"] = strconv.Itoa(vnext.Users[0].AlterIds)
	res["scy"] = vnext.Users[0].Security
	Init_vmessURL_stream (src.StreamSetting, res);
	
	j, err := json.Marshal(res);
	if nil != err {
		panic(err); // it's ours.
	}
	u := &url.URL{
		Scheme: "vmess",
		Host: base64.StdEncoding.EncodeToString(j),
	};
	return u
}
