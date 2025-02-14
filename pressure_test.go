package nssidatacommon

import (
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	if n := len(strings.Split(PRESSURE_PCB_DATA_CSV_HEADERS, ",")); n != PRESSURE_PCB_DATA_CSV_VALUE_COUNT {
		t.Errorf("expected number of comma-separated headers in `PRESSURE_PCB_DATA_CSV_HEADERS` to match `PRESSURE_PCB_DATA_CSV_VALUE_COUNT` (%d), but got %d",
			PRESSURE_PCB_DATA_CSV_VALUE_COUNT, n)
	}
}

func (d PressureSensorData) Equals(o PressureSensorData) bool {
	return (d.PressureMbar == o.PressureMbar) && (d.PressureRaw == o.PressureRaw) && (d.Temp == o.Temp) && (d.Error == o.Error)
}
func (d HtuData) Equals(o HtuData) bool {
	return (d.Hum == o.Hum) && (d.Temp == o.Temp)
}
func (d AnalogData) Equals(o AnalogData) bool {
	return (d.Mv == o.Mv) && (d.Raw == o.Raw)
}
func (d PressurePcbData) Equals(o PressurePcbData) bool {
	return (d.UnixSec == o.UnixSec) && (d.Ps1.Equals(o.Ps1)) && (d.Ps2.Equals(o.Ps2)) && (d.Ps3.Equals(o.Ps3)) && (d.Ps4.Equals(o.Ps4)) && (d.Ps5.Equals(o.Ps5)) && (d.Htu.Equals(o.Htu)) && (d.J1.Equals(o.J1))
}
func TestStoreRecallPressure(t *testing.T) {
	testData := PressurePcbData{
		UnixSec: 123456,
		Ps1: PressureSensorData{
			1.0,
			2.3,
			5.67,
			100,
		},
		Ps2: PressureSensorData{
			0.1,
			2.39,
			5.1,
			0,
		},
		Ps3: PressureSensorData{
			1.1,
			12.39,
			15.1,
			10,
		},
		Ps4: PressureSensorData{
			3.1,
			3.39,
			3.1,
			3,
		},
		Ps5: PressureSensorData{
			12,
			12,
			12,
			12,
		},
		Htu: HtuData{
			12.12,
			42.42,
		},
		J1: AnalogData{
			Raw: 43,
			Mv:  100.1,
		},
	}

	testMemoryName := "testMemory"
	sharedMem, err := testData.CreateStore(testMemoryName)
	if err != nil {
		t.Errorf("failed to create shared memory and store: %v", err)
		if sharedMem != nil {
			sharedMem.Close()
		}
		return
	}
	defer sharedMem.Close()

	resultData := &PressurePcbData{}

	if err := resultData.Recall(testMemoryName); err != nil {
		t.Errorf("failed to recall data: %v", err)
		return
	}
	if !testData.Equals(*resultData) {
		t.Error("test and result data are not equal")
	}
}
