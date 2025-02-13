package gosharedmemory

import (
	"slices"
	"testing"

	"golang.org/x/sys/windows"
)

func TestStoreRecall(t *testing.T) {
	testDataPath := "MySharedData"
	for _, testBytes := range [][]byte{
		[]byte("Hello, World!"),
	} {
		mutex, memHandle, err := CreateSharedMemory(testDataPath)
		if err != nil {
			t.Errorf("error in `CreateSharedMemory()`: %v", err)
		} else {
			if err := StoreData(testBytes, memHandle, mutex); err != nil {
				t.Errorf("error in `StoreData()`: %v", err)
			} else {
				if resultByte, err := RecallData(testDataPath); err != nil {
					t.Errorf("error in `RecallData()`: %v", err)
				} else {
					if !slices.Equal(resultByte, testBytes) {
						t.Errorf("Original bytes not the same after `StoreData()` and `RecallData()`, started with \"%s\", ended up with \"%s\"",
							testBytes, resultByte)
					}
				}
			}
		}
		windows.CloseHandle(mutex)
		windows.CloseHandle(memHandle)

	}
}
