// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"errors"
	"strings"

	"net/url"
	"encoding/base64"

	log "github.com/siamak-amo/v2utils/log"
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

	res[Network] = params.Pop ("type")
	switch (res[Network]) {
	case "ws":
		// TODO: what about WS_Headers?
		res[WS_Path] = params.Pop ("path")
		res[WS_Host] = params.Pop ("host")
		break;

	case "tcp":
		res[TCP_HTTP_Host] = params.Pop ("host")
		res[TCP_HTTP_Path] = params.Pop ("path")
		res[TCP_HeaderType] = params.Pop ("headerType")
		break;

	case "kcp", "mkcp":
		res[KCP_SEED] = params.Pop ("path");
		res[KCP_HType] = params.Pop ("headerType");
		break;

	case "grpc":
		res[GRPC_Mode] = params.Pop ("mode")
		res[GRPC_MultiMode] = params.Pop ("multiMode")
		res[GRPC_ServiceName] = params.Pop ("serviceName")
		break;

	case "xhttp", "splithttp":
		res[XHTTP_Path] = params.Pop ("path")
		res[XHTTP_Mode] = params.Pop ("mode")
		res[XHTTP_Host] = params.Pop ("host")
		break;

	case "httpupgrade":
		// TODO: HTTPUP_Headers
		res[HTTPUP_Host] = params.Pop ("host")
		res[HTTPUP_Path] = params.Pop ("path")
		break;

	default:
		break;
	}

	switch (res[Security]) {
	case "tls":
		res[TLS_fp] = params.Pop ("fp")
		res[TLS_sni] = params.Pop ("sni")
		res[TLS_ALPN] = params.Pop ("alpn")
		res[TLS_AllowInsecure] = params.Pop ("allowInsecure")
		break;

	case "reality":
		res[REALITY_fp] = params.Pop ("fp")
		res[REALITY_sni] = params.Pop ("sni")
		res[REALITY_ShortID] = params.Pop ("sid")
		res[REALITY_SpiderX] = params.Pop ("spx")
		res[REALITY_PublicKey] = params.Pop ("pbk")
		break;

	default:
		break;
	}

	for key, v := range params {
		if len(v) >= 1 {
			log.Warnf("parse_vless_url - parameter '%v' was ignored.\n", key)
		}
	}
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
	dst := make(Str2Str, 0)
	if e = unmarshal_H (&dst, string(decoded)); nil != e {
		return nil, e
	}

	res := make (URLmap, 0)
	res[Protocol] = "vmess"
	res[ServerAddress] = dst.Pop ("add")
	res[ServerPort] = dst.Pop ("port")
	res[Vxess_ID] = dst.Pop ("id")
	res[Network] = dst.Pop ("net")

	switch (res[Network]) {
	case "tcp":
		res[TCP_HTTP_Host] = dst.Pop ("host")
		res[TCP_HTTP_Path] = dst.Pop ("path")
		res[TCP_HeaderType] = dst.Pop ("type")
		break;

	case "grpc":
		res[GRPC_Mode] = dst.Pop ("type")
		break;

	case "ws":
		// TODO: WS_Headers
		res[WS_Path] = dst.Pop ("path")
		res[WS_Host] = dst.Pop ("host")
		break;

	case "xhttp", "splithttp":
		// TODO: XHTTP_Headers
		res[XHTTP_Host] = dst.Pop ("host")
		res[XHTTP_Path] = dst.Pop ("path")
		res[XHTTP_Mode] = dst.Pop ("mode")
		break;

	case "httpupgrade":
		// TODO: HTTPUP_Headers
		res[HTTPUP_Host] = dst.Pop ("host")
		res[HTTPUP_Path] = dst.Pop ("path")
		break;

	default:
		break;
	}

	if s := dst.Pop ("tls"); "tls" == s {
		res[Security] = "tls"
		res[TLS_sni] = dst.Pop ("sni")
		res[TLS_fp] = dst.Pop ("fp")
		res[TLS_ALPN] = dst.Pop ("alpn")
		res[TLS_AllowInsecure] = dst.Pop ("allowInsecure")
	}

	dst.Pop ("aid"); dst.Pop ("scy"); dst.Pop ("ps"); dst.Pop ("v") // unused
	for key, v := range dst {
		if len(v) >= 1 {
			log.Warnf("parse_vmess_url - parameter '%v' was ignored.\n", key)
		}
	}
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
	switch (res[Network]) {
	case "ws":
		// TODO: WS_Headers
		res[WS_Path] = params.Pop ("path");
		res[WS_Host] = params.Pop ("host");
		break;
	case "tcp":
		res[TCP_HTTP_Host] = params.Pop ("host");
		res[TCP_HTTP_Path] = params.Pop ("path");
		res[TCP_HeaderType] = params.Pop ("headerType");
		break;
	case "grpc":
		res[GRPC_Mode] = params.Pop ("mode")
		res[GRPC_MultiMode] = params.Pop ("multiMode")
		res[GRPC_ServiceName] = params.Pop ("serviceName")
		break;
	case "kcp", "mkcp":
		res[KCP_SEED] = params.Pop ("path");
		res[KCP_HType] = params.Pop ("headerType");
		break;
	case "xhttp", "splithttp":
		res[XHTTP_Path] = params.Pop ("path")
		res[XHTTP_Mode] = params.Pop ("mode")
		res[XHTTP_Host] = params.Pop ("host")
		break;
	case "httpupgrade":
		// TODO: HTTPUP_Headers
		res[HTTPUP_Host] = params.Pop ("host")
		res[HTTPUP_Path] = params.Pop ("path")
		break;
	default:
		break;
	}

	switch (res[Security]) {
	case "tls":
		res[TLS_fp] = params.Pop ("fp")
		res[TLS_sni] = params.Pop ("sni")
		res[TLS_ALPN] = params.Pop ("alpn")
		res[TLS_AllowInsecure] = params.Pop ("allowInsecure")
		break;

	case "reality":
		res[REALITY_fp] = params.Pop ("fp")
		res[REALITY_sni] = params.Pop ("sni")
		res[REALITY_ShortID] = params.Pop ("sid")
		res[REALITY_SpiderX] = params.Pop ("spx")
		res[REALITY_PublicKey] = params.Pop ("pbk")
		break;

	default:
		break;
	}

	for key, v := range params {
		if len(v) >= 1 {
			log.Warnf("parse_ss_url - parameter '%v' was ignored.\n", key)
		}
	}
	return res
}
