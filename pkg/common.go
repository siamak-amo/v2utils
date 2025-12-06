package pkg

import (
	"fmt"
	"errors"
	"encoding/json"
)

type URLMapper int
type URLmap map[URLMapper]string

const (
	// Common parts
	ServerAddress URLMapper = iota
	ServerPort
	Protocol
	Network
	Security
	// Security config
	TLS_sni
	TLS_ALPN
	TLS_fp // fingerprint
	TLS_AllowInsecure
	// Stream type specific parts
	TCP_HeaderType         // tcp
	TCP_HTTP_Host
	TCP_HTTP_Path
	WS_Path                // web socket
	WS_Headers
	// Protocol parts
	Vxess_ID  // vless & vmess  we call them vxess
	Vless_ENC
	Vless_Level
	Vmess_Sec
	Vmess_AlterID
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
    return errors.New(fmt.Sprintf ("%s - not implemented", feature))
}
