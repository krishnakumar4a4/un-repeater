package main

import (
	"fmt"
	"time"

	"github.com/krishnakumar4a4/un-repeater/menu"
	"github.com/krishnakumar4a4/un-repeater/worker"

	"github.com/caseymrm/menuet"
)

var sessionManagerMenu *menu.SessionManager
var sessionState menu.CurrentSessionState

func main() {
	go showMenu()
	menuet.App().RunApplication()
}

func getMenuItems(lister *worker.ScriptLister) func() []menuet.MenuItem {
	return func() []menuet.MenuItem {
		items := []menuet.MenuItem{}
		if sessionState == menu.SessionInProgress {
			items = append(items, sessionManagerMenu.StopSessionMenuItem())
		} else {
			items = append(items, sessionManagerMenu.StartSessionMenuItem())
		}
		items = append(items, menu.GetMenuItems(lister)...)
		return items
	}
}

func showTimer(stateChan chan string, closeChan chan int) {
	ticker := time.NewTicker(time.Second * 5)
	initValue := time.Now().Unix()
	currentState := <-stateChan
	for {
		select {
		case t := <-ticker.C:
			menuet.App().SetMenuState(&menuet.MenuState{
				Title: fmt.Sprintf("UnRepeater-%s (%d)s", currentState, t.Unix()-initValue),
			})
		case newState := <-stateChan:
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
				sessionState = menu.SessionInProgress
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater-Running",
				})
				closeChan = make(chan int)
				stateChan = make(chan string, 1)
				stateChan <- "Running"
				go showTimer(stateChan, closeChan)
			case menu.SessionStopInProgress:
				sessionState = menu.SessionStopInProgress
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater-Stopping",
				})
				stateChan <- "Stopping"
			case menu.SessionStopped:
				sessionState = menu.SessionStopped
				menuet.App().SetMenuState(&menuet.MenuState{
					Title: "UnRepeater",
				})
				menuet.App().MenuChanged()
				close(closeChan)
			}
		}
	}
}

func showMenu() {
	sessionMenuActionListenerChan := make(chan menu.CurrentSessionState)
	go updateSessionStateInLoop(sessionMenuActionListenerChan)

	scriptsLister := worker.NewScriptLister()

	taskSessionWorker := worker.NewTaskSession(scriptsLister)
	sessionManagerMenu = menu.NewSessionManager(taskSessionWorker, sessionMenuActionListenerChan)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "UnRepeater",
	})
	menuet.App().Children = getMenuItems(scriptsLister)
	menuet.App().Label = "github.com/krishnakumar4a4/un-repeater"
}
