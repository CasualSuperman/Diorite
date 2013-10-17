package main

import (
	"github.com/conformal/gotk3/gtk"
)

type statusUpdate struct {
	msg  string
	icon *gtk.Image
}

func buildWindow() (*gtk.Window, error) {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	if err != nil {
		return nil, err
	}

	win.SetTitle("Diorite")
	win.Connect("destroy", gtk.MainQuit)

	return win, nil
}
