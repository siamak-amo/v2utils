// logger_debug.go
//go:build debug
// +build debug
package log

import (
	"os"
	"fmt"
	"runtime"
)

func Flogf(out *os.File, typ string, format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		fmt.Fprintf (out,
			"%s:%d: [%s]  " + format,
			append([]interface{}{file, line, typ}, args...)...,
		)
	} else {
		fmt.Fprintf (out,
			"[%s]  " + format,
			append([]interface{}{typ}, args...)...,
		)
	}
}

func MiniFlogf(out *os.File, format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(2); ok {
		fmt.Fprintf (out,
			"%s:%d: " + format,
			append([]interface{}{file, line}, args...)...,
		);
	} else {
		fmt.Fprintf (out, format, args...);
	}
}

func Debugf(format string, args ...interface{}) {
	MiniFlogf (os.Stderr, format, args...);
}
