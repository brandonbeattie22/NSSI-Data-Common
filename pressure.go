package gosharedmemory

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	PRESSURE_PCB_DATA_CSV_HEADERS = "unix_sec,teensy_on_ms,ps1_pressure_mbar,ps1_temp_c,ps1_pressure_raw,ps1_temp_raw,ps1_error,ps2_pressure_mbar,ps2_temp_c,ps2_pressure_raw,ps2_temp_raw,ps2_error," +
		"ps4_pressure_bar,ps4_temp_c,ps4_pressure_raw,ps4_temp_raw,ps4_error,ps5_pressure_bar,ps5_temp_c,ps5_pressure_raw,ps5_temp_raw,ps5_error,htuTemp,htuHum,j1_analog_mv,j1_analog_raw"
	PRESSURE_PCB_DATA_CSV_VALUE_COUNT = 26
)

type PressureSensorData struct {
	PressureMbar float32 `json:"pressure_mbar"`
	Temp         float32 `json:"temp"`
	PressureRaw  float32 `json:"pressure_raw"`
	Error        uint8   `json:"error"`
}

func (d *PressureSensorData) pointers() []interface{} {
	return []interface{}{
		&d.PressureMbar, &d.Temp, &d.PressureRaw, &d.Error,
	}
}

func (d PressureSensorData) csv() string {
	return fmt.Sprintf("%f,%f,%f,%d", d.PressureMbar, d.Temp, d.PressureRaw, d.Error)
}

type HtuData struct {
	Temp float32 `json:"temp"`
	Hum  float32 `json:"hum"`
}

func (d *HtuData) pointers() []interface{} {
	return []interface{}{
		&d.Temp, &d.Hum,
	}
}

func (d HtuData) csv() string {
	return fmt.Sprintf("%f,%f", d.Temp, d.Hum)
}

type AnalogData struct {
	Mv  float32 `json:"mv"`
	Raw uint32  `json:"raw"`
}

func (d *AnalogData) pointers() []interface{} {
	return []interface{}{&d.Mv, &d.Raw}
}

func (d AnalogData) csv() string {
	return fmt.Sprintf("%f,%d", d.Mv, d.Raw)
}

type PressurePcbData struct {
	UnixSec uint32             `json:"unix_sec"`
	Ps1     PressureSensorData `json:"ps1"`
	Ps2     PressureSensorData `json:"ps2"`
	Ps3     PressureSensorData `json:"ps3"`
	Ps4     PressureSensorData `json:"ps4"`
	Ps5     PressureSensorData `json:"ps5"`
	Htu     HtuData            `json:"htu"`
	J1      AnalogData         `json:"j1"`
}

func (d *PressurePcbData) pointers() []interface{} {
	returnSlice := []interface{}{
		&d.UnixSec,
	}
	returnSlice = append(returnSlice, d.Ps1.pointers()...)
	returnSlice = append(returnSlice, d.Ps2.pointers()...)
	returnSlice = append(returnSlice, d.Ps4.pointers()...)
	returnSlice = append(returnSlice, d.Ps5.pointers()...)
	returnSlice = append(returnSlice, d.Htu.pointers()...)
	returnSlice = append(returnSlice, d.J1.pointers()...)
	return returnSlice
}

func (d PressurePcbData) Csv() string {
	return fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s", d.UnixSec, d.Ps1.csv(), d.Ps2.csv(), d.Ps4.csv(), d.Ps5.csv(), d.Htu.csv(), d.J1.csv())
}

func (d *PressurePcbData) FromCsv(str string) error {
	splitStr := strings.Split(str, ",")
	if n := len(splitStr); n != PRESSURE_PCB_DATA_CSV_VALUE_COUNT {
		return fmt.Errorf("expected %d values from csv, but got %d", PRESSURE_PCB_DATA_CSV_VALUE_COUNT, n)
	}

	for pIdx, pData := range d.pointers() {
		if err := readStringIntoPointer(splitStr[pIdx], pData); err != nil {
			return fmt.Errorf("error reading item num %d: %v", pIdx, err)
		}
	}
	return nil
}

func (d PressurePcbData) Json() ([]byte, error) {
	return json.Marshal(&d)
}

func (d *PressurePcbData) FromJson(buf []byte) error {
	return json.Unmarshal(buf, d)
}

/* utility */
func readStringIntoPointer(s string, p interface{}) error {
	switch v := p.(type) {
	case *float32:
		if val, err := strconv.ParseFloat(s, 32); err != nil {
			return fmt.Errorf("error parsing float32: %v", err)
		} else {
			*v = float32(val)
		}
	case *float64:
		if val, err := strconv.ParseFloat(s, 64); err != nil {
			return fmt.Errorf("error parsing float64: %v", err)
		} else {
			*v = val
		}
	case *uint8:
		if val, err := strconv.ParseUint(s, 10, 8); err != nil {
			return fmt.Errorf("error parsing uint8: %v", err)
		} else {
			*v = uint8(val)
		}
	case *uint16:
		if val, err := strconv.ParseUint(s, 10, 16); err != nil {
			return fmt.Errorf("error parsing uint16: %v", err)
		} else {
			*v = uint16(val)
		}
	case *uint32:
		if val, err := strconv.ParseUint(s, 10, 32); err != nil {
			return fmt.Errorf("error parsing uint32: %v", err)
		} else {
			*v = uint32(val)
		}
	case *uint64:
		if val, err := strconv.ParseUint(s, 10, 64); err != nil {
			return fmt.Errorf("error parsing uint64: %v", err)
		} else {
			*v = val
		}
	default:
		return fmt.Errorf("pointer type not dealt with: %T", v)
	}
	return nil
}

func (d *PressurePcbData) CreateStore(storageName string) (SharedMemory, error) {
	sharedMem, err := CreateSharedMemory(storageName)
	if err != nil {
		return sharedMem, fmt.Errorf("failed to create shared memory: %v", err)
	}
	return sharedMem, d.Store(sharedMem)
}

func (d *PressurePcbData) Store(sharedMem SharedMemory) error {
	jsonBytes, err := d.Json()
	if err != nil {
		return fmt.Errorf("failed to json-ify: %v", err)
	}

	if err := sharedMem.StoreData(jsonBytes); err != nil {
		return fmt.Errorf("failed to store data to shared memory: %v", err)
	}
	return nil
}

func (d *PressurePcbData) Recall(storageName string) error {
	sharedMem, err := AccessSharedMemory(storageName)
	if err != nil {
		return fmt.Errorf("failed to access shared memory: %d", err)
	}
	dataBytes, err := sharedMem.RecallData()
	if err != nil {
		return fmt.Errorf("failed to recall data from shared memory: %v", err)
	}
	if err := json.Unmarshal(dataBytes, d); err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return nil
}
