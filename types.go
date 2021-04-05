package api

import (
	"github.com/go-resty/resty/v2"
	"time"
)

type loginResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
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

// State - the exported state of the Litter Robot.
type State struct {
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

// Insight - represents a Litter Robot Insights response.
type Insight struct {
	AverageCycles float64 `json:"averageCycles"`
	TotalCycles   int     `json:"totalCycles"`
}

// Config - Configuration for the Litter Robot client
type Config struct {
	ApiUrl       string
	AuthUrl      string
	ClientId     string
	ClientSecret string
	Email        string
	Password     string
	ApiKey       string
}

// Client - The Client for interacting with the Litter Robot API.
type Client struct {
	*Config
	Expiry       time.Duration
	apiClient    *resty.Client
	authClient   *resty.Client
	token        string
	refreshToken string
	userID       string
	robots       map[string]State
	statusPath   string
	insightsPath string
	cmdPath      string
}

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

