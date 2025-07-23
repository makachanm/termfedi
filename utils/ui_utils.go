package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func WriteTo(t tcell.Screen, x int, y int, content string, style tcell.Style) {
	curr_x, curr_y := x, y

	for _, cstr := range content {
		var ctr []rune

		width := runewidth.RuneWidth(cstr)
		if width == 0 {
			width = 1
			ctr = []rune{cstr}
			cstr = ' '
		}

		t.SetContent(curr_x, curr_y, cstr, ctr, style)
		curr_x += width
	}
}
