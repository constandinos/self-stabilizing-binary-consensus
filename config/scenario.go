package config

import "self-stabilizing-binary-consensus/logger"

var (
	ByzantineScenario string

	byzantine_scenarios = map[int]string{
		0: "NORMAL",    // Normal execution
		1: "IDLE",      // Byzantine processes remain idle (send nothing)
		2: "INVERSE",   // Byzantine processes send inverse values from the ones it should send to the servers
		3: "HALF&HALF", // Byzantine processes send different messages to half the servers
		4: "RANDOM",
	}

	corruption []bool
)

func InitializeScenario(s int) {
	if s >= len(byzantine_scenarios) {
		logger.ErrLogger.Println("Scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	ByzantineScenario = byzantine_scenarios[s]
}
