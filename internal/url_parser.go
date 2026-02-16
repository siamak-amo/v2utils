// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"errors"
	"strings"

	"net/url"
	"encoding/base64"

	log "github.com/siamak-amo/v2utils/log"
)

var (
	trojan_stream_parser    = __stream_parser_vless_trojan
    vless_stream_parser     = __stream_parser_vless_trojan
    vmess_stream_parser     = __stream_parser_vmess

    vless_security_parser   = __security_parser_vless_trojan
	trojan_security_parser  = __security_parser_vless_trojan
	vmess_security_parser   = __security_parser_vmess
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
		return parse_vmess_url (link)
	case "ss":
		return parse_ss_url (u)
	case "trojan":
		return parse_trojan_url (u), nil

	default:
		return nil, errors.New ("Invalid URL scheme")
	}
}

// 	url: "vless://uuid@address:port?key=val..."
func parse_vless_url (u *url.URL) (URLmap) {
	res := make (URLmap, 0)
	params := Str2Strr(u.Query())

	res[Protocol] = "vless"
	res[ServerPort] = u.Port()
	res[Vxess_ID] = u.User.Username()
	res[ServerAddress] = u.Hostname()
	res[Security] = params.Pop ("security")
	res[Vless_ENC] = params.Pop ("encryption")
	res[Vless_Flow] = params.Pop ("flow")
	res[Vless_Level] = params.Pop ("level")

	res[Network] = params.Pop ("type")
	vless_stream_parser (res, params);
	vless_security_parser (res, params);

	extract_unused ("vless", params);
	return res
}

// 	url:  "vmess://BASE64(Json(key: value, ...))"
func parse_vmess_url (input string) (URLmap, error) {
	if (len (input) <= 8) { // 'vmess://'
		return nil, errors.New ("Invalid URL")
	}
	decoded, e := base64.StdEncoding.DecodeString(input[8:])
	if nil != e {
		return nil, e
	}
	src := make(Str2Str, 0)
	if e = unmarshal_H (&src, string(decoded)); nil != e {
		return nil, e
	}

	res := make (URLmap, 0)
	res[Protocol] = "vmess"
	res[ServerAddress] = src.Pop ("add")
	res[ServerPort] = src.Pop ("port")
	res[Vxess_ID] = src.Pop ("id")
	res[Network] = src.Pop ("net")

	vmess_stream_parser (res, src);
	vmess_security_parser (res, src);

	src.Pop ("aid"); src.Pop ("scy"); src.Pop ("ps"); src.Pop ("v") // unused
	extract_unused ("vmess", src);
	return res, nil
}

// 	url:  "ss://BASE64(method:password)@address:port"
func parse_ss_url (u *url.URL) (URLmap, error) {
	res := make (URLmap, 0)

	decoded, e := base64.StdEncoding.DecodeString(u.User.Username())
	if nil != e {
		return nil, e
	}

	mp := strings.Split (string(decoded), ":")
	if len(mp) >= 1 {
		res[SS_Method] = mp[0];
	}
	if len(mp) >= 2 {
		res[SS_Password] = mp[1];
	}
	res[Protocol] = "shadowsocks"
	res[ServerPort] = u.Port()
	res[ServerAddress] = u.Hostname()

	extract_unused ("shadowsocks", u.Query());
	return res, nil
}

// 	url:  "trojan://password@address:port?key=value..."
func parse_trojan_url (u *url.URL) (URLmap) {
	res := make (URLmap, 0)
	params := Str2Strr(u.Query())

	res[Protocol] = "trojan"
	res[ServerPort] = u.Port()
	res[ServerAddress] = u.Hostname()
	res[Trojan_Password] = u.User.Username()

	res[Network] = params.Pop ("type")
	res[Security] = params.Pop ("security")

	trojan_stream_parser (res, params);
	trojan_security_parser (res, params);

	extract_unused ("trojan", params);
	return res
}


// Internal stream/security parser functions
// Only use:  xxx_security_parser and xxx_stream_parser functions

