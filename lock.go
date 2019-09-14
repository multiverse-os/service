package service

import (
	"syscall"
)

func lockFile(fd uintptr) error   { return syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB) }
func unlockFile(fd uintptr) error { return syscall.Flock(int(fd), syscall.LOCK_UN) }
