// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
	colorful "github.com/lucasb-eyer/go-colorful"
)

type demoColorsTrue struct {
	dark bool
}

func mainColorsTrue() {
	_ = os.Setenv("COLORTERM", "truecolor")
	g, err := gocui.NewGui(gocui.OutputTrue, true)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	d := &demoColorsTrue{dark: false}
	g.SetManagerFunc(d.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, d.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if d.dark {
			d.dark = false
		} else {
			d.dark = true
		}
		d.displayHsv(v)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func (d *demoColorsTrue) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	rows := 33
	cols := 182
	if maxY < rows {
		rows = maxY
	}
	if maxX < cols {
		cols = maxX
	}

	if v, err := g.SetView("colors", 0, 0, cols-1, rows-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.FrameColor = gocui.GetColor("#FFAA55")
		d.displayHsv(v)

		if _, err := g.SetCurrentView("colors"); err != nil {
			return err
		}
	}
	return nil
}

func (d *demoColorsTrue) displayHsv(v *gocui.View) {
	v.Clear()
	str := ""
	// HSV color space (lines are value or saturation)
	for i := 50; i > 0; i -= 2 {
		// Hue
		for j := 0; j < 360; j += 2 {
			ir, ig, ib := d.hsv(j, i-1)
			ir2, ig2, ib2 := d.hsv(j, i)
			str += fmt.Sprintf("\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm▀\x1b[0m", ir, ig, ib, ir2, ig2, ib2)
		}
		str += "\n"
		_, _ = fmt.Fprint(v, str)
		str = ""
	}

	_, _ = fmt.Fprintln(v, "\n\x1b[38;5;245mCtrl + R - Switch light/dark mode")
	_, _ = fmt.Fprintln(v, "\nCtrl + C - Exit\n")
	_, _ = fmt.Fprint(v, "Example should enable true color, but if it doesn't work run this command: \x1b[0mexport COLORTERM=truecolor")
}

func (d *demoColorsTrue) hsv(hue, sv int) (uint32, uint32, uint32) {
	if !d.dark {
		ir, ig, ib, _ := colorful.Hsv(float64(hue), float64(sv)/50, float64(1)).RGBA()
		return ir >> 8, ig >> 8, ib >> 8
	}
	ir, ig, ib, _ := colorful.Hsv(float64(hue), float64(1), float64(sv)/50).RGBA()
	return ir >> 8, ig >> 8, ib >> 8
}

func (d *demoColorsTrue) quit(g *gocui.Gui, v *gocui.View) error {
	_ = g
	_ = v
	return gocui.ErrQuit
}
