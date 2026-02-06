// SPDX-License-Identifier: GPL-3.0-or-later
package internal

import (
	"fmt"
	"errors"
	"strconv"
	"strings"

	"net/url"
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

func set_stream_tcp (args URLmap, dst *conf.StreamConfig) (error) {
	switch (args[TCP_HeaderType]) {
	case "none", "":
		return unmarshal_H (&dst.TCPSettings, (`{"header": {"type": "none"}}`));

	case "http":
		return unmarshal_H (&dst.TCPSettings,
			fmt.Sprintf (
				`{
                    "header": {
                        "type": "%s",
                        "request": {
                            "version": "1.1", "path": ["%s"], "headers": {"Host": "%s"}
                        }
                    }
                 }`,
				args[TCP_HeaderType], args[TCP_HTTP_Path], args[TCP_HTTP_Host],
			),
		);

	default:
		return not_implemented ("not implemented header type: " + args[TCP_HeaderType]);
	}
}

func set_stream_ws (args URLmap, dst *conf.StreamConfig) (error) {
	args[WS_Headers] = csv2jsonArray (args[WS_Headers]);
	return unmarshal_H (&dst.WSSettings,
		fmt.Sprintf (`{"path": "%s", "host": "%s", "headers": [%s]}`,
			args[WS_Path], args[WS_Host], args[WS_Headers]),
	);
}

func set_stream_kcp (args URLmap, dst *conf.StreamConfig) (error) {
	return unmarshal_H (&dst.KCPSettings,
		fmt.Sprintf (`{"seed": "%s", "header": {"type": "%s"}}`,
			args[KCP_SEED], args[KCP_HType]),
	);
}

func set_stream_grpc (args URLmap, dst *conf.StreamConfig) (error) {
	return unmarshal_H (&dst.GRPCSettings,
		fmt.Sprintf (`{"serviceName": "%s", "multiMode": %s, "mode": "%s"}`,
			args[GRPC_ServiceName], args[GRPC_MultiMode], args[GRPC_Mode]),
	);
}

func set_stream_xhttp (args URLmap, dst *conf.StreamConfig) (error) {
	args[XHTTP_Headers] = csv2jsonArray (args[XHTTP_Headers]);
	return unmarshal_H (&dst.SplitHTTPSettings,
		fmt.Sprintf (`{"host": "%s", "path": "%s", "mode": "%s", "headers": {%s}}`,
			args[XHTTP_Host], args[XHTTP_Path], args[XHTTP_Mode], args[XHTTP_Headers],
		),
	);
}

func set_stream_httpupgrade (args URLmap, dst *conf.StreamConfig) (error) {
	args[HTTPUP_Headers] = csv2jsonArray (args[HTTPUP_Headers]);
	return unmarshal_H (&dst.HTTPUPGRADESettings,
		fmt.Sprintf (`{"host": "%s", "path": "%s", "headers": {%s}}`,
			args[HTTPUP_Host], args[HTTPUP_Path], args[HTTPUP_Headers],
		),
	);
}

func set_sec_tls (args URLmap, dst *conf.StreamConfig) (error) {
	args[TLS_ALPN] = csv2jsonArray (args[TLS_ALPN]);
	return unmarshal_H (&dst.TLSSettings,
		fmt.Sprintf (
			`{"servername": "%s", "allowInsecure": %s, "alpn": [%s], "fingerprint": "%s"}`,
			args[TLS_sni], args[TLS_AllowInsecure], args[TLS_ALPN], args[TLS_fp],
		),
	);
}

func set_sec_reality (args URLmap, dst *conf.StreamConfig) (error) {
	return unmarshal_H (&dst.REALITYSettings,
		fmt.Sprintf (
			`{"serverName": "%s", "fingerprint": "%s", "show": %s,
              "publicKey": "%s", "shortId": "%s", "spiderX": "%s"}`,
			args[REALITY_sni], args[REALITY_fp], args[REALITY_Show],
			args[REALITY_PublicKey], args[REALITY_ShortID], args[REALITY_SpiderX],
		),
	);
}

func set_stream_settings(args URLmap, dst *conf.StreamConfig) (e error) {
	switch (args[Network]) {
	case "ws":
		e = set_stream_ws (args, dst);
		break;
	case "tcp":
		e = set_stream_tcp (args, dst);
		break;
	case "kcp", "mkcp":
		e = set_stream_kcp (args, dst);
		break;
	case "grpc":
		map_normal (args, GRPC_MultiMode, "false")
		e = set_stream_grpc (args, dst)
		break;
	case "xhttp":
		e = set_stream_xhttp (args, dst)
		break;
	case "httpupgrade":
		e = set_stream_httpupgrade (args, dst)
		break;
	default:
		return not_implemented ("network " + args[Network])
	}
	if nil != e {
		// log
		return
	}

	switch (args[Security]) {
	case "", "none":
		break

	case "tls":
		map_normal (args, TLS_AllowInsecure, "true")
		map_normal (args, TLS_ALPN, "h2,http/1.1")
		e = set_sec_tls (args, dst)
		break;

	case "reality":
		map_normal (args, REALITY_Show, "false")
		e = set_sec_reality (args, dst);
		break;

	case "xtls":
		return not_implemented ("security " + args[Security]);

	default:
		return errors.New ("invalid security protocol: " + args[Security]);
	}
	return
}


func Gen_streamSettings(args URLmap) (dst *conf.StreamConfig, e error) {
	// Set the default network to tcp and security to none
	map_normal (args, Network, "tcp")
	map_normal (args, Security, "none")
	map_normal (args, TCP_HeaderType, "none")
	dst = &conf.StreamConfig{}
	if e = unmarshal_H (dst,
		fmt.Sprintf (`{"network": "%s", "security": "%s"}`,
			args[Network], args[Security],
		)); nil != e {
		// log
		return
	}
	if e = set_stream_settings (args, dst); nil != e {
		// log
		return
	}
	return
}



// Only for generating URLs //

// This type is not compatible with xray-core.
type TCPHeaderConfig struct {
	Type string						`json:"type"`
	Request struct {
		Path string					`json:"path"`
		Headers map[string]string	`json:"headers"`
	}
}
type KCPHeaderConfig struct {
	Type string					    `json:"type"`
}

func encode_tcp_header(src []byte) (TCPHeaderConfig, error) {
	v := TCPHeaderConfig{}
	if e := json.Unmarshal(src, &v); nil != e {
		return v,e
	}
	return v,nil
}

func encode_kcp_header(src []byte) (KCPHeaderConfig, error) {
	v := KCPHeaderConfig{}
	if e := json.Unmarshal(src, &v); nil != e {
		return v,e
	}
	return v,nil
}

type VlessURL_stream_handler  func(src *conf.StreamConfig, dst url.Values)
type VmessURL_stream_handler  func(src *conf.StreamConfig, dst map[string]string)
type SSURL_stream_handler     func(src *conf.StreamConfig, dst url.Values)
type TrojanURL_stream_handler func(src *conf.StreamConfig, dst url.Values)

var (
	Init_vlessURL_stream  VlessURL_stream_handler  = __set_kv_stream_vless_trojan
	Init_vmessURL_stream  VmessURL_stream_handler  = __set_kv_stream_vmess
	Init_ssURL_stream     SSURL_stream_handler     = __set_kv_stream_ss
	Init_trojanURL_stream TrojanURL_stream_handler = __set_kv_stream_vless_trojan
)

// vless / trojan compatible
func __set_kv_stream_vless_trojan(src *conf.StreamConfig, dst url.Values) {
	if nil != src {
		if nil != src.Network {
			net := string(*src.Network)
			AddQuery (dst, "type", net)
			switch (net) {
			case "tcp":
				if v,e := encode_tcp_header(src.TCPSettings.HeaderConfig); nil == e {
					AddQuery (dst, "headerType", v.Type)
					AddQuery (dst, "path", v.Request.Path)
					AddQuery (dst, "host", v.Request.Headers["Host"])
				}
				break;
			case "mkcp", "kcp":
				if nil != src.KCPSettings.Seed {
					AddQuery (dst, "path", *src.KCPSettings.Seed);
				}
				if v,e := encode_kcp_header(src.KCPSettings.HeaderConfig); nil == e {
					AddQuery (dst, "headerType", v.Type)
				}
				break;
			case "grpc":
				AddQuery (dst, "serviceName", src.GRPCSettings.ServiceName)
				AddQuery (dst, "authority", src.GRPCSettings.Authority)
				if src.GRPCSettings.MultiMode {
					AddQuery (dst, "mode", "multi")
				}
				break;
			case "ws":
				AddQuery (dst, "host", src.WSSettings.Host)
				AddQuery (dst, "path", src.WSSettings.Path)
				break;
			case "httpupgrade":
				AddQuery (dst, "host", src.HTTPUPGRADESettings.Host);
				AddQuery (dst, "path", src.HTTPUPGRADESettings.Path);
				break;
			case "xhttp":
				var _cfg *conf.SplitHTTPConfig;
				if nil != src.XHTTPSettings {
					_cfg = src.XHTTPSettings;
				} else if nil != src.SplitHTTPSettings {
					_cfg = src.SplitHTTPSettings;
				} else {
					break; // No xhttp settings is provided.
				}
				AddQuery (dst, "host", _cfg.Host)
				AddQuery (dst, "mode", _cfg.Mode)
				AddQuery (dst, "path", _cfg.Path)
				break;
			}
		} else {
			AddQuery (dst, "type", "tcp")
		}

		sec := src.Security
		AddQuery (dst, "security", sec)
		switch (sec) {
		case "tls":
			AddQuery (dst, "allowInsecure", strconv.FormatBool(src.TLSSettings.Insecure))
			AddQuery (dst, "sni", src.TLSSettings.ServerName)
			AddQuery (dst, "fp", src.TLSSettings.Fingerprint)
			if nil != src.TLSSettings.ALPN {
				AddQuery (dst, "alpn", strings.Join(*src.TLSSettings.ALPN, ","))
			}
			break;
		case "reality":
			AddQuery (dst, "fp", src.REALITYSettings.Fingerprint)
			AddQuery (dst, "spx", src.REALITYSettings.SpiderX)
			AddQuery (dst, "pbk", src.REALITYSettings.PublicKey)
			AddQuery (dst, "sid", src.REALITYSettings.ShortId)
			AddQuery (dst, "sni", src.REALITYSettings.ServerNames[0])
			AddQuery (dst, "mode", src.REALITYSettings.Type)
			break;
		}
	}
}

// vmess compatible
func __set_kv_stream_vmess(src *conf.StreamConfig, dst map[string]string) {
	if nil != src && nil != src.Network {
		net := string(*src.Network);
		dst["net"] = net
		switch (net) {
		case "tcp":
			if v,e := encode_tcp_header(src.TCPSettings.HeaderConfig); nil == e {
				dst["type"] = v.Type
				dst["path"] = v.Request.Path
				dst["host"] = v.Request.Headers["Host"]
			}
			break;
		case "grpc":
			dst["path"] = src.GRPCSettings.ServiceName
			dst["authority"] = src.GRPCSettings.Authority
			if src.GRPCSettings.MultiMode {
				dst["type"] = "multi"
			}
			break;

		case "ws":
			dst["host"] = src.WSSettings.Host
			dst["path"] = src.WSSettings.Path
			dst["type"] = "none"
			break;
		}
		if sec := src.Security; "tls" == sec {
			dst["tls"] = "tls";
			dst["fp"] = src.TLSSettings.Fingerprint
			dst["allowInsecure"] = strconv.FormatBool(src.TLSSettings.Insecure)
			dst["sni"] = src.TLSSettings.ServerName
			if nil != src.TLSSettings.ALPN {
				dst["alpn"] = strings.Join(*src.TLSSettings.ALPN, ",")
			}
		}
	} else {
		dst["net"] = "tcp";
	}
}

func __set_kv_stream_ss(src *conf.StreamConfig, dst url.Values) {
	// ??
}
