// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"fmt"
	"errors"
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

func set_stream_grpc (args URLmap, dst *conf.StreamConfig) (error) {
	return unmarshal_H (&dst.GRPCSettings,
		fmt.Sprintf (`{"serviceName": "%s", "multiMode": %s, "mode": "%s"}`,
			args[GRPC_ServiceName], args[GRPC_MultiMode], args[GRPC_Mode]),
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
	case "grpc":
		map_normal (args, GRPC_MultiMode, "false")
		e = set_stream_grpc (args, dst)
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

func encode_tcp_header(src []byte) (TCPHeaderConfig, error) {
	v := TCPHeaderConfig{}
	if e := json.Unmarshal(src, &v); nil != e {
		return v,e
	}
	return v,nil
}