func __security_parser_vless_trojan (dst URLmap, src Str2Strr) {
	switch (dst[Security]) {
	case "tls":
		dst[TLS_fp] = src.Pop ("fp")
		dst[TLS_sni] = src.Pop ("sni")
		dst[TLS_ALPN] = src.Pop ("alpn")
		dst[TLS_AllowInsecure] = cbool(src.Pop ("allowInsecure"))
		break;

	case "reality":
		dst[REALITY_fp] = src.Pop ("fp")
		dst[REALITY_sni] = src.Pop ("sni")
		dst[REALITY_ShortID] = src.Pop ("sid")
		dst[REALITY_SpiderX] = src.Pop ("spx")
		dst[REALITY_PublicKey] = src.Pop ("pbk")
		break;

	default:
		break;
	}
}

func __security_parser_vmess (dst URLmap, src Str2Str) {
	if s := src.Pop ("tls"); "tls" == s {
		dst[Security] = "tls"
		dst[TLS_sni] = src.Pop ("sni")
		dst[TLS_fp] = src.Pop ("fp")
		dst[TLS_ALPN] = src.Pop ("alpn")
		dst[TLS_AllowInsecure] = cbool(src.Pop ("allowInsecure"))
	}
}

func __stream_parser_vless_trojan (dst URLmap, src Str2Strr) {
	switch (dst[Network]) {
	case "ws":
		// TODO: what about WS_Headers?
		dst[WS_Path] = src.Pop ("path")
		dst[WS_Host] = src.Pop ("host")
		break;

	case "tcp":
		dst[TCP_HTTP_Host] = src.Pop ("host")
		dst[TCP_HTTP_Path] = src.Pop ("path")
		dst[TCP_HeaderType] = src.Pop ("headerType")
		break;

	case "kcp", "mkcp":
		dst[KCP_SEED] = src.Pop ("path");
		dst[KCP_HType] = src.Pop ("headerType");
		break;

	case "grpc":
		dst[GRPC_Mode] = src.Pop ("mode")
		dst[GRPC_MultiMode] = src.Pop ("multiMode")
		dst[GRPC_ServiceName] = src.Pop ("serviceName")
		break;

	case "xhttp", "splithttp":
		dst[XHTTP_Path] = src.Pop ("path")
		dst[XHTTP_Mode] = src.Pop ("mode")
		dst[XHTTP_Host] = src.Pop ("host")
		break;

	case "httpupgrade":
		// TODO: HTTPUP_Headers
		dst[HTTPUP_Host] = src.Pop ("host")
		dst[HTTPUP_Path] = src.Pop ("path")
		break;

	default:
		break;
	}
}

func __stream_parser_vmess (dst URLmap, src Str2Str) {
	switch (dst[Network]) {
	case "tcp":
		dst[TCP_HTTP_Host] = src.Pop ("host")
		dst[TCP_HTTP_Path] = src.Pop ("path")
		dst[TCP_HeaderType] = src.Pop ("type")
		break;

	case "grpc":
		dst[GRPC_Mode] = src.Pop ("type")
		break;

	case "ws":
		// TODO: WS_Headers
		dst[WS_Path] = src.Pop ("path")
		dst[WS_Host] = src.Pop ("host")
		break;

	case "xhttp", "splithttp":
		// TODO: XHTTP_Headers
		dst[XHTTP_Host] = src.Pop ("host")
		dst[XHTTP_Path] = src.Pop ("path")
		dst[XHTTP_Mode] = src.Pop ("mode")
		break;

	case "httpupgrade":
		// TODO: HTTPUP_Headers
		dst[HTTPUP_Host] = src.Pop ("host")
		dst[HTTPUP_Path] = src.Pop ("path")
		break;

	default:
		break;
	}
}

func extract_unused (name string, params any) {
	switch x := params.(type) {
	case Str2Str:
		for key, v := range x {
			if 0 != len(v) {
				log.Warnf("%s parser - parameter '%v' was ignored.\n", name, key)
			}
		}
		break;
	case Str2Strr:
		for key, v := range x {
			if 0 != len(v) {
				log.Warnf("%s parser - parameter '%v' was ignored.\n", name, key);
			}
		}
		break;
	case url.Values:
		for key, _ := range x {
			if 0 != len(key) {
				log.Warnf("%s parser - parameter '%v' was ignored.\n", name, key);
			}
		}
		break;
	}
}
