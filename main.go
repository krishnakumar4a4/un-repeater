package main

import (
	"unrepeater/menu"
	"unrepeater/worker"

	"github.com/caseymrm/menuet"
)

var sessionMenu *menu.StartStopSession

func main() {
	go showMenu()
	menuet.App().RunApplication()
}

func showMenu() {
	workerSession := worker.NewSession()
	sessionMenu = menu.NewStartStopSession(workerSession)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "UnRepeater",
	})
	menuet.App().Children = getMenuItems
	menuet.App().Label = "UnRepeater"
}

func getMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, sessionMenu.StartSessionMenuItem())
	return items
}
