// SPDX-License-Identifier: GPL-3.0-or-later
package log

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
	COLOR_a
	COLOR_b
);


type Color struct {
	BG_color int // background
	FG_color int // foreground
};

// equivalent to: "\033[xxm"
var COLOR_RESET string = "\033[0m";

func (c Color) Color_Fmt() string {
	// this is equivalent to: fmt.Sprintf ("\033[3%d;4%dm", c.FG_color, c.BG_color)
	// but it should be more efficient
	var templateANSI = [8]byte{033, '[', 'x','x', 'm', '4','x', 'm'};

	if c.BG_color == COLOR_EMPTY { // setting fg
		templateANSI[2] = '3';
		templateANSI[3] = (byte)(c.FG_color - 1 + '0');
		return string(templateANSI[0:5]);
	} else if c.FG_color == COLOR_EMPTY { // setting bg
		templateANSI[2] = '4';
		templateANSI[3] = (byte)(c.BG_color - 1 + '0');
		return string(templateANSI[0:5]);
	} else {
		templateANSI[2] = '3';
		templateANSI[3] = (byte)(c.FG_color - 1 + '0');
		templateANSI[4] = ';';
		templateANSI[6] = (byte)(c.BG_color - 1 + '0');
		return string(templateANSI[:])
	}
}

func (c Color) Printfmt(val string) string {
	return c.Color_Fmt() + val + COLOR_RESET;
}
