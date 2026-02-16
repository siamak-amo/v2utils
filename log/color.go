// SPDX-License-Identifier: GPL-3.0-or-later
package log

import (
	"fmt"
)

const (
	COLOR_EMPTY int = iota
	COLOR_BLACK
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_PURPLE
	COLOR_CYAN
	COLOR_WHITE
);

type Color struct {
	BG_color int // background
	FG_color int // foreground
};

var COLOR_RESET string = "\033[0m";

func (c Color) Color_Fmt() string {
	if c.BG_color == COLOR_EMPTY && c.FG_color == COLOR_EMPTY {
		return ""; // nothing to do
	}
	if c.BG_color == COLOR_EMPTY { // setting fg
		return fmt.Sprintf ("\033[3%dm", c.FG_color-1);
	} else if c.FG_color == COLOR_EMPTY { //setting bg
		return fmt.Sprintf ("\033[4%dm", c.BG_color-1);
	} else { // setting both
		return fmt.Sprintf ("\033[3%d;4%dm", c.FG_color-1, c.BG_color-1);
	}
}

func (c Color) Printfmt(val string) string {
	return c.Color_Fmt() + val + COLOR_RESET;
}
