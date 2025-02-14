package gosharedmemory

const SHARED_DATA_BUFF_SIZE = uint32(1024)

type SharedMemory interface {
	StoreData([]byte) error
	RecallData() ([]byte, error)
	Close()
}

func CreateSharedMemory(memoryName string) (SharedMemory, error) {
	return createSharedMemory(memoryName)
}
