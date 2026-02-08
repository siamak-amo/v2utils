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
	MiniFlogf (os.Stderr, format, args...)
}

func Verbosef(format string, args ...interface{}) {
	if (Verbose == LogLevel) {
		MiniFlogf (os.Stderr, format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if Must (Info) {
		Flogf (os.Stderr, strcolor(COLOR_CYAN, "INFO"), format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if Must (Warning) {
		Flogf (os.Stderr, strcolor(COLOR_YELLOW, "WARNING"), format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Must (Error) {
		Flogf (os.Stderr, strcolor(COLOR_RED, "ERROR"), format, args...)
	}
}
