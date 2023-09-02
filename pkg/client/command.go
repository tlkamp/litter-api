package client

const (
	on            = "1"
	off           = "0"
	powerCmd      = "<P"
	cycleCmd      = "<C"
	nightLightCmd = "<N"
	panelLockCmd  = "<L"
	waitCmd       = "<W"
)

type commandBody struct {
	Command       string `json:"command"`
	LitterRobotId string `json:"litterRobotId"`
}
