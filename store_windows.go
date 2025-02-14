package gosharedmemory

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type SharedMemoryWindows struct {
	MutexPath   string
	MemPath     string
	ReadOnly    bool
	MutexHandle windows.Handle
	MemHandle   windows.Handle
}

func accessSharedMemory(sharedDataName string) (*SharedMemoryWindows, error) {
	memName := "Global\\" + sharedDataName
	mutName := memName + "Mutex"

	return &SharedMemoryWindows{
		MutexPath: mutName,
		MemPath:   memName,
		ReadOnly:  true,
	}, nil
}
func createSharedMemory(sharedDataName string) (*SharedMemoryWindows, error) {
	memName := "Global\\" + sharedDataName
	mutName := memName + "Mutex"
	// Define the names for shared memory and mutex.
	sharedMemName, err := windows.UTF16PtrFromString(memName)
	if err != nil {
		return nil, err
	}

	mutexName, err := windows.UTF16PtrFromString(mutName)
	if err != nil {
		return nil, err
	}

	// Create a named mutex for synchronization.
	mutex, err := windows.CreateMutex(nil, false, mutexName)
	if (err != nil) && (err != windows.ERROR_ALREADY_EXISTS) {
		return nil, fmt.Errorf("CreateMutex failed: %v", err)
	}

	// Create a named file mapping object backed by the system paging file.
	handle, err := windows.CreateFileMapping(windows.InvalidHandle, nil,
		windows.PAGE_READWRITE, 0, SHARED_DATA_BUFF_SIZE, sharedMemName)
	if err != nil {
		return nil, fmt.Errorf("CreateFileMapping failed: %v", err)
	}

	return &SharedMemoryWindows{
		MutexPath:   mutName,
		MemPath:     memName,
		ReadOnly:    false,
		MutexHandle: mutex,
		MemHandle:   handle,
	}, nil
}

// StoreData writes data to the shared memory.
func (d *SharedMemoryWindows) StoreData(data []byte) error {
	if d.ReadOnly {
		return fmt.Errorf("readonly instance")
	}
	if n := len(data); n > int(SHARED_DATA_BUFF_SIZE) {
		return fmt.Errorf("data length %d exceeds maximum %d", n, SHARED_DATA_BUFF_SIZE)
	}

	// Map a view of the file into the address space.
	addr, err := windows.MapViewOfFile(d.MemHandle, windows.FILE_MAP_READ|windows.FILE_MAP_WRITE, 0, 0, uintptr(SHARED_DATA_BUFF_SIZE))
	if err != nil {
		return fmt.Errorf("MapViewOfFile failed: %v", err)
	}
	defer windows.UnmapViewOfFile(addr)

	// Acquire the mutex before writing to shared memory.
	waitResult, err := windows.WaitForSingleObject(d.MutexHandle, windows.INFINITE)
	if err != nil {
		return fmt.Errorf("WaitForSingleObject failed: %v", err)
	}
	if waitResult != windows.WAIT_OBJECT_0 {
		return fmt.Errorf("failed to acquire mutex")
	}
	// Ensure the mutex is released after we're done.
	defer windows.ReleaseMutex(d.MutexHandle)

	// Write data to the shared memory.
	dst := (*[SHARED_DATA_BUFF_SIZE]byte)(unsafe.Pointer(addr))
	copy(dst[:], data)
	return nil
}

// Lazy-load kernel32.dll and the OpenFileMappingW procedure.
var (
	modkernel32          = windows.NewLazySystemDLL("kernel32.dll")
	procOpenFileMappingW = modkernel32.NewProc("OpenFileMappingW")
)

// OpenFileMapping wraps the OpenFileMappingW Win32 API call.
func OpenFileMapping(desiredAccess uint32, inheritHandle bool, name *uint16) (windows.Handle, error) {
	var bInherit uint32
	if inheritHandle {
		bInherit = 1
	}
	ret, _, err := procOpenFileMappingW.Call(
		uintptr(desiredAccess),
		uintptr(bInherit),
		uintptr(unsafe.Pointer(name)),
	)
	if ret == 0 {
		// If the call fails, err will be non-nil.
		if err != nil && err != syscall.Errno(0) {
			return 0, err
		}
		return 0, syscall.EINVAL
	}
	return windows.Handle(ret), nil
}

// RecallData reads data from the shared memory and prints it.
func (d *SharedMemoryWindows) RecallData() ([]byte, error) {
	// Define the names for shared memory and mutex.
	sharedMemName, err := windows.UTF16PtrFromString(d.MemPath)
	if err != nil {
		return nil, err
	}

	mutexName, err := windows.UTF16PtrFromString(d.MutexPath)
	if err != nil {
		return nil, err
	}

	// Open the named mutex.
	mutex, err := windows.CreateMutex(nil, false, mutexName)
	if err != nil && err != windows.ERROR_ALREADY_EXISTS {
		return nil, fmt.Errorf("CreateMutex failed: %v", err)
	}

	// Acquire the mutex before reading.
	waitResult, err := windows.WaitForSingleObject(mutex, windows.INFINITE)
	if err != nil {
		return nil, fmt.Errorf("WaitForSingleObject failed: %v", err)
	}
	if waitResult != windows.WAIT_OBJECT_0 {
		return nil, fmt.Errorf("failed to acquire mutex")
	}
	defer windows.ReleaseMutex(mutex)

	// Open the existing file mapping object.
	handle, err := OpenFileMapping(windows.FILE_MAP_READ, false, sharedMemName)
	if err != nil {
		return nil, fmt.Errorf("OpenFileMapping failed: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Map a view of the file into the address space.
	addr, err := windows.MapViewOfFile(handle, windows.FILE_MAP_READ, 0, 0, uintptr(SHARED_DATA_BUFF_SIZE))
	if err != nil {
		return nil, fmt.Errorf("MapViewOfFile failed: %v", err)
	}
	defer windows.UnmapViewOfFile(addr)

	// Read data from the shared memory.
	src := (*[SHARED_DATA_BUFF_SIZE]byte)(unsafe.Pointer(addr))
	// Determine the length of the data by finding the first null byte.
	var dataLen int
	for i, b := range src[:] {
		if b == 0 {
			dataLen = i
			break
		}
	}
	data := make([]byte, dataLen)
	copy(data, src[:dataLen])
	return data, nil
}

func (d *SharedMemoryWindows) Close() {
	if d.ReadOnly {
		return
	}
	println("Close called!")
	windows.CloseHandle(d.MemHandle)
	windows.CloseHandle(d.MutexHandle)
}
