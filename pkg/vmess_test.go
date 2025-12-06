package pkg

import (
	"testing"
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
