package worker

import (
	"bufio"
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
	// Also save logs if needed

	log.Println("session started")

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
		log.Fatalf("unable to get current working dir: %s", err.Error())
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
	log.Println("session stopped")
}

func logStdOut(wg *sync.WaitGroup, readCloser io.ReadCloser) {
	defer readCloser.Close()
	defer wg.Done()

	fileScanner := bufio.NewScanner(readCloser)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		log.Println("STDOUT: " + fileScanner.Text())
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("unable to read script stdout: %s", err.Error())
		return
	}
}

func logStdErr(wg *sync.WaitGroup, readCloser io.ReadCloser) {
	defer readCloser.Close()
	defer wg.Done()

	fileScanner := bufio.NewScanner(readCloser)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		log.Println("STDERR: " + fileScanner.Text())
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("unable to read script stderr: %s", err.Error())
		return
	}
}
