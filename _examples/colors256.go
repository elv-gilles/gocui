// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

type demoColors256 struct{}

func mainColors256() {
	g, err := gocui.NewGui(gocui.Output256, true)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	d := demoColors256{}
	g.SetManagerFunc(d.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, d.quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func (d *demoColors256) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("colors", -1, -1, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		// 256-colors escape codes
		for i := 0; i < 256; i++ {
			str := fmt.Sprintf("\x1b[48;5;%dm\x1b[30m%3d\x1b[0m ", i, i)
			str += fmt.Sprintf("\x1b[38;5;%dm%3d\x1b[0m ", i, i)

			if (i+1)%10 == 0 {
				str += "\n"
			}

			_, _ = fmt.Fprint(v, str)
		}

		_, _ = fmt.Fprint(v, "\n\n")

		// 8-colors escape codes
		ctr := 0
		for i := 0; i <= 7; i++ {
			for _, j := range []int{1, 4, 7} {
				str := fmt.Sprintf("\x1b[3%d;%dm%d:%d\x1b[0m ", i, j, i, j)
				if (ctr+1)%20 == 0 {
					str += "\n"
				}

				_, _ = fmt.Fprint(v, str)

				ctr++
			}
		}
		if _, err := g.SetCurrentView("colors"); err != nil {
			return err
		}
	}
	return nil
}

func (d *demoColors256) quit(g *gocui.Gui, v *gocui.View) error {
	_ = g
	_ = v
	return gocui.ErrQuit
}
