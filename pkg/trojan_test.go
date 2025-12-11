// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"testing"
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
