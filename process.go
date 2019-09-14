package service

import (
	"fmt"
	"os"
	"syscall"
	"time"

	pid "github.com/multiverse-os/starshipyard/framework/service/pid"
	signal "github.com/multiverse-os/starshipyard/framework/service/signal"
)

// TODO: May want to add a weight concept so that certain functions will be
// esnured to be ran before others
type Process struct {
	Pid                int
	PidFile            *pid.File
	User               string
	UID                int
	GID                int
	TempDirectory      string
	UserCacheDirectory string
	Env                map[string]string
	Executable         string
	WorkingDirectory   string
	StartedAt          time.Time
	Signals            signal.Handler
	Shutdown           []func()
}

// NOTE: These functions are intended to be specific to this process, as this
// application framework expands to running multiple child processes, this will
// enable dynamic editing of shutdown sequence per child procses. (The
// functionality for removing ShutdownFunctions has not yet been added)
func (self *Process) ShutdownProcess() {
	for _, function := range self.Shutdown {
		function()
	}
}

func (self *Process) AppendToShutdown(exitFunction func()) {
	self.Shutdown = append(self.Shutdown, exitFunction)
}

func ParseProcess() *Process {
	if executable, err := os.Executable(); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to process executable:", err))
	} else {
		return &Process{
			Pid:        os.Getpid(),
			UID:        os.Getuid(),
			GID:        os.Getgid(),
			Executable: executable,
			StartedAt:  time.Now(),
		}
	}
}

//[ Method for process ]///////////////////////////////////////////////////////
// NOTE: Returns the close function so that it can be called easily added to a
// defer. This is important because since we are doing OS based locks on the
// pidfile we may need to unlock the file
func (self *Process) WritePid(path string) *pid.File {
	fmt.Println("[starship] writing pid:", self.Pid)
	if pidFile, err := pid.Write(path); err != nil {
		panic(fmt.Sprintf("[fatal error] failed to write pid:", err))
	} else {
		self.PidFile = pidFile
		return self.PidFile
	}
}

func (self *Process) CleanPid() error {
	return self.PidFile.Clean()
}

//[ General process utilities ]////////////////////////////////////////////////
func isProcessRunning(pid int) bool {
	if process, err := os.FindProcess(pid); err == nil {
		return false
	} else {
		err = process.Signal(syscall.Signal(0))
		return (err != nil)
	}
}

//func SetUserId(procAttr *syscall.SysProcAttr, uid uint32, gid uint32) {
//	procAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid, NoSetGroups: true}
//}
//
