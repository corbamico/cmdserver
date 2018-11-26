package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	processName string = "shadowsocks_local"
)

func findAndKillProcess(path string, info os.FileInfo, err error) error {
	pattern := `/proc/([0-9]+)/status`
	reg := regexp.MustCompile(pattern)

	if err != nil {
		return nil
	}

	//we only walk for /proc/<pid>/status
	if pids := reg.FindStringSubmatch(path); pids != nil {
		pid, err := strconv.Atoi(pids[0])
		if err != nil {
			return nil
		}
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		name := string(f[6:bytes.IndexByte(f, '\n')])
		// we find it
		if name == processName {
			proc, err := os.FindProcess(pid)
			if err == nil {
				proc.Kill()
			}
			return io.EOF
		}
	}

	return nil
}

//StopProc kill the process of "processName"
func StopProc() {
	filepath.Walk("/proc", findAndKillProcess)
}

//RunProc run the process of "processName"
func RunProc() {
	cmd := exec.Command(processName)
	if err := cmd.Start(); err == nil {
		time.Sleep(500 * time.Millisecond)
		cmd.Process.Release()
	}
}

//RestartProc re-start the process
func RestartProc() {
	StopProc()
	time.Sleep(500 * time.Millisecond)
	RunProc()
}

var statProcess int32 // 0:unknwon,1:running
//StatProc return stat of process
func StatProc() int32 {
	statProcess = 0
	filepath.Walk("/proc", findAndStatProcess)
	return statProcess
}

func findAndStatProcess(path string, info os.FileInfo, err error) error {
	pattern := `/proc/([0-9]+)/status`
	reg := regexp.MustCompile(pattern)

	if err != nil {
		return nil
	}

	//we only walk for /proc/<pid>/status
	if pids := reg.FindStringSubmatch(path); pids != nil {
		//pid, err := strconv.Atoi(pids[0])
		if err != nil {
			return nil
		}
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		name := string(f[6:bytes.IndexByte(f, '\n')])
		// we find it
		if name == processName {
			statProcess = 1 //running
			return io.EOF
		}
	}

	return nil
}
