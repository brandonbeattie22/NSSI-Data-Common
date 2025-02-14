package nssidatacommon

import (
	"slices"
	"testing"
	"time"
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

func TestStoreRecallProc(t *testing.T) {
	testDataPath := "MySharedData"
	for _, testBytes := range [][]byte{
		[]byte("Hello, World!"),
	} {
		sharedMem, err := CreateSharedMemory(testDataPath)
		if err != nil {
			t.Errorf("failed to create shared memory for test: %v", err)
		} else {
			if err := sharedMem.StoreData(testBytes); err != nil {
				t.Errorf("failed to store into shared memory for test: %v", err)
			} else {
				resultChan := make(chan []byte, 1)
				go func() {
					mem, err := AccessSharedMemory(testDataPath)
					if err != nil {
						resultChan <- nil
						return
					}

					res, err := mem.RecallData()
					if err != nil {
						resultChan <- nil
						return
					}
					resultChan <- res
				}()
				select {
				case resultBytes := <-resultChan:
					if resultBytes == nil {
						t.Error("recieved no bytes")
						continue
					} else {
						if !slices.Equal(resultBytes, testBytes) {
							t.Error("result bytes are not equal")
						}
					}
				case <-time.After(time.Second * 2):
					t.Error("operation timed out")
				}
			}
		}
		sharedMem.Close()
	}
}

func TestStoreRecallData(t *testing.T) {

}
