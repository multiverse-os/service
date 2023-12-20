package service

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
)

// TODO: Abstract the logic from each of the dependency free subpackages so that
// a developer can access all the necessary functionality from each package to
// setup a quality Linux service in an intuitive way and doing this only
// importing the base service package. Even if functions are essentially
// duplicated, this is fine for a simple interface to 5 subpackages

func syscallDup(oldfd int, newfd int) (err error) {
	// linux_arm64 platform doesn't have syscall.Dup2
	// so use the nearly identical syscall.Dup3 instead.
	return syscall.Dup3(oldfd, newfd, 0)
}

// [ Test Function ]////////////////////////////////////////////////////////////
func Daemonize(function func()) error {
	fmt.Println("Daemonizing")

	daemon := Process{LogFile: "/dev/stdout"}
	if child, err := daemon.Run(); err != nil {
		return err
	} else {
		if child != nil {
			return fmt.Errorf("[error] failed to create daemon process")
		}
	}

	// TODO: This would setup the pid file previously, should require user to
	// supply the PID file location or default to system default
	//defer daemon.Release()
	function()
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func (p *Process) Run() (childProcess *os.Process, err error) {
	fmt.Println("Run() command")
	if p.IsRunning() {
		if err = p.child(); err != nil {
			return nil, err
		}
	} else {
		if childProcess, err = p.parent(); err != nil {
			return nil, err
		}
	}
	return childProcess, nil
}

func (p *Process) parent() (parentProcess *os.Process, err error) {
	if err = p.prepareEnv(); err != nil {
		return
	}
	attr := &os.ProcAttr{
		Dir:   p.WorkingDirectory,
		Env:   p.Environment(),
		Files: p.files(),
		Sys: &syscall.SysProcAttr{
			//Chroot:     p.Chroot, // To bad there is no comments on this
			Credential: p.Credential,
			Setsid:     true,
		},
	}
	if parentProcess, err = os.StartProcess(p.abspath, p.Args, attr); err != nil {
		// TODO: Replace this with our pid code
		//if p.pidFile != nil {
		//	p.pidFile.Remove()
		//}
		return
	}
	p.rpipe.Close()
	encoder := json.NewEncoder(p.wpipe)
	if err = encoder.Encode(p); err != nil {
		return
	}
	_, err = fmt.Fprint(p.wpipe, "\n\n")
	return
}

func (p *Process) prepareEnv() (err error) {
	if p.abspath, err = os.Executable(); err != nil {
		return err
	}
	if len(p.Args) == 0 {
		p.Args = os.Args
	}
	mark := fmt.Sprintf("%s=%s", p.EnvVar.Key, p.EnvVar.Value)
	if len(p.Env) == 0 {
		p.EnvVars = os.Environ()
	}
	p.EnvVars = append(p.EnvVars, mark)
	return nil
}

func (p *Process) files() (processFiles []*os.File) {
	log := p.nullFile
	if p.logFile != nil {
		log = p.logFile
	}
	// TODO: Default should be nil
	processFiles = []*os.File{
		p.rpipe,    // (0) stdin
		log,        // (1) stdout
		log,        // (2) stderr
		p.nullFile, // (3) dup on fd 0 after initialization
	}
	// TODO: PID file defintion was here.
	//if p.pidFile != nil {
	//	f = append(f, p.pidFile.File) // (4) pid file
	//}

	// TODO: This is the method using our pid system
	//pid.Write(p.PIDFile)

	return processFiles
}

func (p *Process) child() error {
	if p.Initialized {
		return os.ErrInvalid
	}
	p.Initialized = true
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(p); err != nil {
		//pid.Clean(p.PIDFile)
		return err
	}
	// create PID file after context decoding to know PID file full path.
	// TODO: Previously file locking was used, but this may not actually be
	// necessary unless we scale this up to handle several processes
	// TODO: we should just have a cleanup function and a encapsulating function
	// that will run it if error, instead of calling pid.Clean more than once.
	//pid.Write(p.PIDFile)
	if err := syscall.Close(0); err != nil {
		//pid.Clean(p.PIDFile)
		return err
	}
	if err := syscallDup(3, 0); err != nil {
		//pid.Clean(p.PIDFile)
		return err
	}
	if p.Umask != 0 {
		syscall.Umask(int(p.Umask))
	}
	if len(p.Chroot) > 0 {
		if err := syscall.Chroot(p.Chroot); err != nil {
			//pid.Clean(p.PIDFile)
			return err
		}
	}
	return nil
}
