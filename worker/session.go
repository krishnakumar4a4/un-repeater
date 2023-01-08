package worker

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"sync"
)

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

	entries, err := os.ReadDir(path.Join("scripts"))
	if err != nil {
		log.Fatalf("unable to open scripts dir: %s", err.Error())
		return
	}
	fileNames := make([]string, 0, len(entries))
	for _, elem := range entries {
		fName := elem.Name()
		fileNames = append(fileNames, fName)
	}
	sort.Strings(fileNames)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to current working dir: %s", err.Error())
		return
	}
	scriptsDir := filepath.Join(cwd, "scripts")
	for _, fName := range fileNames {
		scriptPath := filepath.Join(scriptsDir, fName)
		cmd := exec.Command(scriptPath)
		wg := &sync.WaitGroup{}
		stdOutReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("unable to open stdout pipe: %s", err.Error())
			return
		}
		wg.Add(1)
		go logStdOut(wg, stdOutReader)

		stdErrReader, err := cmd.StderrPipe()
		if err != nil {
			log.Fatalf("unable to open stderr pipe: %s", err.Error())
			return
		}
		wg.Add(1)
		go logStdErr(wg, stdErrReader)

		if err := cmd.Run(); err != nil {
			log.Fatalf("unable to execute script: %s", err.Error())
		}
	}
}

func (s *Session) Stop() {
	fmt.Println("session stopped")
}

// TODO: Can read line by line
func logStdOut(wg *sync.WaitGroup, readCloser io.ReadCloser) {
	defer readCloser.Close()
	defer wg.Done()
	outData, err := io.ReadAll(readCloser)
	if err != nil {
		log.Fatalf("unable to read script stdout: %s", err.Error())
		return
	}
	log.Println("STDOUT: " + string(outData))
}

// TODO: Can read line by line
func logStdErr(wg *sync.WaitGroup, readCloser io.ReadCloser) {
	defer readCloser.Close()
	defer wg.Done()
	outData, err := io.ReadAll(readCloser)
	if err != nil {
		log.Fatalf("unable to read script stderr: %s", err.Error())
		return
	}
	log.Println("STDERR: " + string(outData))
}
