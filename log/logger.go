// logger.go
//go:build !debug
// +build !debug
package log

import (
	"os"
	"fmt"
)

func Flogf(out *os.File, typ string, format string, args ...interface{}) {
	fmt.Fprintf (out,
		"[%s]  " + format,
		append([]interface{}{typ}, args...)...,
	)
}

func MiniFlogf(out *os.File, format string, args ...interface{}) {
	fmt.Fprintf (out, format, args...)
}
