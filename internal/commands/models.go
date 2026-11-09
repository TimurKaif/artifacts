package commands

type CharacterInfo struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	HP   int    `json:"hp"`
}

type Cooldown struct {
	TotalSeconds     int    `json:"total_seconds"`
	RemainingSeconds int    `json:"remaining_seconds"`
	Reason           string `json:"reason"`
}
