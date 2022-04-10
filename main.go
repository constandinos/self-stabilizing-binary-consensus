package main

import (
	"log"
	"os"
	"os/signal"
	"self-stabilizing-binary-consensus/config"
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/messenger"
	"self-stabilizing-binary-consensus/modules"
	"self-stabilizing-binary-consensus/threshenc"
	"self-stabilizing-binary-consensus/variables"
	"strconv"
	"syscall"
)

// initializer(id, n, m, clients, remote, byzantine_scenario, corrupt_scenario, binValue)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, m int, clients int, rem int, byzantine_scenario int, self_stabilization_flag bool,
	corrupt_scenario int) {
	variables.Initialize(id, n, m, clients, rem)
	logger.InitializeLogger("./logs/out/", "./logs/error/")

	config.InitializeByzantineScenario(byzantine_scenario)
	config.InitializeCorruptionScenario(corrupt_scenario)

	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}

	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F, " | M:", variables.M, " | Clients:",
		variables.Clients, " | Byzantine scenario:", config.ByzantineScenario, " | Byzantine processor:", variables.Byzantine,
		" | Self-stabilizing:", self_stabilization_flag, " | Corruption scenario:", config.CorruptionScenario, " | Remote:",
		variables.Remote, "\n\n",
	)

	threshenc.ReadKeys("./keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()
	messenger.TransmitMessages()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for range terminate {
			for i := 0; i < variables.N; i++ {
				if i == variables.ID {
					continue // Not myself
				}
				messenger.ReceiveSockets[i].Close()
				messenger.SendSockets[i].Close()
			}

			for i := 0; i < variables.Clients; i++ {
				messenger.ServerSockets[i].Close()
				messenger.ResponseSockets[i].Close()
			}
			os.Exit(0)
		}
	}()
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && string(args[0]) == "generate_keys" {
		N, _ := strconv.Atoi(args[1])
		threshenc.GenerateKeys(N, "./keys/")

	} else if len(args) == 9 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		m, _ := strconv.Atoi(args[2])
		clients, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])
		byzantine_scenario, _ := strconv.Atoi(args[5])
		corruption_scenario, _ := strconv.Atoi(args[6])
		self_stabilization, _ := strconv.Atoi(args[7])
		binValue, _ := strconv.Atoi(args[8])

		var self_stabilization_flag bool
		if self_stabilization == 0 {
			self_stabilization_flag = false
		} else {
			self_stabilization_flag = true
		}

		initializer(id, n, m, clients, remote, byzantine_scenario, self_stabilization_flag, corruption_scenario)

		logger.OutLogger.Println("Initial estimate value: ", uint(binValue))

		if self_stabilization == 0 {
			modules.BinaryConsensus(1, uint(binValue))
		} else {
			modules.SelfStabilizingBinaryConsensus(int(binValue))
			//modules.SelfStabilizingMultivaluedConsensus(int(binValue))
		}

		done := make(chan interface{}) // To keep the server running
		<-done

	} else {
		log.Fatal("Arguments should be '<ID> <N> <Clients> <Scenario> <Remote>'")
	}
}
