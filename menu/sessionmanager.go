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
	SessionStopInProgress
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
			ss.workerSession.Start()
			ss.sessionNotifyChan <- SessionInProgress
		},
	}
}

func (ss *SessionManager) StopSessionMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text: string(StopSession),
		Clicked: func() {
			ss.sessionNotifyChan <- SessionStopInProgress
			ss.workerSession.Stop()
			ss.sessionNotifyChan <- SessionStopped
		},
	}
}
