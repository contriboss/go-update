//go:build windows

package update

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

const (
	lockfileExclusiveLock   = 0x02
	lockfileFailImmediately = 0x01
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

	var overlapped syscall.Overlapped
	r1, _, err := procLockFileEx.Call(
		f.Fd(),
		lockfileExclusiveLock,
		0,
		1, 0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

// TryLock attempts to acquire the lock without blocking
// Returns true if lock was acquired, false if already held by another process
func (l *fileLock) TryLock() (bool, error) {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return false, err
	}
	l.file = f

	var overlapped syscall.Overlapped
	r1, _, err := procLockFileEx.Call(
		f.Fd(),
		lockfileExclusiveLock|lockfileFailImmediately,
		0,
		1, 0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if r1 == 0 {
		f.Close()
		l.file = nil
		// ERROR_LOCK_VIOLATION = 33
		if errno, ok := err.(syscall.Errno); ok && errno == 33 {
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

	var overlapped syscall.Overlapped
	r1, _, err := procUnlockFileEx.Call(
		l.file.Fd(),
		0,
		1, 0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if r1 == 0 {
		return err
	}

	if err := l.file.Close(); err != nil {
		return err
	}
	l.file = nil
	_ = os.Remove(l.path)
	return nil
}
