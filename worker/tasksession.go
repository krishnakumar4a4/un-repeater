package worker

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var sessionDataRoot string

func init() {
	sessionDataRoot = os.Getenv("SESSION_ROOT")
	if strings.TrimSpace(sessionDataRoot) == "" {
		sessionDataRoot = "session-data"
	}
}

type TaskSession struct {
	doneChan          chan int
	sessionFolderName string
}

func NewTaskSession() *TaskSession {
	return &TaskSession{}
}

func (s *TaskSession) Start() {
	// Start capture session
	// Collect bug reports
	// Get BLE snoop from report
	// Get Scenario description
	// Save/discard capture
	// Also save logs if needed

	log.Println("session started")
	s.sessionFolderName = filepath.Join(sessionDataRoot, time.Now().Format(time.RFC3339))
	if err := os.MkdirAll(s.sessionFolderName, os.ModePerm); err != nil {
		log.Fatalf("error creating session folder: %s", err.Error())
		return
	}
	s.doneChan = make(chan int)
	s.execScripts(path.Join("scripts", "start-hooks"))
}

func (s *TaskSession) Stop() {
	log.Println("session stopped")
	close(s.doneChan)
	s.execScripts(path.Join("scripts", "stop-hooks"))
}

func execCmd(ctx context.Context, scriptPath, sessionFolderName string) {
	cmd := exec.CommandContext(ctx, scriptPath, sessionFolderName)
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
	log.Println("Finished executing command waiting for stdout and stderr to close")
	wg.Wait()
	log.Println("Finished executing command")
}

func (s *TaskSession) execScripts(dir string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get current working dir: %s", err.Error())
		return
	}
	scriptsDir := filepath.Join(cwd, dir)

	entries, err := os.ReadDir(scriptsDir)
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
	for _, fName := range fileNames {
		scriptPath := filepath.Join(scriptsDir, fName)
		ctx := context.Background()
		go s.cancelContext(ctx)
		if strings.Contains(fName, ".noblock") {
			log.Printf("running %s as non block mode\n", scriptPath)
			go execCmd(ctx, scriptPath, s.sessionFolderName)
		} else {
			log.Printf("running %s as blocking mode\n", scriptPath)
			execCmd(ctx, scriptPath, s.sessionFolderName)
		}
	}
}

func (s *TaskSession) cancelContext(ctx context.Context) {
	for {
		select {
		case <-s.doneChan:
			log.Println("closing script context")
			ctx.Done()
			return
		}
	}
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
		log.Printf("unable to read script stdout: %s \n", err.Error())
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
		log.Printf("unable to read script stderr: %s \n", err.Error())
		return
	}
}
