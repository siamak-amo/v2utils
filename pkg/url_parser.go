package pkg

import (
	"fmt"
	"net/url"
)

func ParseURL(link string) (URLmap, error) {
	u, e := url.Parse(link)
	if nil != e {
		return nil, e
	}

	switch (u.Scheme) {
	case "vless":
		return parse_vless_url (u), nil
	case "vmess":
		return nil, nil
	default:
		return nil, not_implemented (u.Scheme)
	}
}

func Pop(m map[string][]string, key string) (string) {
	if v, ok := m[key]; ok {
		// delete (m, key)
		m[key] = []string{}
		if len(v) >= 1 {
			return v[0]
		} else {
			return ""
		}
	}
	return ""
}

// 	link := fmt.Sprintf("vless://%s@%s:%d", uuid, address, port)
func parse_vless_url (u *url.URL) (URLmap) {
	res := make (URLmap, 0)
	params := u.Query ()

	res[Protocol] = "vless"
	res[ServerPort] = u.Port()
	res[Vxess_ID] = u.User.Username()
	res[ServerAddress] = u.Hostname()
	res[Security] = Pop(params, "security")
	res[Vless_ENC] = Pop(params, "encryption")

	res[Network] = Pop(params, "type")
	switch (res[Network]) {
	case "ws":
		res[WS_Path] = Pop(params, "path")
		res[WS_Headers] = Pop(params, "host")
		break;

	case "tcp":
		res[TCP_HTTP_Host] = Pop(params, "host")
		res[TCP_HTTP_Path] = Pop(params, "path")
		res[TCP_HeaderType] = Pop(params, "headerType")
		break;

	default:
		break;
	}

	switch (res[Security]) {
	case "tls":
		res[TLS_fp] = Pop(params, "fp")
		res[TLS_sni] = Pop(params, "sni")
		res[TLS_ALPN] = Pop(params, "alpn")
		break;

	default:
		break;
	}

	for key, v := range params {
		if len(v) >= 1 {
			fmt.Printf ("parse_vless_url: parameter '%v' was ignored.\n", key)
		}
	}
	return res
}
