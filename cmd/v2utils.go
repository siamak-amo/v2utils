/* V2utils provides xray-core compatible utilities
   Copyright 2025 Ahmad <edu.siamak@gmail.com>

   V2utils is free software: you can redistribute it and/or modify it
   under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   V2utils is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"os"
	"fmt"

	log "github.com/siamak-amo/v2utils/log"
	getopt "github.com/siamak-amo/v2utils/getopt"
)

const Version = "1.1";

func (opt *Opt) GetArgs() {
	const optstr = "u:f:t:o:c:rvh"
	lopts := []getopt.Option{
		{"url",           true,  'u'},
		{"config",        true,  'c'},
		{"template",      true,  't'},
		{"output",        true,  'o'},
		{"input",         true,  'i'},
		{"rm",            false, 'r'},
		{"verbose",       false, 'v'},
		{"help",          false, 'h'},
	}
	argv := os.Args
	for idx := 0; -1 != idx; {
		idx = getopt.Getopt_long (argv, optstr, lopts);

		switch (idx) {
		case 'u':
			opt.url = getopt.Optarg; break;
		case 'c':
			opt.configs = getopt.Optarg; break;
		case 't':
			opt.template_file = getopt.Optarg; break;
		case 'i':
			opt.in_file = getopt.Optarg; break;
		case 'o':
			opt.output_dir = getopt.Optarg; break;
		case 'r':
			opt.rm = true; break;
		case 'v':
			opt.verbose = true; break;
		case 'h':
			fmt.Fprintf(os.Stderr, `v2utils v%s - xray-core compatible utility
Usage:  v2utils COMMAND [OPTIONS]

COMMAND:
      Run:  to execute a Xray intense based on the given configuration
     Test:  to test the current configuration has internet access
  Convert:  to convert the current configuration to a different format

OPTIONS:
    -u, --url             VPN url (e.g. vless:// trojan://)
    -c, --config          path to config file or folder
    -t, --template        path to template file
                          (for Run and Convert commands)
    -i, --input           path to input URL file
    -o, --output          path to output folder
    -r, --rm              to remove broken config files
                          (only in Test command)
    -v, --verbose         verbose

Examples:
    # run xray by URL:
    $ v2utils run --url vless://id@1.2.3.4:1234

    # test json files and remove broken ones
    $ v2utils test --config /path/to/configs/ --rm

    # convert URLs to json:
    $ cat url.txt | v2utils convert -o /path/to/configs

    # convert outbound of json files to URL:
    $ v2utils convert --config /path/to/configs_dir
`, Version);
			os.Exit(0);
		}
	}
}

// returns negative on fatal failures
func (opt *Opt) HandleArgs() int {
	argv := os.Args
	if len(argv) < 2 {
		fmt.Fprintln(os.Stderr, "error:  missing COMMAND")
		fmt.Fprintln(os.Stderr, "usage:  v2utils [COMMAND] [OPTIONS]")
		return -1
	}
	switch (argv[1]) {
	case "convert","CONVERT", "conv", "c","C":
		if "" != opt.configs {
			opt.Cmd = CMD_CONVERT_CFG;
		} else {
			// Default: converting URLs (stdin / --url)
			opt.Cmd = CMD_CONVERT;
		}
		break;

	case "test","Test","TEST", "t","T":
		if "" != opt.configs {
			// For Testing config (.json) files
			opt.Cmd = CMD_TEST_CFG;
		} else {
			// Default: testing URLs
			opt.Cmd = CMD_TEST;
		}
		break;

	case "run","Run","RUN", "r","R":
		if "" != opt.configs {
			opt.Cmd = CMD_RUN_CFG;
		} else if "" != opt.url {
			opt.Cmd = CMD_RUN;
		} else {
			log.Errorf("Run command needs a URL (--url) or config file (--config).");
			return -1
		}
		break;

	default:
		println ("Invalid command.");
		return -1
	}
	return 0;
}

// returns negative on fatal failures
func (opt *Opt) Init() int {
	if opt.verbose {
		log.LogLevel = log.Verbose;
	} else {
		log.LogLevel = log.Warning; // Default level
	}

	switch (opt.Cmd) {
	case CMD_RUN:
		opt.Set_rd_url()
		break;

	case CMD_TEST, CMD_CONVERT:
		if "" != opt.output_dir {
			if err := os.MkdirAll(opt.output_dir, 0o755); nil != err {
				log.Errorf ("Could not create dir - %v\n", err);
				return -1
			}
		}
		if "" != opt.url {
			opt.Set_rd_url();
		} else if "" != opt.in_file {
			if e := opt.Set_rd_file(); nil != e {
				log.Errorf ("%v\n", e);
				return -1
			}
		} else {
			opt.Set_rd_stdin();
		}
		break;

	case CMD_CONVERT_CFG, CMD_TEST_CFG, CMD_RUN_CFG:
		if "-" == opt.configs {
			opt.Set_rd_cfg_stdin();
		} else {
			opt.Set_rd_cfg();
		}
		break;
	}
	return 0;
}

func (opt Opt) Do() {
	EOF := false
	var ln string

	for true != EOF {
		ln, EOF = opt.GetInput();
		if len(ln) == 0 || ln[0] <= ' ' || ln[0] == '#' {
			continue;
		}

		switch (opt.Cmd) {
		case CMD_CONVERT:
			if e := opt.Convert_url2json(ln); nil != e {
				if opt.verbose {
					log.Warnf("Convert error: %v\n", e)
				}
			}
			break;

		case CMD_RUN:
			if "" == opt.template_file {
				log.Errorf("No template provided, using the default template\n");
				fmt.Printf("Default template:%s\n", opt.Get_Default_Template()); // should print this
			}
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Invalid template - %v\n", e)
				return;
			}
			if e := opt.Init_Outbound_byURL(ln); nil != e {
				log.Errorf("Invalid or unsupported URL - %v\n", e)
				break;
			}
			if e := opt.Exec_Xray(); nil != e {
				log.Errorf("Exec xray-core failed - %v\n", e)
				break;
			}
			break;

		case CMD_TEST:
			if opt.Test_URL(ln) {
				if opt.verbose {
					log.Infof("`%s` OK.\n", ln)
				} else {
					fmt.Println(ln)
				}
				if "" != opt.output_dir {
					opt.CFG_Out(ln); // Also generate json files
				}
			} else {
				log.Infof("Broken URL '%s'\n", ln);
			}
			break;

			// For xxx_CFG commands, @ln is path to a file or `-` for stdin
		case CMD_TEST_CFG:
			if opt.Test_CFG(ln) {
				if opt.verbose {
					fmt.Printf("config file `%s':  OK.\n", ln)
				} else {
					fmt.Println(ln);
				}
			} else {
				if opt.rm {
					fmt.Printf("Broken file %s was removed.\n", ln)
					if e := os.Remove(ln); nil != e {
						log.Errorf("Could not remove %s - %v\n", ln, e)
					}
				} else {
					log.Infof("Broken config file %s\n", ln)
				}
			}
			break;

		case CMD_RUN_CFG:
			// Run only uses the first json file @ln, so we used return here
			opt.template_file = ln
			if "-" != ln {
				log.Infof("Running xray-core with config: %s\n", ln)
			}
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Loading config failed - %v\n", e)
				return;
			}
			if e := opt.Exec_Xray(); nil != e {
				log.Errorf("Exec xray-core failed - %v\n", e)
				return;
			}
			return;

		case CMD_CONVERT_CFG:
			opt.template_file = ln
			if e := opt.Init_CFG(); nil != e {
				log.Errorf("Loading config failed - %v\n", e)
				return;
			}
			res, e := opt.Convert_conf2url();
			if nil != e {
				log.Warnf ("Converting to URL failed - %v\n", e);
			} else {
				fmt.Println(res);
			}
			break;
		}
	}
}

func main() {
	opt := Opt{};
	opt.GetArgs();
	if ret := opt.HandleArgs(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}

	opt.Do(); // The main loop
}
