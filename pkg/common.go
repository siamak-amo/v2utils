// SPDX-License-Identifier: GPL-3.0-or-later
package pkg

import (
	"errors"
	"strings"
	"encoding/json"
)

type URLMapper int
type URLmap map[URLMapper]string

type Str2Str map[string]string
type Str2Strr map[string][]string

const (
	// Common parts
	ServerAddress URLMapper = iota
	ServerPort
	Protocol
	Network
	Security
	// Security: TLS
	TLS_sni
	TLS_ALPN
	TLS_fp // fingerprint
	TLS_AllowInsecure
	// Security: reality
	REALITY_fp // fingerprint
	REALITY_sni
	REALITY_Show
	REALITY_PublicKey
	REALITY_ShortID
	REALITY_SpiderX
	// Stream type specific parts
	TCP_HeaderType         // tcp
	TCP_HTTP_Host
	TCP_HTTP_Path
	WS_Path                // web socket
	WS_Headers
	GRPC_Mode               // GRPC
	GRPC_MultiMode
	GRPC_ServiceName
	// Protocol parts
	Vxess_ID  // vless & vmess  we call them vxess
	Vless_ENC
	Vless_Level
	Vmess_Sec
	Vmess_AlterID
	SS_Password
	SS_Method
	Trojan_Password
)

func unmarshal_H (t interface{}, input string) (error) {
	err := json.Unmarshal ([]byte(input), t)
	if err != nil {
		return err
	}
	return nil
}

func map_normal (m URLmap, key URLMapper, def_val string) (string){
	val, ok := m[key]
	if !ok || "" == val {
		m[key] = def_val
		return def_val
	} else {
		return val
	}
}

func not_implemented (feature string) error {
    return errors.New(feature + " not implemented")
}

// converts: `x,y,z` -> `"x", "y", "z"`
func csv2jsonArray (csv string) string {
	var res string
	for _, key := range strings.Split(csv, ",") {
		res += `"` + key + `",`
	}
	if len(res) >= 1 {
		res = res[:len(res)-1]
	}
	return res
}

func (m Str2Strr) Pop(key string) (string) {
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

func (m Str2Str) Pop(key string) (string) {
	if v, ok := m[key]; ok {
		// delete (m, key)
		m[key] = ""
		if len(v) >= 1 {
			return v
		} else {
			return ""
		}
	}
	return ""
}
