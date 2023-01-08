package worker

import "fmt"

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) Start() {
	// Start capture session
	// Collect bug reports
	// Get BLE snoop from report
	// Get Scenario description
	// Save/discard capture
	// Also save logs of needed

	fmt.Println("session started")
}

func (s *Session) Stop() {
	fmt.Println("session stopped")
}
