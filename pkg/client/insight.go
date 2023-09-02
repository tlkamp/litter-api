package client

type History struct {
	Date            string `json:"date"`
	CyclesCompleted int    `json:"cyclesCompleted"`
}

// Insight - represents a Litter Robot Insights response.
type Insight struct {
	AverageCycles float64   `json:"averageCycles"`
	TotalCycles   int       `json:"totalCycles"`
	CycleHistory  []History `json:"cycleHistory"`
}
