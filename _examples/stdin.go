// Copyright 2015 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example doesn't work when running `go run stdin.go`, you are suposed to pipe someting to this like: `/bin/ls | go run stdin.go`

package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
)

type demoStdin struct{}

func mainStdin() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.Cursor = true

	d := &demoStdin{}
	g.SetManagerFunc(d.layout)

	if err := d.initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatalln(err)
	}
}

func (d *demoStdin) layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("help", maxX-23, 0, maxX-1, 5, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		_, _ = fmt.Fprintln(v, "KEYBINDINGS")
		_, _ = fmt.Fprintln(v, "↑ ↓: Seek input")
		_, _ = fmt.Fprintln(v, "a: Enable autoscroll")
		_, _ = fmt.Fprintln(v, "^C: Exit")
	}

	if v, err := g.SetView("stdin", 0, 0, 80, 35, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = true

		if _, err := io.Copy(hex.Dumper(v), os.Stdin); err != nil {
			return err
		}

		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
		}
	}

	return nil
}

func (d *demoStdin) initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, d.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", 'a', gocui.ModNone, d.autoscroll); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_ = d.scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_ = d.scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (d *demoStdin) quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (d *demoStdin) autoscroll(_ *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = true
	return nil
}

func (d *demoStdin) scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
