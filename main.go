package main

import (
	"fmt"
	"time"
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
	sessionListenerChan := make(chan menu.CurrentSessionState)
	go updateSessionState(sessionListenerChan)
	sessionMenu = menu.NewStartStopSession(workerSession, sessionListenerChan)
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

func updateSessionState(sessionStateListenerChan chan menu.CurrentSessionState) {
	var closeChan chan int
	for {
		select {
		case state := <-sessionStateListenerChan:
			switch state {
			case menu.SessionInProgress:
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater-Running",
				})
				closeChan = make(chan int)
				go showTimer(closeChan)
			case menu.SessionStopped:
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater",
				})
				close(closeChan)
			}
		}
	}
}

func showTimer(closeChan chan int) {
	ticker := time.NewTicker(time.Second * 5)
	initValue := time.Now().Unix()
	for {
		select {
		case t := <-ticker.C:
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: fmt.Sprintf("UnRepeater-Running (%d)s", t.Unix()-initValue),
			})
		case <-closeChan:
			<-ticker.C
			return
		}
	}
}
