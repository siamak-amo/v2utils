// SPDX-License-Identifier: GPL-3.0-or-later
package log

import "os"

const (
	Verbose int = iota
	Info
	Warning
	Error
	None // no log at all
)

var (
	LogLevel int
	ColorEnabled bool = true

	Log_File = os.Stderr
)

func Must(level int) bool {
	return (Verbose == LogLevel) ||
		(None != LogLevel  &&  level >= LogLevel);
}

func strcolor(color int, str string) string {
	if ! ColorEnabled {
		return str;
	} else {
		c := Color{FG_color: color}
		return c.Printfmt(str);
	}
}

// bypass the level, always print the log
func Logf(format string, args ...interface{}) {
	MiniFlogf (Log_File, format, args...)
}

func Debugf(format string, args ...interface{}) {
	// NOP()  only print it when `-tags debug' is used
}

func Verbosef(format string, args ...interface{}) {
	if (Verbose == LogLevel) {
		MiniFlogf (Log_File, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if Must (Info) {
		Flogf (Log_File, strcolor(COLOR_CYAN, "INFO"), format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if Must (Warning) {
		Flogf (Log_File, strcolor(COLOR_YELLOW, "WARNING"), format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Must (Error) {
		Flogf (Log_File, strcolor(COLOR_RED, "ERROR"), format, args...)
	}
}
