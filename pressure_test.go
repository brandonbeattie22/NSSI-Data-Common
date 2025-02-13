package gosharedmemory

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
