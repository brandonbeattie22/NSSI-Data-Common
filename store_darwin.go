package nssidatacommon

import (
	"fmt"
	"io"
	"os"
)

const SHARED_TEENSY_DATA_1 = "Global\\Teensy1SharedMemory"

type SharedMemoryDarwin struct {
	filepath string
}

func accessSharedMemory(sharedDataName string) (SharedMemoryDarwin, error) {
	return SharedMemoryDarwin{
		filepath: "/var/run/" + sharedDataName,
	}, nil
}

func createSharedMemory(sharedDataName string) (SharedMemoryDarwin, error) {
	filepath := "/var/run/" + sharedDataName
	if f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0677); err != nil {
		return SharedMemoryDarwin{}, fmt.Errorf("failed to create file: %v", err)
	} else {
		f.Close()
	}

	return SharedMemoryDarwin{filepath: filepath}, nil
}

// StoreData writes data to the shared memory.
func (d SharedMemoryDarwin) StoreData(data []byte) error {
	f, err := os.OpenFile(d.filepath, os.O_TRUNC|os.O_RDWR, 0677)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	} else if m := len(data); n != m {
		return fmt.Errorf("failed to write correct number of bytes to file, wrote %d bytes, data was %d bytes long", n, m)
	}

	return nil
}

func (d SharedMemoryDarwin) RecallData() ([]byte, error) {
	f, err := os.Open(d.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file due to error: %v", err)
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (d SharedMemoryDarwin) Close() {
}
