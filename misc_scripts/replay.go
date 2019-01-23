package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	REPLAY_DIR = "replays_local"
	FLUORINE_DIR = "C:\\Users\\Owner\\github\\fluorine"
)

func main() {

	files, _ := ioutil.ReadDir(REPLAY_DIR)

	if len(files) == 0 {
		return
	}

	var latest_time time.Time
	var latest_filename string

	for _, file := range files {
		if file.ModTime().After(latest_time) {
			latest_time = file.ModTime()
			latest_filename = file.Name()
		}
	}

	full_name := filepath.Join(REPLAY_DIR, latest_filename)

	cmd := fmt.Sprintf("electron %s -o %s", FLUORINE_DIR, full_name)

	cmd_split := strings.Fields(cmd)

	exec_command := exec.Command(cmd_split[0], cmd_split[1:]...)
	exec_command.Start()
}
