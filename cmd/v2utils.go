// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"os"
	"fmt"

	flag "github.com/spf13/pflag"
	log "github.com/siamak-amo/v2utils/log"
)

func (opt *Opt) RegisterFlag() {
	opt.url = flag.String ("url", "", "proxy URL e.g. vless://xxx");
	opt.in_file = flag.String ("input", "", "path to proxy URLs file");
	opt.template_file = flag.String ("template", "", "path to json template file");
	opt.output_dir = flag.String ("output", "", "output directory path");
	opt.configs = flag.String ("config", "", "path to config file or dir");
	opt.Verbose = flag.Bool ("verbose", false, "verbose");
	opt.rm = flag.Bool ("rm", false, "remove broken and invalid files")

	flag.Parse();
}

// returns negative on fatal failures
func (opt *Opt) ParseFlags() int {
	opt.RegisterFlag();
	argv := flag.Args()
	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "error:  missing COMMAND")
		fmt.Fprintln(os.Stderr, "usage:  v2utils [COMMAND] [OPTIONS]")
		return -1
	}
	switch (argv[0]) {
	case "convert","CONVERT", "conv", "c","C":
		if "" != *opt.configs {
			opt.Cmd = CMD_CONVERT_CFG;
		} else {
			// Default: converting URLs (stdin / --url)
			opt.Cmd = CMD_CONVERT;
		}
		break;

	case "test","Test","TEST", "t","T":
		if "" != *opt.configs {
			// For Testing config (.json) files
			opt.Cmd = CMD_TEST_CFG;
		} else {
			// Default: testing URLs
			opt.Cmd = CMD_TEST;
		}
		break;

	case "run","Run","RUN", "r","R":
		if "" != *opt.configs {
			opt.Cmd = CMD_RUN_CFG;
		} else if "" != *opt.url {
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
	switch (opt.Cmd) {
	case CMD_RUN:
		opt.Set_rd_url()
		break;

	case CMD_TEST, CMD_CONVERT:
		if "" != *opt.output_dir {
			if err := os.MkdirAll(*opt.output_dir, 0o755); nil != err {
				log.Errorf ("Could not create dir - %v\n", err);
				return -1
			}
		}
		if "" != *opt.url {
			opt.Set_rd_url();
		} else if "" != *opt.in_file {
			if e := opt.Set_rd_file(); nil != e {
				log.Errorf ("%v\n", e);
				return -1
			}
		} else {
			opt.Set_rd_stdin();
		}
		break;

	case CMD_CONVERT_CFG, CMD_TEST_CFG, CMD_RUN_CFG:
		if "-" == *opt.configs {
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
				return
			}
			break;

		case CMD_RUN:
			if "" == *opt.template_file {
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
				if *opt.Verbose {
					log.Infof("`%s` OK.\n", ln)
				} else {
					println(ln)
				}
				if "" != *opt.output_dir {
					opt.CFG_Out(ln); // Also generate json files
				}
			} else {
				log.Infof("Broken URL '%s'\n", ln);
			}
			break;

			// For xxx_CFG commands, @ln is path to a file or `-` for stdin
		case CMD_TEST_CFG:
			if opt.Test_CFG(ln) {
				if *opt.Verbose {
					fmt.Printf("config file `%s':  OK.\n", ln)
				} else {
					println(ln)
				}
			} else {
				if *opt.rm {
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
			opt.template_file = &ln
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
			log.Errorf("(%s) CMD_CONVERT_CFG -- Not Implemented.\n", ln);
			break;
		}
	}
}

func main() {
	opt := Opt{};
	if ret := opt.ParseFlags(); ret < 0 {
		os.Exit (-ret);
	}
	if ret := opt.Init(); ret < 0 {
		os.Exit (-ret);
	}
	if *opt.Verbose {
		log.LogLevel = log.Verbose;
	}
	opt.Do(); // The main loop
}
