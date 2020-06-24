package main

import (
	"github.com/jroimartin/gocui"
)

var DEBUG_DISPLAYED bool = false
var NAMESPACES_DISPLAYED bool = false

// Global action: Quit
func actionGlobalQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Global action: Toggle debug
func actionGlobalToggleViewDebug(g *gocui.Gui, v *gocui.View) error {
	vn := "debug"

	if !DEBUG_DISPLAYED {
//		debug(g, "Action: Display debug popup")
		g.SetViewOnTop(vn)
		g.SetCurrentView(vn)
	} else {
//		debug(g, "Action: Hide debug popup")
		g.SetViewOnBottom(vn)
		g.SetCurrentView("pods")
	}

	DEBUG_DISPLAYED = !DEBUG_DISPLAYED

	return nil
}

// View namespaces: Toggle display
func actionGlobalToggleViewNamespaces(g *gocui.Gui, v *gocui.View) error {
	vn := "namespaces"

	if !NAMESPACES_DISPLAYED {
//		debug(g, "Action: Display namespaces popup")
		g.SetViewOnTop(vn)
		g.SetCurrentView(vn)
		changeStatusContext(g, "SE")
	} else {
//		debug(g, "Action: Hide namespaces popup")
		g.SetViewOnBottom(vn)
		g.SetCurrentView("pods")
		changeStatusContext(g, "D")
	}

	NAMESPACES_DISPLAYED = !NAMESPACES_DISPLAYED

	return nil
}

// View pods: Up
func actionViewPodsUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 2)
//	debug(g, "Select up in pods view")
	return nil
}

// View pods: Down
func actionViewPodsDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
//	debug(g, "Select down in pods view")
	return nil
}

func actionViewPodsSelect(g *gocui.Gui, v *gocui.View) error {
	line,err  := getViewLine(g,v)
	if err != nil {
		return err
	}
//	maxX, maxY := g.Size()	
	LOG_MOD = "pod"
	errr := showViewPodsLogs(g)

	changeStatusContext(g, "SL")
//	viewLogs(g, maxX, maxY)
	displayConfirmation(g, line+" Pod selected")
	return errr

}

// View namespaces: Up
func actionViewNamespacesUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 0)
//	debug(g, "Select up in namespaces view")
	return nil
}

// View namespaces: Down
func actionViewNamespacesDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
//	debug(g, "Select down in namespaces view")
	return nil
}

// Namespace: Choose
func actionViewNamespacesSelect(g *gocui.Gui, v *gocui.View) error {
	line, err := getViewLine(g, v)
//	debug(g, "Select namespace: "+line)
	NAMESPACE = line
	go viewPodsRefreshList(g)
	actionGlobalToggleViewNamespaces(g, v)
	displayConfirmation(g, line+" namespace selected")
	return err
}
