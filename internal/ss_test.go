// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"testing"
	"encoding/base64"
	"github.com/xtls/xray-core/infra/conf"
)


func TestGen_ss (t *testing.T) {
	tc := TestCase[ServerCFG] {T: t,
		Input: map[URLMapper]string {
	    	    Protocol:       "shadowsocks",
				ServerAddress:	"vpn.net",
				ServerPort:		"1234",
				SS_Method:      "chacha20-ietf-poly1305",
				SS_Password:    "p@ssw0rd",
			},
		Output: ServerCFG{},
	}

	v, e := Gen_ss (tc.Input)
	if nil != e {
		t.Fatalf ("gen_ss failed: %v\n", e)
		return
	}
	tc.Do(v);

	ss := tc.Output.Settings.Servers[0];
	tc.Assert (tc.Output.Protocol,    tc.Input[Protocol])
	tc.Assert (ss.Address,            tc.Input[ServerAddress])
	tc.Assert (ss.Port,               tc.Input[ServerPort])
	tc.Assert (ss.Method,             tc.Input[SS_Method])
	tc.Assert (ss.Password,           tc.Input[SS_Password])
}

// Generating shadowsocks URL
const (
	FMT_ss =`
		{
            "protocol": "shadowsocks", "settings": {
                "servers": [{
                    "address": "1.2.3.4", "port": 1234,
                    "method": "aes-777", "password": "p@ssw0rd"
                }]
            }, "streamSettings": {%s}
        }`;
)

func TestGen_ss_1 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "shadowsocks" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss, "")); nil != e {
		panic (e);
	}
	u := Gen_ss_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}
	v,e := base64.StdEncoding.DecodeString(u.User.Username())
	if nil != e {
		t.Fatal("Invalid URL")
	}

	Assert (t, u.Scheme, "ss");
	Assert (t, u.Host, "1.2.3.4:1234");
	Assert (t, string(v), "aes-777:p@ssw0rd")
}
