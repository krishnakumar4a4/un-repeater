package main

import (
	"fmt"
	"time"

	"github.com/krishnakumar4a4/un-repeater/menu"
	"github.com/krishnakumar4a4/un-repeater/worker"

	"github.com/caseymrm/menuet"
)

var sessionManagerMenu *menu.SessionManager

func main() {
	go showMenu()
	menuet.App().RunApplication()
}

func getMenuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}
	items = append(items, sessionManagerMenu.StartSessionMenuItem())
	return items
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

func updateSessionStateInLoop(sessionStateListenerChan chan menu.CurrentSessionState) {
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

func showMenu() {
	sessionListenerChan := make(chan menu.CurrentSessionState)
	go updateSessionStateInLoop(sessionListenerChan)

	taskSessionWorker := worker.NewTaskSession()
	sessionManagerMenu = menu.NewSessionManager(taskSessionWorker, sessionListenerChan)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "UnRepeater",
	})
	menuet.App().Children = getMenuItems
	menuet.App().Label = "github.com/krishnakumar4a4/un-repeater"
}
