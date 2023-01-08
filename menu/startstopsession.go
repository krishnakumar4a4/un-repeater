package menu

import (
	"unrepeater/worker"

	"github.com/caseymrm/menuet"
)

type SessionStates string

const (
	StartSession string = "Start Session"
	StopSession  string = "Stop Session"
)

type StartStopSession struct {
	workerSession *worker.Session
}

func NewStartStopSession(worker *worker.Session) *StartStopSession {
	return &StartStopSession{
		workerSession: worker,
	}
}

func (ss *StartStopSession) StartSessionMenuItem() menuet.MenuItem {
	return menuet.MenuItem{
		Text: StartSession,
		Clicked: func() {
			removeStartMenuItem()
			ss.addStopMenuItem()
			ss.workerSession.Start()
		},
	}
}

func removeStartMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == StartSession {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *StartStopSession) addStopMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == StopSession {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	updatedItems = append(updatedItems, ss.stopSessionMenuItem())
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *StartStopSession) stopSessionMenuItem() menuet.MenuItem {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == StartSession {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }

	return menuet.MenuItem{
		Text: StopSession,
		Clicked: func() {
			removeStopMenuItem()
			ss.addStartMenuItem()
			ss.workerSession.Stop()
		},
	}
}

func removeStopMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == StopSession {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}

func (ss *StartStopSession) addStartMenuItem() {
	items := menuet.App().Children()
	updatedItems := []menuet.MenuItem{}
	for _, item := range items {
		if item.Text == StopSession {
			continue
		}
		updatedItems = append(updatedItems, item)
	}
	updatedItems = append(updatedItems, ss.StartSessionMenuItem())
	menuet.App().Children = func() []menuet.MenuItem { return updatedItems }
}