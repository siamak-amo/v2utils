// SPDX-License-Identifier: GPL-3.0-or-later
package getopt

import (
	"os"
	"fmt"
	"strings"
)

const (
	No_argument       bool = false
	Required_argument bool = true
);

var (
	Optind int      = 1    // the index of the next element to be processed in argv
	Optarg string   = ""   // the text of the following argv-element
	Optopt byte     = '?'  // the current option argument
	Opterr bool     = true // if it's set to false, getopt does not print error message
                           // the caller can determine errors by getopt return value == '?'
);

type Option struct {
	Name     string  // long name --Name
	HasArg   bool
	Value    byte  // the value to return
};

func Getopt_reset() {
	Optind = 1
	Optopt = '?'
	Optarg = ""
}

func errlog(format string, args ...any) {
	if Opterr {
		fmt.Fprintf(os.Stderr, format, args...);
	}
}

func Getopt_long(argv []string, optstring string, longopts []Option) int {
beginning_of_parse:
	if Optind >= len(argv) {
		return -1
	}
	arg := argv[Optind]
	if len(arg) <= 1 {
		return -1;
	}
	Optarg = ""

	if arg[0] == '-' && arg[1] != '-' {
		// Checking for options: -x
		Optopt = arg[1]
		if idx := strings.IndexByte(optstring, arg[1]); idx != -1 {
			if idx+1 < len(optstring) && optstring[idx+1] == ':' { // accepts arg
				if len(arg) == 2 {  // -x
					Optind += 1
					if Optind < len(argv) && argv[Optind][0] != '-' {
						Optarg = argv[Optind]
						Optind += 1
					} else {
						if Opterr {
							errlog("%s: option %s: requires parameter\n", argv[0], arg);
							goto beginning_of_parse;
						} else {
							return '?';
						}
					}
				} else if len(arg) > 2 { // -xXXX GNU style
					Optarg = argv[Optind][2:]
					Optind += 1
				}
			} else { // no arg
				Optind += 1
			}
			return (int)(arg[1])
		} else {
			Optind += 1
			if Opterr {
				errlog("%s: invalid option  -- '%s'\n", argv[0], arg[1:2]);
				goto beginning_of_parse;
			} else {
				return '?'
			}
		}
	} else if arg[0] == '-' && arg[1] == '-' {
		// Checking for long options --xxx
		if len(arg) == 2 {
			Optopt = '?'
			return -1 // End of Options (--)
		}
		for _,v := range longopts {
			if arg[2:] == v.Name { // found
				Optopt = v.Value
				Optind += 1
				if v.HasArg && Optind < len(argv) && argv[Optind][0] != '-' {
					Optarg = argv[Optind]
					Optind += 1
				} else if v.HasArg {
					if Opterr {
						errlog("%s: option %s: requires parameter\n", argv[0], arg);
						goto beginning_of_parse;
					} else {
						return '?'
					}
				}
				return (int)(v.Value)
			}
		}
		Optopt = '?'
		Optind += 1
		if Opterr {
			errlog("%s: unrecognized option '%s'\n", argv[0], arg);
			goto beginning_of_parse;
		} else {
			return '?'
		}
	} else {
		Optind += 1
		goto beginning_of_parse;
	}

	return 0
}

func Getopt(argv []string, optstring string) int {
beginning_of_parse:
	if Optind >= len(argv) {
		return -1
	}
	arg := argv[Optind]
	if len(arg) <= 1 {
		return -1;
	}
	Optarg = ""

	if arg[0] == '-' && arg[1] != '-' {
		// Checking for options: -x
		Optopt = arg[1]
		if idx := strings.IndexByte(optstring, arg[1]); idx != -1 {
			if idx+1 < len(optstring) && optstring[idx+1] == ':' { // accepts arg
				if len(arg) == 2 {  // -x
					Optind += 1
					if Optind < len(argv) && argv[Optind][0] != '-' {
						Optarg = argv[Optind]
						Optind += 1
					} else {
						if Opterr {
							errlog("%s: option %s: requires parameter\n", argv[0], arg);
							goto beginning_of_parse;
						} else {
							return '?'
						}
					}
				} else if len(arg) > 2 { // -xXXX GNU style
					Optarg = argv[Optind][2:]
					Optind += 1
				}
			} else { // no arg
				Optind += 1
			}
			return (int)(arg[1])
		} else {
			Optind += 1
			if Opterr {
				errlog("%s: invalid option  -- '%s'\n", argv[0], arg[1:2]);
				goto beginning_of_parse;
			} else {
				return '?'
			}
		}
	} else if arg[0] == '-' && arg[1] == '-' {
		if Opterr {
			errlog("%s: unrecognized option '%s'\n", argv[0], arg);
		} else {
			return '?'
		}
	}
	return 0
}
