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
	tc.Assert (tc.Output.Network,    tc.Input[Network])
	tc.Assert (tc.Output.Security,   tc.Input[Security])
	tc.Assert (tls.ServerName,       tc.Input[TLS_sni])
	tc.Assert (tls.Insecure,         tc.Input[TLS_AllowInsecure])


	// Test the default value of allowInsecure == True
	delete (tc.Input, TLS_AllowInsecure)
	v, e = Gen_streamSettings (tc.Input)
	if nil != e {
		t.Fatalf ("Gen_streamSettings failed: %v\n", e)
		return
	}

	tc.Do(v);
	tls = tc.Output.TLSSettings
	tc.Assert (tc.Output.Network,    tc.Input[Network])
	tc.Assert (tc.Output.Security,   tc.Input[Security])
	tc.Assert (tls.ServerName,       tc.Input[TLS_sni])
	tc.Assert (tls.Insecure,         "true")
}
