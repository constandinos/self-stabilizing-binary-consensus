package config

import (
	"math/rand"
	"self-stabilizing-binary-consensus/logger"
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
		1: "RANDOM",
		2: "ALL",
	}

	Corruptions      []bool
	corruption_cases int = 6
)

func InitializeByzantineScenario(s int) {
	if s >= len(byzantine_scenarios) {
		logger.ErrLogger.Println("Byzantine scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	ByzantineScenario = byzantine_scenarios[s]
}

func InitializeCorruptionScenario(s int) {
	if s >= len(corruption_scenarios) {
		logger.ErrLogger.Println("Corruption scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	CorruptionScenario = corruption_scenarios[s]

	Corruptions = make([]bool, corruption_cases)
	if CorruptionScenario == "NORMAL" {
		for i := 0; i < corruption_cases; i++ {
			Corruptions[i] = false
		}
	} else if CorruptionScenario == "RANDOM" {
		for i := 0; i < corruption_cases; i++ {
			rand_num := rand.Intn(2)
			if rand_num == 0 {
				Corruptions[i] = false
			} else {
				Corruptions[i] = true
			}
		}
	} else {
		for i := 0; i < corruption_cases; i++ {
			Corruptions[i] = true
		}
	}
}
