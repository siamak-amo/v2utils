// SPDX-License-Identifier: GPL-3.0-or-later
package log

import "os"

func Infof(format string, args ...interface{}) {
	Flogf (os.Stderr, "Info", format, args...)
}

func Warnf(format string, args ...interface{}) {
	Flogf (os.Stderr, "Warning", format, args...)
}

func Errorf(format string, args ...interface{}) {
	Flogf (os.Stderr, "Error", format, args...)
}
