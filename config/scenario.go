package config

import (
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/variables"
)

var (
	ByzantineScenario string

	byzantine_scenarios = map[int]string{
		0: "NORMAL",    // Normal execution
		1: "IDLE",      // Byzantine processor remain idle / send nothing (crash)
		2: "INVERSE",   // Byzantine processor send inverse values from the ones it should send to the servers
		3: "HALF&HALF", // Byzantine processor send different messages to half the servers
		4: "RANDOM",    // Byzantine processor send random messages to different processors
	}

	CorruptionScenario string

	corruption_scenarios = map[int]string{
		0: "NORMAL", // Normal execution
		1: "RANDOM", // Apply random corruptions
	}
)

func InitializeByzantineScenario(s int) {
	if s >= len(byzantine_scenarios) {
		logger.ErrLogger.Println("Byzantine scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	ByzantineScenario = byzantine_scenarios[s]

	if ByzantineScenario == "NORMAL" {
		variables.Byzantine = false
	} else {
		if variables.ID < variables.F {
			variables.Byzantine = true
		} else {
			variables.Byzantine = false
		}
	}
}

func InitializeCorruptionScenario(s int) {
	if s >= len(corruption_scenarios) {
		logger.ErrLogger.Println("Corruption scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	CorruptionScenario = corruption_scenarios[s]
}
