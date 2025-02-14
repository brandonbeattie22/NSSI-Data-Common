package gosharedmemory

import (
	"slices"
	"testing"
)

func TestStoreRecall(t *testing.T) {
	testDataPath := "MySharedData"
	for _, testBytes := range [][]byte{
		[]byte("Hello, World!"),
	} {
		sharedMem, err := CreateSharedMemory(testDataPath)
		if err != nil {
			t.Error(err)
		} else {
			if err := sharedMem.StoreData(testBytes); err != nil {
				t.Error(err)
			} else {
				if resultBytes, err := sharedMem.RecallData(); err != nil {
					t.Error(err)
				} else {
					if !slices.Equal(resultBytes, testBytes) {
						t.Error("test bytes do not equal result bytes")
					}
				}
			}
		}
		sharedMem.Close()
	}
}
