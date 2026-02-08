/* This module provides GNU like getopt and getopt_long functions
   Copyright 2025-2026 Ahmad <edu.siamak@gmail.com>

   This software is free software: you can redistribute it and/or modify it
   under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   This software is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
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
	optoff int      = 0    // offset of the current option from -xyz options
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

func is_opt(s string) bool {
	len := len(s)
	if len == 0 || len == 1 {
		return false;
	}
	if s[0] == '-' {
		if len == 1 {
			return false;
		}
		if len == 2 {
			return s[1] != '-';
		}
		return true;
	}
	return false;
}

func parse_lopt(arg string) (opt string, value string, has_value bool) {
	eq := strings.LastIndexByte (arg, '=');
	if eq != -1 {
		opt = arg[2:eq];
		has_value = true;
		if len(arg) > eq+1 {
			value = arg[eq+1:]
		} else {
			value = "" // empty
		}
	} else {
		opt = arg[2:];
		has_value = false;
	}
	return;
}

// @return:  index, does_accept_param
func arg_lookup(arg, optstring string) (int, bool) {
	if idx := strings.IndexByte(optstring, arg[1+optoff]); idx != -1 {
		Optopt = arg[1+optoff];
		if idx+1 < len(optstring) && optstring[idx+1] == ':' {
			return idx, true;
		} else {
			return idx, false;
		}
	} else {
		return -1, false;
	}
}

func Getopt_long(argv []string, optstring string, longopts []Option) int {
beginning_of_parse:
	if Optind >= len(argv) {
		return -1
	}
	arg := argv[Optind]
	if len(arg) < 2 || arg[0] != '-' {
		Optind += 1
		goto beginning_of_parse;
	}
	Optarg = ""

	// cases:  '-x', '-x123', '-xyz'
	if arg[1] != '-' {
		if 1+optoff >= len(arg) {
			optoff = 0
			Optind += 1
			goto beginning_of_parse;
		}
		idx, acc_param := arg_lookup (arg, optstring);
		if idx < 0 {
			Optind += 1
			if Opterr {
				errlog ("%s: invalid option  -- '%s'\n", argv[0], arg[1:2]);
				goto beginning_of_parse;
			} else {
				Optopt = '?'
				return (int)(Optopt);
			}
		}
		if acc_param {
			Optind += 1
			if 2+optoff < len(arg) { // consider the rest of this parameter as option
				Optarg = arg[2+optoff:]
			} else { // use the next parameter
				if Optind < len(argv) && !is_opt(argv[Optind]) {
					Optarg = argv[Optind]
				} else if Opterr {
					errlog("%s: option -%c: requires parameter\n", argv[0], Optopt);
					goto beginning_of_parse;
				} else {
					Optopt = '?'
					return (int)(Optopt);
				}
				Optind += 1
			}
			optoff = 0
		} else {
			optoff += 1
		}
		return (int)(Optopt);
	}
	// cases:  '--key val', '--key=val', '--'
	if arg[1] == '-' {
		if len(arg) == 2 {
			Optopt = '?'
			return -1 // End of Options (--)
		}
		for _,v := range longopts {
			opt, value, has_value := parse_lopt (arg);
			if opt == v.Name { // found
				Optopt = v.Value
				Optind += 1
				if has_value {
					// detected: '--option=value'
					Optarg = value
					return (int)(v.Value);
				} else if v.HasArg && Optind < len(argv) && !is_opt (argv[Optind]) {
					// detected: '--option' 'value'
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
	}

	return 0
}

func Getopt(argv []string, optstring string) int {
beginning_of_parse:
	if Optind >= len(argv) {
		return -1
	}
	arg := argv[Optind]
	if len(arg) < 2 || arg[0] != '-' ||
		              (arg[0] == '-' && arg[1] == '-') {
		Optind += 1
		goto beginning_of_parse;
	}
	Optarg = ""

	// cases:  '-x', '-x123', '-xyz'
	if arg[1] != '-' {
		if 1+optoff >= len(arg) {
			optoff = 0
			Optind += 1
			goto beginning_of_parse;
		}
		idx, acc_param := arg_lookup (arg, optstring);
		if idx < 0 {
			Optind += 1
			if Opterr {
				errlog ("%s: invalid option  -- '%s'\n", argv[0], arg[1:2]);
				goto beginning_of_parse;
			} else {
				Optopt = '?'
				return (int)(Optopt);
			}
		}
		if acc_param {
			Optind += 1
			if 2+optoff < len(arg) { // consider the rest of this parameter as option
				Optarg = arg[2+optoff:]
			} else { // use the next parameter
				if Optind < len(argv) && !is_opt(argv[Optind]) {
					Optarg = argv[Optind]
				} else if Opterr {
					errlog("%s: option -%c: requires parameter\n", argv[0], Optopt);
					goto beginning_of_parse;
				} else {
					Optopt = '?'
					return (int)(Optopt);
				}
				Optind += 1
			}
			optoff = 0
		} else {
			optoff += 1
		}
		return (int)(Optopt);
	}	

	return 0
}
