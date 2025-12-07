package pkg

import (
	"fmt"
	"runtime"
	"testing"
	"strconv"
	"encoding/json"

	core "github.com/xtls/xray-core/transport/internet"
)

type TestCase[T any] struct {
	*testing.T
	Input       URLmap
	Output      T
}

func (tc *TestCase[T]) Do (data any) {
	jsonData, err := json.Marshal (data)
	if nil != err {
		tc.Fatalf ("TestCase.Do failed: %v\n", err)
		return
	}
	if e := unmarshal_H (&tc.Output, string(jsonData)); nil != e {
		tc.Fatalf ("unmarshal_H failed:  %v\n", e)
		return
	}
}

// Call it with ONLY ONE function call depth
// this happens if your language doesn't support macro
func printFailure(actual, expected interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		fmt.Printf ("%s:%d:  test failed.\n\t", file, line)
	} else {
		fmt.Printf ("test failed.\n\t")
	}
	fmt.Printf ("(expected '%v')  !=  (actual '%v').\n", expected, actual)
}

func (umap URLmap) Assert (t *testing.T, id URLMapper, expected string) {
	v, ok := umap[id];
	if !ok || v != expected {
		printFailure (v, expected)
		t.Fail()
	}
}

func (tc TestCase[T]) Assert (val interface{}, expected string) {
	var value string
	if i, ok := val.(int); ok {
		value = strconv.Itoa (i)
	} else if s, ok := val.(string); ok {
		value = s
	} else if b, ok := val.(bool); ok {
		if b { value = "true" } else { value = "false" }
	} else {
		tc.Fatalf ("Assertion failed: type not implemented.\n")
		tc.Fail()
	}

	if value != expected {
		printFailure (value, expected)
		tc.Fail()
	}
}


// outbound, settings can be: any of xxxOutboundConfig types
type OutboundDetourConfig[SETTINGS any] struct {
	Protocol       string           `json:"protocol"`
	Tag            string           `json:"tag"`
	Settings       SETTINGS         `json:"settings"`
	StreamSetting  *core.StreamConfig `json:"streamSettings"`
}

// vnext for vless/vmess, account can be any of xxxAccount types
type VXessOutboundConfig[T any] struct {
	Vnext []T						`json:"vnext"`
}
type VXessOutboundVnext[ACCOUNT any] struct {
	Address string					`json:"address"`
	Port    int						`json:"port"`
	Users   []ACCOUNT				`json:"users"`
}

type VMessAccount struct {
	ID          string				`json:"id"`
	AlterIds    int					`json:"alterId"`
	Security    string				`json:"security"`
}
type VLessAccount struct {
	ID          string				`json:"id"`
	Encryption  string				`json:"encryption"`
	Level       int					`json:"level"`
}

// Complete types to be used in testings
type VLessCFG OutboundDetourConfig[VXessOutboundConfig[VXessOutboundVnext[VLessAccount]]]
type VMessCFG OutboundDetourConfig[VXessOutboundConfig[VXessOutboundVnext[VMessAccount]]]

type StreamConfig struct {
	Network             string               `json:"network"`
	Security            string               `json:"security"`
	TCPSettings         *TCPConfig           `json:"tcpSettings"`
	TLSSettings         *TLSConfig           `json:"tlsSettings"`
	REALITYSettings     *REALITYConfig       `json:"realitySettings"`
	WSSettings          *WebSocketConfig     `json:"wsSettings"`
	HTTPSettings        *HTTPConfig          `json:"httpSettings"`
	Hy2Settings         *Hy2Config           `json:"hy2Settings"`
	QUICSettings        *QUICConfig          `json:"quicSettings"`
	GRPCConfig          *GunConfig           `json:"grpcSettings"`
}

type TCPConfig struct {
	HeaderConfig map[string]string  `json:"header"`
}
type TLSConfig struct {
	Insecure    bool				`json:"allowInsecure"`
	ServerName  string				`json:"serverName"`
	ALPN        []string			`json:"alpn"`
}
type REALITYConfig struct {
	Show          bool          `json:"show"`
	ServerName    string	    `json:"serverName"`
	ServerNames   []string      `json:"serverNames"`
	PrivateKey    string        `json:"privateKey"`
	ShortId       string	    `json:"shortId"`
	ShortIds      []string      `json:"shortIds"`
	Fingerprint   string	    `json:"fingerprint"`
	PublicKey     string	    `json:"publicKey"`
	SpiderX       string	    `json:"spiderX"`
}
type WebSocketConfig struct {
	Path       string				`json:"path"`
	Headers    map[string]string	`json:"headers"`
}
type HTTPConfig struct {
	Host    []string				`json:"host"`
	Path    string					`json:"path"`
	Method  string					`json:"method"`
	Headers map[string][]string		`json:"headers"`
}
type Hy2Config struct {
	Password         string			`json:"password"`
	UseUDPExtension  bool			`json:"use_udp_extension"`
}
type QUICConfig struct {
	Header   map[string]string		`json:"header"`
	Security string					`json:"security"`
	Key      string					`json:"key"`
}
type GunConfig struct {
	ServiceName string				`json:"serviceName"`
}
