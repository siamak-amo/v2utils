// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"fmt"
	"time"
	"errors"
	"net/http"

	pkg "github.com/siamak-amo/v2utils/pkg"
	log "github.com/siamak-amo/v2utils/log"
	"golang.org/x/net/proxy"
)

const (
    HTTP_Test_Endpoint = "http://google.com"
    FMT_DefaultTemplate = `{
            "log": {"loglevel": "none"},
            "inbounds": [{
                "listen": "127.0.0.1", "port": %d,
                "protocol": "socks"
            }]
         }`
)

// Generates a minimal default config to test
// Returns the listening port number (and -1 on failure)
func (opt *Opt) MK_Default_TestConfig(url string) int {
	port := 2054 // TODO: make it random
	cf, e := pkg.Gen_main(fmt.Sprintf(FMT_DefaultTemplate, port))
	if nil != e || 0 == len(cf.InboundConfigs){
		panic (errors.New("Making test config failed."))
	}
	opt.CFG = *cf;

	if e := opt.Init_Outbound_byURL(url); nil != e {
		return -1
	}
	return port
}

// Tests the socks proxy of the Default_CFG
func (opt Opt) test_Proxy(port int) bool {
	if e := opt.Run_Xray(); nil != e {
		log.Errorf("Could not run the proxy server - %v\n", e)
		return false
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		return false
	}
	// time.Sleep(50*time.Second)
	res := test_http(dialer);
	opt.Kill_Xray()
	return (res == nil)
}

func test_http(dialer proxy.Dialer) error {
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   5 * time.Second,
	}

	// TODO: Is this a reliable way to check dialer works?
	_, err := client.Get (HTTP_Test_Endpoint)
	if err != nil {
		log.Infof("Broken VPN :( %v\n", err)
		return err
	}
	return nil
}

func (opt *Opt) Test_URL(url string) bool {
	if port := opt.MK_Default_TestConfig(url); port > 0 {
		return opt.test_Proxy(port);
	} else {
		return false
	}
}
