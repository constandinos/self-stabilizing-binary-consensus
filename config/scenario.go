package config

import "self-stabilizing-binary-consensus/logger"

var (
	Scenario string

	scenarios = map[int]string{
		0: "NORMAL",      // Normal execution
		1: "IDLE",        // Byzantine processes remain idle (send nothing)
		2: "INVERSE",     // Byzantine processes send inverse values from the ones it should send to the servers
		3: "HALF_&_HALF", // Byzantine processes send different messages to half the servers
	}
)

func InitializeScenario(s int) {
	if s >= len(scenarios) {
		logger.ErrLogger.Println("Scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	Scenario = scenarios[s]
}
