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

func showTimer(stateChan chan string, closeChan chan int) {
	ticker := time.NewTicker(time.Second * 5)
	initValue := time.Now().Unix()
	currentState := <- stateChan
	for {
		select {
		case t := <-ticker.C:
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: fmt.Sprintf("UnRepeater-%s (%d)s", currentState, t.Unix()-initValue),
			})
		case newState := <- stateChan:
			currentState = newState
		case <-closeChan:
			<-ticker.C
			return
		}
	}
}

func updateSessionStateInLoop(sessionStateListenerChan chan menu.CurrentSessionState) {
	var closeChan chan int
	var stateChan chan string

	for {
		select {
		case state := <-sessionStateListenerChan:
			switch state {
			case menu.SessionInProgress:
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater-Running",
				})
				closeChan = make(chan int)
				stateChan = make(chan string, 1)
				stateChan <- "Running"
				go showTimer(stateChan, closeChan)
			case menu.SessionStopInProgress:
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater-Stopping",
				})
				stateChan <- "Stopping"
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
	sessionMenuActionListenerChan := make(chan menu.CurrentSessionState)
	go updateSessionStateInLoop(sessionMenuActionListenerChan)

	taskSessionWorker := worker.NewTaskSession()
	sessionManagerMenu = menu.NewSessionManager(taskSessionWorker, sessionMenuActionListenerChan)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "UnRepeater",
	})
	menuet.App().Children = getMenuItems
	menuet.App().Label = "github.com/krishnakumar4a4/un-repeater"
}
