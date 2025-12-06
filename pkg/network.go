package pkg

import (
	"fmt"
	"strings"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
)

func set_tls_alpn (args URLmap) {
	var res string
	if _, ok := args[TLS_ALPN]; !ok {
		args[TLS_ALPN] = `"h2", "http/1.1"`
		return
	}
	for _, key := range strings.Split(args[TLS_ALPN], ",") {
		res += `"` + key + `",`
	}
	if len(res) >= 1 {
		res = res[:len(res)-1]
	}
	args[TLS_ALPN] = res
}

func set_stream_settings(args URLmap, dst *v4.StreamConfig) (e error) {
	switch (args[Network]) {
	case "ws":
		e = unmarshal_H (&dst.WSSettings,
			fmt.Sprintf (`{"path": "%s", "headers": {"Host": "%s"}}`,
				args[WS_Path], args[WS_Headers]),
		);
		break

	case "tcp":
		switch (args[TCP_HeaderType]) {
		case "none", "":
			e = unmarshal_H (&dst.TCPSettings, (`{"header": {"type": "none"}}`));
			break;

		case "http":
			e = unmarshal_H (&dst.TCPSettings, fmt.Sprintf (
				`{
                    "header": {"type": "%s"}, "request": {
                      "version": "1.1", "path": ["%s"], "headers": {"Host": "%s"}
                    }
                 }`,
				args[TCP_HeaderType], args[TCP_HTTP_Path], args[TCP_HTTP_Host]),
			);
			break;

		default:
			break;
		} // End of `switch (args[TCP_HeaderType])`
		break

	default:
		return not_implemented ("network " + args[Network])
	}
	if nil != e {
		// log
	}

	switch (args[Security]) {
	case "", "none":
		break

	case "tls":
		map_normal (args, TLS_AllowInsecure, "true")
		set_tls_alpn (args)
		if e = unmarshal_H (&dst.TLSSettings,
			fmt.Sprintf (
				`{"servername": "%s", "allowInsecure": %s, "alpn": [%s], "fingerprint": "%s"}`,
				args[TLS_sni], args[TLS_AllowInsecure], args[TLS_ALPN], args[TLS_fp],
			),
		); e != nil {
			// log
		}
		break;
	case "reality":
		return not_implemented ("reality")
	default:
		return not_implemented ("security " + args[Security])
	}
	return
}


func Gen_streamSettings(args URLmap) (dst *v4.StreamConfig, e error) {
	// Set the default network to tcp and security to none
	map_normal (args, Network, "tcp")
	map_normal (args, Security, "none")
	map_normal (args, TCP_HeaderType, "none")
	dst = &v4.StreamConfig{}
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
