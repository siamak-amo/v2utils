// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"testing"
)

// Test empty keys
func Test_Gen_StreamSettings_1 (t *testing.T) {
	tc := TestCase[StreamConfig] {T: t,
		Input: map[URLMapper]string {},
		Output: StreamConfig{},
	}

	v, e := Gen_streamSettings (tc.Input)
	if nil != e {
		t.Fatalf ("Gen_streamSettings failed: %v\n", e)
		return
	}
	tc.Do(v);

	tc.Assert (tc.Output.Network,  "tcp")
	tc.Assert (tc.Output.Security, "none")
}


// Test TLS
func Test_Gen_StreamSettings_2 (t *testing.T) {
	tc := TestCase[StreamConfig] {T: t,
		Input: map[URLMapper]string {
	    	Network:    		"tcp",
			Security:			"tls",
			TLS_sni:			"x.com",
			TLS_AllowInsecure:	"false",
			TLS_ALPN:           "p1,p2,p3",
	    },
		Output: StreamConfig{},
	}

	v, e := Gen_streamSettings (tc.Input)
	if nil != e {
		t.Fatalf ("Gen_streamSettings failed: %v\n", e)
		return
	}

	tc.Do(v);
	tls := tc.Output.TLSSettings
	if nil == tls {
		t.Fatalf ("TLSSettings is uninitialized")
	}

	tc.Assert (tc.Output.Network,    tc.Input[Network])
	tc.Assert (tc.Output.Security,   tc.Input[Security])
	tc.Assert (tls.ServerName,       tc.Input[TLS_sni])
	tc.Assert (tls.Insecure,         tc.Input[TLS_AllowInsecure])

	tc.Assert (tls.ALPN[0], "p1")
	tc.Assert (tls.ALPN[1], "p2")
	tc.Assert (tls.ALPN[2], "p3")


	// Test the default values of allowInsecure and ALPN
	delete (tc.Input, TLS_AllowInsecure)
	delete (tc.Input, TLS_ALPN)
	v, e = Gen_streamSettings (tc.Input)
	if nil != e {
		t.Fatalf ("Gen_streamSettings failed: %v\n", e)
		return
	}

	tc.Do(v);
	tls = tc.Output.TLSSettings
	if nil == tls {
		t.Fatalf ("TLSSettings is uninitialized")
	}
	tc.Assert (tc.Output.Network,    tc.Input[Network])
	tc.Assert (tc.Output.Security,   tc.Input[Security])
	tc.Assert (tls.ServerName,       tc.Input[TLS_sni])
	tc.Assert (tls.Insecure,         "true")
}

// Reality test
func Test_Gen_StreamSettings_3 (t *testing.T) {
	tc := TestCase[StreamConfig] {T: t,
		Input: map[URLMapper]string {
		    Network:    		"grpc",
			Security:			"reality",
			REALITY_fp:	        "firefox-66",
			REALITY_sni:	    "sec.vpn.net",
			REALITY_SpiderX:    "%2Fvideos%2F",
			REALITY_PublicKey:  "OBR2JYROQB8odK5glVW_KLnsWl3UZ",
		},
		Output: StreamConfig{},
	}
	v, e := Gen_streamSettings (tc.Input)
	if nil != e {
		t.Fatalf ("Gen_streamSettings failed: %v\n", e)
		return
	}

	tc.Do(v);
	reality := tc.Output.REALITYSettings
	if nil == reality {
		t.Fatalf ("RealitySettings is uninitialized")
	}

	tc.Assert (tc.Output.Network,           tc.Input[Network])
	tc.Assert (tc.Output.Security,          tc.Input[Security])
	tc.Assert (reality.Show,                "false");
	tc.Assert (reality.Fingerprint,         tc.Input[REALITY_fp]);
	tc.Assert (reality.ServerName,          tc.Input[REALITY_sni]);
	tc.Assert (reality.SpiderX,             tc.Input[REALITY_SpiderX]);
	tc.Assert (reality.PublicKey,           tc.Input[REALITY_PublicKey]);
}
