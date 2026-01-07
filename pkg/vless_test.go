// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	"testing"
	"github.com/xtls/xray-core/infra/conf"
)


func TestGen_vless(t *testing.T) {
	tc := TestCase[VLessCFG] {T: t,
		Input: map[URLMapper]string {
	    	Protocol: "vless",
				ServerAddress:	"vpn.net",
				ServerPort:		"1234",
				Vxess_ID:		"68d3bd6c-829a-4ce9-92a6-edbebf4c5260",
				Vless_ENC:		"x509",
				Vless_Level:	"12",
			},
		Output: VLessCFG{},
	}
	
	v, e := Gen_vless (tc.Input)
	if nil != e {
		t.Fatalf ("gen_vless failed: %v\n", e)
		return
	}
	tc.Do(v);

	vnext := tc.Output.Settings.Vnext[0];
	user0 := vnext.Users[0];
	tc.Assert (tc.Output.Protocol,    tc.Input[Protocol])
	tc.Assert (vnext.Address,         tc.Input[ServerAddress])
	tc.Assert (vnext.Port,            tc.Input[ServerPort])
	tc.Assert (user0.ID,              tc.Input[Vxess_ID])
	tc.Assert (user0.Encryption,      tc.Input[Vless_ENC])
	tc.Assert (user0.Level,           tc.Input[Vless_Level])
}

// Test URL generator
const (
	FMT_Vless =`
		{
            "protocol": "vless", "settings": {
                "vnext": [{
                    "address": "1.2.3.4", "port": 1234, "users": [{
                        "id": "my_uuid", "encryption": "none_enc", "level": 666
                    }]
                }]
            }, "streamSettings": {%s}
        }`;
)

// Elementary test
func Test_Gen_vless_URL_1(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless, "")); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	Assert (t, u.Scheme, "vless");
	Assert (t, u.User.Username(), "my_uuid");
	Assert (t, u.Hostname(), "1.2.3.4");
	Assert (t, u.Port(), "1234");
	Assert (t, u.Query().Get("encryption"), "none_enc");
	Assert (t, u.Query().Get("level"), "666");
}

// Basic vless TCP + TLS test
func Test_Gen_vless_URL_2(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "tcp", "security": "tls",
         "tcpSettings": {"header": {"type": "none"}},
         "tlsSettings": {"allowInsecure": true, "serverName": "x.com",
                         "alpn": "h2,http/1.1", "fingerprint": "firefox-66"}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "tcp");
	Assert (t, q.Get("security"), "tls");
	Assert (t, q.Get("allowInsecure"), "true");
	Assert (t, q.Get("sni"), "x.com");
	Assert (t, q.Get("alpn"), "h2,http/1.1");
	Assert (t, q.Get("fp"), "firefox-66");
}

// Advanced vless TCP test
func Test_Gen_vless_URL_3(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "tcp", "security": "none", "tcpSettings": {
             "header": {
                 "type": "http",
                 "request": {"path": "/test", "headers": { "Host": "test.com" }}
             }
         }`),
	); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "tcp");
	Assert (t, q.Get("security"), "none");
	Assert (t, q.Get("headerType"), "http");
	Assert (t, q.Get("path"), "/test");
	Assert (t, q.Get("host"), "test.com");
}

// Vless over GRPC test
func Test_Gen_vless_URL_4(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "grpc", "security": "none_sec", "grpcSettings": {
             "authority": "abcd", "serviceName": "srv_name", "multiMode": true
         }`),
	); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "grpc");
	Assert (t, q.Get("security"), "none_sec");
	Assert (t, q.Get("authority"), "abcd");
	Assert (t, q.Get("serviceName"), "srv_name");
	Assert (t, q.Get("mode"), "multi");
}

// Vless over WS + Reality test
func Test_Gen_vless_URL_5(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "ws", "security": "reality",
         "wsSettings": {"host": "x.com", "path": "/test", "headers": {}},
         "realitySettings": {"serverNames": ["x.com", "y.com"],
                             "publicKey": "public.key", "shortId": "shortID",
                             "spiderX": "/test/spider", "fingerprint": "firefox-66"}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "ws");
	Assert (t, q.Get("security"), "reality");
	Assert (t, q.Get("host"), "x.com");
	Assert (t, q.Get("path"), "/test");
	Assert (t, q.Get("fp"), "firefox-66");
	Assert (t, q.Get("sni"), "x.com"); // we just picked the first one
	Assert (t, q.Get("pbk"), "public.key");
	Assert (t, q.Get("sid"), "shortID");
	Assert (t, q.Get("spx"), "/test/spider");
}

// Vless over KCP
func Test_Gen_vless_URL_6(t *testing.T) {
		cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "kcp", "security": "none", "kcpSettings": {
             "seed": "p@ssw0rd", "header": {"type": "utp"}
         }`),
	); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "kcp");
	Assert (t, q.Get("security"), "none");
	Assert (t, q.Get("path"), "p@ssw0rd");
	Assert (t, q.Get("headerType"), "utp");
}

// Vless over XHTTP
func Test_Gen_vless_URL_7(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "xhttp", "security": "none", "xhttpSettings": {
             "host": "x.com", "path": "/xpath"
         }`),
	); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "xhttp");
	Assert (t, q.Get("security"), "none");
	Assert (t, q.Get("path"), "/xpath");
	Assert (t, q.Get("host"), "x.com");
}

// Vless over HTTPUpgrade
func Test_Gen_vless_URL_8(t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "vless" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_Vless,
		`"network": "httpupgrade", "security": "none",
             "httpupgradeSettings": {
                 "host": "x.com", "path": "/http_upgrade"
             }`,
	)); nil != e {
		panic (e);
	}
	u := Gen_vless_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "httpupgrade");
	Assert (t, q.Get("security"), "none");
	Assert (t, q.Get("path"), "/http_upgrade");
	Assert (t, q.Get("host"), "x.com");
}
