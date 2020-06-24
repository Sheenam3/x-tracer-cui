package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"github.com/jroimartin/gocui"
)

var version = "1.14.4"
var LOG_MOD string = "pod"
var NAMESPACE string = "default"

// Configure globale keys
var keys []Key = []Key{
	Key{"", gocui.KeyCtrlC, actionGlobalQuit},
	Key{"", gocui.KeyCtrlD, actionGlobalToggleViewDebug},
	Key{"pods", gocui.KeyCtrlN, actionGlobalToggleViewNamespaces},
	Key{"pods", gocui.KeyArrowUp, actionViewPodsUp},
	Key{"pods", gocui.KeyArrowDown, actionViewPodsDown},
	Key{"pods", gocui.KeyEnter, actionViewPodsSelect},
	//Key{"pods", 'd', actionViewPodsDelete},
	//Key{"pods", 'l', actionViewPodsLogs},
	//Key{"logs", 'l', actionViewPodsLogsHide},
	//Key{"logs", gocui.KeyArrowUp, actionViewPodsLogsUp},
	//Key{"logs", gocui.KeyArrowDown, actionViewPodsLogsDown},
	Key{"namespaces", gocui.KeyArrowUp, actionViewNamespacesUp},
	Key{"namespaces", gocui.KeyArrowDown, actionViewNamespacesDown},
	Key{"namespaces", gocui.KeyEnter, actionViewNamespacesSelect},
}

// Main or not main, that's the question^^
func main() {
	c := getConfig()

	// Ask version
	if c.askVersion {
		fmt.Println(versionFull())
		os.Exit(0)
	}

	// Ask Help
	if c.askHelp {
		fmt.Println(versionFull())
		fmt.Println(HELP)
		os.Exit(0)
	}

	// Only used to check errors
	getClientSet()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(uiLayout)

	if err := uiKey(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}


func uiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

//	viewDebug(g, maxX, maxY)
	viewLogs(g, maxX, maxY)
	viewInfo(g,maxX,maxY)
	viewNamespaces(g, maxX, maxY)
	viewOverlay(g, maxX, maxY)
	viewTitle(g, maxX, maxY)
	viewPods(g, maxX, maxY)
	viewStatusBar(g, maxX, maxY)

	return nil
}



// Move view cursor to the bottom
func moveViewCursorDown(g *gocui.Gui, v *gocui.View, allowEmpty bool) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, err := getNextViewLine(g, v)
	if err != nil {
		return err
	}
	if !allowEmpty && nextLine == "" {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// Move view cursor to the top
func moveViewCursorUp(g *gocui.Gui, v *gocui.View, dY int) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if cy > dY {
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get view line (relative to the cursor)
func getViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, err
}

// Get the next view line (relative to the cursor)
func getNextViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}

// Set view cursor to line
func setViewCursorToLine(g *gocui.Gui, v *gocui.View, lines []string, selLine string) error {
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	for y, line := range lines {
		if line == selLine {
			if err := v.SetCursor(ox, y); err != nil {
				if err := v.SetOrigin(cx, y); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get pod name form line
func getPodNameFromLine(line string) string {
	if line == "" {
		return ""
	}

	i := strings.Index(line, " ")
	if i == -1 {
		return line
	}

	return line[0:i]
}

// Get selected pod
func getSelectedPod(g *gocui.Gui) (string, error) {
	v, err := g.View("pods")
	if err != nil {
		return "", err
	}
	l, err := getViewLine(g, v)
	if err != nil {
		return "", err
	}
	p := getPodNameFromLine(l)

	return p, nil
}



//Get Host Name
func getHostName() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Println("Hostname : ", name)
	return name
}

func showViewPodsLogs(g *gocui.Gui) error {
	vn := "logs"

	switch LOG_MOD {
	case "pod":
		// Get current selected pod
		p, err := getSelectedPod(g)
		if err != nil {
			return err
		}

		// Display pod containers
		vLc, err := g.View(vn + "-containers")
		if err != nil {
			return err
		}
		vLc.Clear()

		var conName []string
		for _, c := range getPodContainers(p) {
			fmt.Fprintln(vLc, c)
			conName = append(conName, c)
		}
		vLc.SetCursor(0, 0)


  	        //Display Container IDs
		lv, err := g.View(vn)
		if err != nil {
                        return err
                }
		lv.Clear()

		fmt.Fprintln(lv, "Containers ID are:")
		for i, conId := range getPodContainersID(p){
			fmt.Fprintln(lv, conName[i] + "->" + conId)

		}


	       

		

		// Display logs
		//refreshPodsLogs(g)
	}

	

//	debug(g, "Action: Show view logs")
	g.SetViewOnTop(vn)
	g.SetViewOnTop(vn + "-containers")
	g.SetCurrentView(vn)

	return nil
}


func displayError(g *gocui.Gui, e error) error {
	lMaxX, lMaxY := g.Size()
	minX := lMaxX / 6
	minY := lMaxY / 6
	maxX := 5 * (lMaxX / 6)
	maxY := 5 * (lMaxY / 6)

	if v, err := g.SetView("errors", minX, minY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Title = " ERROR "
		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true
		v.BgColor = gocui.ColorRed
		v.FgColor = gocui.ColorWhite

		// Content
		v.Clear()
		fmt.Fprintln(v, e.Error())

		// Send to forground
		g.SetCurrentView(v.Name())
	}

	return nil
}

// Hide error box
func hideError(g *gocui.Gui) {
	g.DeleteView("errors")
}

// Display confirmation message
func displayConfirmation(g *gocui.Gui, m string) error {
	lMaxX, lMaxY := g.Size()

	if v, err := g.SetView("confirmation", -1, lMaxY-3, lMaxX, lMaxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false

		// Content
		fmt.Fprintln(v, textPadCenter(m, lMaxX))

		// Auto-hide message
		hide := func() {
			hideConfirmation(g)
		}
		time.AfterFunc(time.Duration(2)*time.Second, hide)
	}

	return nil
}

// Hide confirmation message
func hideConfirmation(g *gocui.Gui) {
	g.DeleteView("confirmation")
}

