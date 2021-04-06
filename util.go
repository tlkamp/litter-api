package api

import (
	"math"
	"strconv"
)

func getBool(i interface{}) bool {
	switch unk := i.(type) {
	case bool:
		return unk
	case string:
		b, err := strconv.ParseBool(unk)
		if err != nil {
			return false
		}
		return b
	default:
		return false
	}
}

func getFloat(i interface{}) float64 {
	switch i := i.(type) {
	case string:
		f, err := strconv.ParseFloat(i, 32)
		if err != nil {
			return math.NaN()
		}
		return f
	case int:
		return float64(i)
	case float32:
		return float64(i)
	case float64:
		return i
	default:
		return math.NaN()
	}
}

func hexToFloat(s string) float64 {
	f, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return math.NaN()
	}
	return float64(f)
}

// NewState creates a struct to hold the status of the litterbox
func NewState(r robotResponse) State {
	s := State{
		LitterRobotID:             r.LitterRobotID.(string),
		LitterRobotSerial:         r.LitterRobotSerial.(string),
		Name:                      r.LitterRobotNickname.(string),
		PowerStatus:               r.PowerStatus.(string),
		UnitStatus:                statusMap[r.UnitStatus.(string)],
		CycleCount:                getFloat(r.CycleCount),
		CycleCapacity:             getFloat(r.CycleCapacity),
		CyclesUntilFull:           getFloat(r.CyclesUntilFull),
		CyclesAfterDrawerFull:     getFloat(r.CyclesAfterDrawerFull),
		DFICycleCount:             getFloat(r.DFICycleCount),
		CleanCycleWaitTimeMinutes: hexToFloat(r.CleanCycleWaitTimeMinutes.(string)),
		PanelLockActive:           getBool(r.PanelLockActive),
		NightLightActive:          getBool(r.NightLightActive),
		SleepModeActive:           getBool(r.SleepModeActive),
		DFITriggered:              getBool(r.IsDFITriggered),
	}
	return s
}
