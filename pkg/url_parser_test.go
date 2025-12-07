package pkg

import (
	"testing"
)


func Test_parse_vless_url_1 (t *testing.T) {
	const VLESS_TEST_1 = "vless://3eae724f-9256@222.22.222.2:80?path=%2FTelegram%3A%40UnlimitedDev%2Fwww&security=none&encryption=none&host=vpn.com&type=ws#FreeInternet4You"
	umap, e := ParseURL(VLESS_TEST_1);
	if nil != e {
		t.Fatalf ("parse_vless_url failed: %v\n", e)
	}

	umap.Assert (t, Protocol,             "vless")
	umap.Assert (t, ServerAddress,        "222.22.222.2")
	umap.Assert (t, ServerPort,		      "80")
	umap.Assert (t, Vxess_ID,		      "3eae724f-9256")
	umap.Assert (t, Vless_ENC,		      "none")
	
	umap.Assert (t, Network,		      "ws")
	umap.Assert (t, Security,		      "none")
	umap.Assert (t, WS_Headers,		      "vpn.com")
	umap.Assert (t, WS_Path,		      "/Telegram:@UnlimitedDev/www")
}

func Test_parse_vless_url_2 (t *testing.T) {
	const VLESS_TEST_2 = "vless://884f9b2c-f0c5@boys.dev:443?path=%2F%3Fed%3D2048&security=tls&encryption=none&host=www.boys.dev&type=ws&sni=cf-wkrs-pages-vless-boys.dev#%40v2rayng_org"
	umap, e := ParseURL(VLESS_TEST_2);
	if nil != e {
		t.Fatalf ("parse_vless_url failed: %v\n", e)
	}

	umap.Assert (t, Protocol,             "vless")
	umap.Assert (t, ServerAddress,        "boys.dev")
	umap.Assert (t, ServerPort,		      "443")
	umap.Assert (t, Vxess_ID,		      "884f9b2c-f0c5")
	umap.Assert (t, Vless_ENC,		      "none")
	
	umap.Assert (t, Security,		      "tls")
	umap.Assert (t, TLS_sni,		      "cf-wkrs-pages-vless-boys.dev")

	umap.Assert (t, Network,		      "ws")
	umap.Assert (t, WS_Path,		      "/?ed=2048")
}

func Test_parse_vless_url_3 (t *testing.T) {
	const VLESS_TEST_3 = "vless://6378d738-1ed3@1.2.3.4:666?security=none&encryption=none2&host=varzesh3.com&headerType=http&type=tcp#%40v2rayng_org"
	umap, e := ParseURL(VLESS_TEST_3);
	if nil != e {
		t.Fatalf ("parse_vless_url failed: %v\n", e)
	}

	umap.Assert (t, Protocol,             "vless")
	umap.Assert (t, ServerAddress,        "1.2.3.4")
	umap.Assert (t, ServerPort,		      "666")
	umap.Assert (t, Vxess_ID,		      "6378d738-1ed3")
	umap.Assert (t, Vless_ENC,		      "none2")
	umap.Assert (t, Security,		      "none")
	umap.Assert (t, Network,		      "tcp")
}

func Test_parse_vless_url_4 (t *testing.T) {
	const VLESS_TEST_4 = "vless://0d0245cd-8412@usa-join.outline-vpn.fun:444?mode=gun&security=tls&encryption=test_algorithm&alpn=h2,http/1.1&fp=chrome&type=grpc&serviceName=@Gigiv2ray&sni=RadioTeHran.oRG#Test4"
	umap, e := ParseURL(VLESS_TEST_4);
	if nil != e {
		t.Fatalf ("parse_vless_url failed: %v\n", e)
	}

	umap.Assert (t, Protocol,			"vless")
	umap.Assert (t, ServerAddress,		"usa-join.outline-vpn.fun")
	umap.Assert (t, ServerPort,			"444")
	umap.Assert (t, Vless_ENC,			"test_algorithm")

	umap.Assert (t, Security,			"tls")
	umap.Assert (t, TLS_sni,			"RadioTeHran.oRG")
	umap.Assert (t, TLS_fp,				"chrome")
	umap.Assert (t, TLS_ALPN,			"h2,http/1.1")

	umap.Assert (t, Network,			"grpc")
	umap.Assert (t, GRPC_Mode,			"gun")
	umap.Assert (t, GRPC_ServiceName,	"@Gigiv2ray")
}

func Test_parse_vless_url_5 (t *testing.T) {
	const VLESS_TEST_4 = "vless://eb3fdd8f-bcc9@config4vpn.fun:2053?mode=gun&security=reality&encryption=none&pbk=z1tAnqd5RA4I99LrK5FCJgjCd&fp=chrome&spx=%2Fvpn4You%2Ftest&type=grpc&serviceName=@configforVPN&sni=discordapp.com&sid=b1803d25#%40ConfigforVPN"
	umap, e := ParseURL(VLESS_TEST_4);
	if nil != e {
		t.Fatalf ("parse_vless_url failed: %v\n", e)
	}

	umap.Assert (t, Network,			"grpc")
	umap.Assert (t, GRPC_Mode,			"gun")
	umap.Assert (t, GRPC_ServiceName,	"@configforVPN")

	umap.Assert (t, Security,			"reality")
	umap.Assert (t, REALITY_fp,			"chrome")
	umap.Assert (t, REALITY_sni,		"discordapp.com")
	umap.Assert (t, REALITY_ShortID,	"b1803d25")
	umap.Assert (t, REALITY_SpiderX,	"/vpn4You/test")
	umap.Assert (t, REALITY_PublicKey,	"z1tAnqd5RA4I99LrK5FCJgjCd")
}
