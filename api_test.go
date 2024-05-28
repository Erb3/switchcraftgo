package switchcraftgo

import (
	"testing"
)

func TestGetPlayerCounts(t *testing.T) {
	res, err := GetPlayerCounts()

	if err != nil {
		t.Fatalf("GetPlayerCounts() returned error %s", err.Error())
	}

	if res.Active > res.Total {
		t.Fatalf("GetPlayerCounts() returned more active players than total players. Received %d active out of %d total", res.Active, res.Total)
	}

	if res.Active < 0 {
		t.Fatalf("GetPlayerCounts() retured active below zero (received %d)", res.Active)
	}

	if res.Total < 0 {
		t.Fatalf("GetPlayerCounts() retured total below zero (received %d)", res.Total)
	}
}

func TestGetTps(t *testing.T) {
	res, err := GetTps()

	if err != nil {
		t.Fatalf("GetTps() returned error %s", err.Error())
	}

	if res.AverageTps == 0 {
		t.Fatalf("GetTps() returned TPS of zero, less than zero, or nil (got %f)", res.AverageTps)
	}

	if res.AverageMillisecondPerTick == 0 {
		t.Fatalf("GetTps() returned average milliseconds of zero, less than zero, or nil (got %f)", res.AverageMillisecondPerTick)
	}

	if res.MillisecondsLastTick == 0 {
		t.Fatalf("GetTps() returned last milliseconds of zero, less than zero, or nil (got %f)", res.MillisecondsLastTick)
	}
}
