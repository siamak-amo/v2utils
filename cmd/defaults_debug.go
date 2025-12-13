// defaults_debug.go
//go:build debug
// +build debug
package main

const (
	// Default template for run command
	DEF_Run_Template =`
         {
              "log": {"loglevel": "info"},
              "inbounds": [
                  {"listen": "127.0.0.1", "port": 1080, "protocol": "socks"},
                  {"listen": "127.0.0.1", "port": 8080, "protocol": "http"}
              ]
         }`

	// Default template for test command
    DEF_Test_Template =`
         {
              "log": {"loglevel": "info"}
         }`
)
