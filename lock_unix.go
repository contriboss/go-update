//go:build !windows

package update

import (
	"os"
	"syscall"
)

// fileLock represents a file-based lock for coordinating updates across processes
type fileLock struct {
	path string
	file *os.File
}

// newFileLock creates a new file lock at the given path
func newFileLock(path string) *fileLock {
	return &fileLock{path: path + ".lock"}
}

// Lock acquires an exclusive lock, blocking until available
func (l *fileLock) Lock() error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	l.file = f

	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

// TryLock attempts to acquire the lock without blocking
// Returns true if lock was acquired, false if already held by another process
func (l *fileLock) TryLock() (bool, error) {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return false, err
	}
	l.file = f

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = f.Close()
		l.file = nil
		if err == syscall.EWOULDBLOCK {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Unlock releases the lock
func (l *fileLock) Unlock() error {
	if l.file == nil {
		return nil
	}
	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}
	if err := l.file.Close(); err != nil {
		return err
	}
	l.file = nil
	_ = os.Remove(l.path)
	return nil
}
