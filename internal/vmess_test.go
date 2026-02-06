// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"errors"
	"net/url"
	"testing"
	"encoding/json"
	"encoding/base64"
	"github.com/xtls/xray-core/infra/conf"
)


func TestGen_vmess(t *testing.T) {
	tc := TestCase[VMessCFG] {T: t,
		Input: map[URLMapper]string {
	    	Protocol: "vmess",
				ServerAddress:	"vpn.net",
				ServerPort:		"1234",
				Vxess_ID:		"68d3bd6c-829a-4ce9-92a6-edbebf4c5260",
				Vmess_Sec:		"secure",
				Vmess_AlterID:	"666",
			},
		Output: VMessCFG{},
	}

	v, e := Gen_vmess (tc.Input)
	if nil != e {
		t.Fatalf ("gen_vless failed: %v\n", e)
		return
	}
	tc.Do(v);

	vnext := tc.Output.Settings.Vnext[0];
	user0 := vnext.Users[0];
	tc.Assert (tc.Output.Protocol,     tc.Input[Protocol])
	tc.Assert (vnext.Address,          tc.Input[ServerAddress])
	tc.Assert (vnext.Port,             tc.Input[ServerPort])
	tc.Assert (user0.ID,               tc.Input[Vxess_ID])
	tc.Assert (user0.Security,         tc.Input[Vmess_Sec])
	tc.Assert (user0.AlterIds,         tc.Input[Vmess_AlterID])
}

// URL generator tests
const (
	FMT_Vmess =`
		{
            "protocol": "vmess", "settings": {
                "vnext": [{
                    "address": "1.2.3.4", "port": 1234, "users": [{
                        "id": "my_uuid", "security": "aes-128-gcm"
                    }]
                }]
            }, "streamSettings": {%s}
        }`;
)

func mk_vmess_map(src *url.URL) (map[string]string, error) {
	if nil == src || len(src.String()) < 8 {
		return nil, errors.New("Invalid src")
	}
	j, e := base64.StdEncoding.DecodeString (src.String()[8:]);
	if nil != e {
		return nil, e
	}
	res := make (map[string]string, 0)
	if e := json.Unmarshal (j, &res);nil != e {
		return nil, e
	}
	return res, nil
}

// Elementary test
func Test_Gen_vmess_URL_1(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vmess" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vmess, "")); nil != e {
		panic (e);
	}
	u := Gen_vmess_URL (cfg);
	res, e := mk_vmess_map(u);
	if nil != e {
		t.Fatalf("Converting url to json failed - %v\n", e)
	}

	Assert (t, u.Scheme, "vmess");
	Assert (t, res["id"], "my_uuid");
	Assert (t, res["add"], "1.2.3.4");
	Assert (t, res["scy"], "aes-128-gcm");
}

// Basic vmess over TCP + TLS test
func Test_Gen_vmess_URL_2(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vmess" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vmess,
		`"network": "tcp", "security": "tls",
         "tcpSettings": {"header": {"type": "none"}},
         "tlsSettings": {"allowInsecure": true, "serverName": "x.com",
                         "alpn": "h2,http/1.1", "fingerprint": "firefox-66"}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vmess_URL (cfg);
	res, e := mk_vmess_map(u);
	if nil != e {
		t.Fatalf("Converting url to json failed - %v\n", e)
	}

	Assert (t, res["net"], "tcp");
	Assert (t, res["tls"], "tls");
	Assert (t, res["sni"], "x.com");
	Assert (t, res["alpn"], "h2,http/1.1");
	Assert (t, res["fp"], "firefox-66");
	Assert (t, res["allowInsecure"], "true");
}

// Advanced TCP test
func Test_Gen_vmess_URL_3(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vmess" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vmess,
		`"network": "tcp", "security": "none", "tcpSettings": {
             "header": {
                 "type": "http",
                 "request": {"path": "/test", "headers": { "Host": "test.com" }}
             }
         }`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vmess_URL (cfg);
	res, e := mk_vmess_map(u);
	if nil != e {
		t.Fatalf("Converting url to json failed - %v\n", e)
	}

	Assert (t, res["net"], "tcp");
	Assert (t, res["tls"], "");
	Assert (t, res["type"], "http");
	Assert (t, res["path"], "/test");
	Assert (t, res["host"], "test.com");
}

// Vmess over WebSocket test
func Test_Gen_vmess_URL_4(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vmess" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vmess,
		`"network": "ws", "security": "none",
         "wsSettings": {"host": "x.com", "path": "/path", "headers": {}}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vmess_URL (cfg);
	res, e := mk_vmess_map(u);
	if nil != e {
		t.Fatalf("Converting url to json failed - %v\n", e)
	}

	Assert (t, res["net"], "ws");
	Assert (t, res["tls"], "");
	Assert (t, res["type"], "none");
	Assert (t, res["path"], "/path");
	Assert (t, res["host"], "x.com");
}

// Vmess over GRPC test
func Test_Gen_vmess_URL_5(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vmess" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vmess,
		`"network": "grpc", "security": "none", "grpcSettings": {
             "authority": "abcd", "serviceName": "srv_name", "multiMode": true
        }`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vmess_URL (cfg);
	res, e := mk_vmess_map(u);
	if nil != e {
		t.Fatalf("Converting url to json failed - %v\n", e)
	}

	Assert (t, res["net"], "grpc");
	Assert (t, res["tls"], "");
	Assert (t, res["type"], "multi");
	Assert (t, res["authority"], "abcd");
	Assert (t, res["path"], "srv_name");
}
