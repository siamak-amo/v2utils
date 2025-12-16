// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	"testing"
	"github.com/xtls/xray-core/infra/conf"
)

func TestGen_trojan (t *testing.T) {
	tc := TestCase[ServerCFG] {T: t,
		Input: map[URLMapper]string {
	    	    Protocol:           "trojan",
				ServerAddress:	    "vpn.net",
				ServerPort:		    "1234",
				Trojan_Password:    "p@ssw0rd",
			},
		Output: ServerCFG{},
	}

	v, e := Gen_trojan (tc.Input)
	if nil != e {
		t.Fatalf ("gen_trojan failed: %v\n", e)
		return
	}
	tc.Do(v);

	trojan := tc.Output.Settings.Servers[0];
	tc.Assert (tc.Output.Protocol,        tc.Input[Protocol])
	tc.Assert (trojan.Address,            tc.Input[ServerAddress])
	tc.Assert (trojan.Port,               tc.Input[ServerPort])
	tc.Assert (trojan.Password,           tc.Input[Trojan_Password])
}


// Generating trojan URL
const (
	FMT_trojan =`
		{
            "protocol": "trojan", "settings": {
                "servers": [{
                    "address": "1.2.3.4", "port": 1234, "password": "p@ssw0rd"
                }]
            }, "streamSettings": {%s}
        }`;
)

// Elementary test
func TestGen_trojan_1 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "trojan" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss, "")); nil != e {
		panic (e);
	}
	u := Gen_trojan_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	Assert (t, u.Scheme, "trojan");
	Assert (t, u.Host, "1.2.3.4:1234");
	Assert (t, u.User.Username(), "p@ssw0rd")
}

// Basic TCP + TLS test
func TestGen_trojan_2 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "trojan" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss,
		`"network": "tcp", "security": "tls",
         "tcpSettings": {"header": {"type": "none"}},
         "tlsSettings": {"allowInsecure": true, "serverName": "x.com",
                         "alpn": "h2,http/1.1", "fingerprint": "firefox-66"}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_trojan_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

// trojan://auto@104.18.12.229:2053?path=%2F&security=tls&host=e9464f45.trauma-2r4.pages.dev&type=ws&sni=e9464f45.trauma-2r4.pages.dev
	q := u.Query()
	Assert (t, q.Get("type"), "tcp")
	Assert (t, q.Get("headerType"), "none")
	Assert (t, q.Get("security"), "tls")
	Assert (t, q.Get("allowInsecure"), "true")
	Assert (t, q.Get("sni"), "x.com")
	Assert (t, q.Get("alpn"), "h2,http/1.1")
	Assert (t, q.Get("fp"), "firefox-66")
}

// Advanced TCP test
func TestGen_trojan_3 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "trojan" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss,
		`"network": "tcp", "security": "none_sec",
         "tcpSettings": {"header": {
             "type": "http",
             "request": {"path": "/test", "headers": { "Host": "test.com" }}
         }}`,
	)); nil != e {
		panic (e);
	}
	u := Gen_trojan_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "tcp");
	Assert (t, q.Get("security"), "none_sec");	
	Assert (t, q.Get("headerType"), "http");
	Assert (t, q.Get("path"), "/test");
	Assert (t, q.Get("host"), "test.com");
}

// Trojan over GRPC + Reality test
func TestGen_trojan_4 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "trojan" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss,
		`"network": "grpc", "security": "reality",
         "grpcSettings": {
             "authority": "abcd", "serviceName": "srv_name", "multiMode": true
         }, "realitySettings": {
             "serverNames": ["x.com", "y.com"],
             "publicKey": "public.key", "shortId": "shortID",
             "spiderX": "/test/spider", "fingerprint": "firefox-66"
         }`,
	)); nil != e {
		panic (e);
	}
	u := Gen_trojan_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "grpc");
	Assert (t, q.Get("security"), "reality");	
	Assert (t, q.Get("serviceName"), "srv_name");
	Assert (t, q.Get("mode"), "multi");
	Assert (t, q.Get("authority"), "abcd");
	
	Assert (t, q.Get("sni"), "x.com");
	Assert (t, q.Get("pbk"), "public.key");
	Assert (t, q.Get("sid"), "shortID");
	Assert (t, q.Get("spx"), "/test/spider");
	Assert (t, q.Get("fp"), "firefox-66");
}

// Websocket test
func TestGen_trojan_5 (t *testing.T) {
	cfg := &conf.OutboundDetourConfig{ Protocol: "trojan" }
	if e := unmarshal_H (cfg, fmt.Sprintf(FMT_ss,
		`"network": "ws", "security": "none_sec",
         "wsSettings": {
             "host": "x.com", "path": "/test/", "headers": {}
         }`,
	)); nil != e {
		panic (e);
	}
	u := Gen_trojan_URL (cfg);
	if nil == u {
		t.Fatal("failed")
	}

	q := u.Query()
	Assert (t, q.Get("type"), "ws");
	Assert (t, q.Get("security"), "none_sec");	
	Assert (t, q.Get("path"), "/test/");
	Assert (t, q.Get("host"), "x.com");
}
