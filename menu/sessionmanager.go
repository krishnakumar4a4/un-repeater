package menu

import (
	"github.com/krishnakumar4a4/un-repeater/worker"

	"github.com/caseymrm/menuet"
)

type SessionStates string

const (
	StartSession SessionStates = "Start Session"
	StopSession  SessionStates = "Stop Session"
)

type CurrentSessionState int

const (
	SessionInProgress CurrentSessionState = iota + 1
	SessionStopped
)

type SessionManager struct {
	workerSession     *worker.TaskSession
	sessionNotifyChan chan CurrentSessionState
}

func NewSessionManager(worker *worker.TaskSession, sessionStateListenerChan chan CurrentSessionState) *SessionManager {
	return &SessionManager{
		workerSession:     worker,
		sessionNotifyChan: sessionStateListenerChan,
	}
}

func (ss *SessionManager) StartSessionMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text: string(StartSession),
		Clicked: func() {
			removeStartMenuItem()
			ss.addStopMenuItem()
			ss.workerSession.Start()
			ss.sessionNotifyChan <- SessionInProgress
		},
	}
}

func removeStartMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == string(StartSession) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *SessionManager) addStopMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == string(StopSession) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	updatedItems = append(updatedItems, ss.stopSessionMenuItem())
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *SessionManager) stopSessionMenuItem() menuet.MenuItem {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == string(StartSession) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }

	return menuet.MenuItem{
		Text: string(StopSession),
		Clicked: func() {
			removeStopMenuItem()
			ss.workerSession.Stop()
			ss.addStartMenuItem()
			ss.sessionNotifyChan <- SessionStopped
		},
	}
}

func removeStopMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == string(StopSession) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *SessionManager) addStartMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == string(StopSession) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	updatedItems = append(updatedItems, ss.StartSessionMenuItem())
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}
