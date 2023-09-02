package client

import (
	"github.com/tlkamp/litter-api/internal/util"
)

var status = []string{
	"RDY",
	"CCP",
	"CCC",
	"CSF",
	"DF1",
	"DF2",
	"CST",
	"CSI",
	"BR",
	"P",
	"OFF",
	"SDF",
	"DFS",
}

var statusMap map[string]float64

func init() {
	statusMap = make(map[string]float64, len(status))
	for i, v := range status {
		statusMap[v] = float64(i)
	}
}

type robotResponse struct {
	AutoOfflineDisabled       interface{} `json:"autoOfflineDisabled"`
	CleanCycleWaitTimeMinutes interface{} `json:"cleanCycleWaitTimeMinutes"`
	CyclesAfterDrawerFull     interface{} `json:"cyclesAfterDrawerFull"`
	CycleCapacity             interface{} `json:"cycleCapacity"`
	CycleCount                interface{} `json:"cycleCount"`
	CyclesUntilFull           interface{} `json:"cyclesUntilFull"`
	DeviceType                interface{} `json:"deviceType"`
	DFICycleCount             interface{} `json:"DFICycleCount"`
	DidNotifyOffline          interface{} `json:"didNotifyOffline"`
	IsDFITriggered            interface{} `json:"isDFITriggered"`
	IsOnboarded               interface{} `json:"isOnboarded"`
	LastSeen                  interface{} `json:"lastSeen"`
	LitterRobotID             interface{} `json:"litterRobotId"`
	LitterRobotNickname       interface{} `json:"litterRobotNickname"`
	LitterRobotSerial         interface{} `json:"litterRobotSerial"`
	NightLightActive          interface{} `json:"nightLightActive"`
	PanelLockActive           interface{} `json:"panelLockActive"`
	PowerStatus               interface{} `json:"powerStatus"`
	SetupDate                 interface{} `json:"setupDate"`
	SleepModeActive           interface{} `json:"sleepModeActive"`
	SleepModeEndTime          interface{} `json:"sleepModeEndTime"`
	SleepModeStartTime        interface{} `json:"sleepModeStartTime"`
	UnitStatus                interface{} `json:"unitStatus"`
}

// Robot is the exported state of the LitterRobot.
type Robot struct {
	CleanCycleWaitTimeMinutes float64
	CyclesAfterDrawerFull     float64
	CycleCapacity             float64
	CycleCount                float64
	CyclesUntilFull           float64
	DidNotifyOffline          bool
	DFICycleCount             float64
	DFITriggered              bool
	LitterRobotID             string
	LitterRobotSerial         string
	Name                      string
	NightLightActive          bool
	PanelLockActive           bool
	PowerStatus               string
	SleepModeActive           bool
	UnitStatus                float64
}

func newRobot(r robotResponse) Robot {
	s := Robot{
		LitterRobotID:             r.LitterRobotID.(string),
		LitterRobotSerial:         r.LitterRobotSerial.(string),
		Name:                      r.LitterRobotNickname.(string),
		PowerStatus:               r.PowerStatus.(string),
		UnitStatus:                statusMap[r.UnitStatus.(string)],
		CycleCount:                util.Float(r.CycleCount),
		CycleCapacity:             util.Float(r.CycleCapacity),
		CyclesUntilFull:           util.Float(r.CyclesUntilFull),
		CyclesAfterDrawerFull:     util.Float(r.CyclesAfterDrawerFull),
		DFICycleCount:             util.Float(r.DFICycleCount),
		CleanCycleWaitTimeMinutes: util.HexToFloat(r.CleanCycleWaitTimeMinutes.(string)),
		PanelLockActive:           util.Bool(r.PanelLockActive),
		NightLightActive:          util.Bool(r.NightLightActive),
		SleepModeActive:           util.Bool(r.SleepModeActive),
		DFITriggered:              util.Bool(r.IsDFITriggered),
	}
	return s
}
