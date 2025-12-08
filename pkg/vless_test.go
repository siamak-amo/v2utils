// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"testing"
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
