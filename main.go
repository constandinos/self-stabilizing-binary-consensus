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

func cleanup() {
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

// Initializer - Method that initializes all required processes
func initializer(id int, n int, m int, clients int, remote int, byzantine_scenario int, self_stabilizing int, corruption int,
	debug int, optimization int) {
	variables.Initialize(id, n, m, clients, remote, debug, optimization)
	logger.InitializeLogger("./logs/out/", "./logs/error/")

	config.InitializeByzantineScenario(byzantine_scenario)
	config.InitializeCorruptionScenario(corruption)

	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}

	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F, " | M:", variables.M, " | Clients:",
		variables.Clients, " | Byzantine scenario:", config.ByzantineScenario, " | Byzantine processor:", variables.Byzantine,
		" | Self-Stabilization:", self_stabilizing == 1, " | Corruption scenario:", config.CorruptionScenario, " | Remote:",
		variables.Remote, " | Debug:", variables.Debug, " | Optimization:", variables.Optimization, "\n\n",
	)

	threshenc.ReadKeys("./keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()

	// IDLE Attack Scenario
	if (config.ByzantineScenario == "IDLE") && (variables.Byzantine) {
		logger.OutLogger.Println("IDLE ATTACK")
		return
	}

	messenger.TransmitMessages()

	cleanup()
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && string(args[0]) == "generate_keys" {
		N, _ := strconv.Atoi(args[1])
		threshenc.GenerateKeys(N, "./keys/")

	} else if len(args) == 11 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		m, _ := strconv.Atoi(args[2])
		clients, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])
		byzantine_scenario, _ := strconv.Atoi(args[5])
		self_stabilizing, _ := strconv.Atoi(args[6])
		corruption, _ := strconv.Atoi(args[7])
		debug, _ := strconv.Atoi(args[8])
		optimization, _ := strconv.Atoi(args[9])
		init_value, _ := strconv.Atoi(args[10])

		initializer(id, n, m, clients, remote, byzantine_scenario, self_stabilizing, corruption, debug, optimization)

		if variables.Debug {
			logger.OutLogger.Print("Initial estimate value:", init_value, "\n\n")
		}

		if self_stabilizing == 0 {
			modules.BinaryConsensus(1, uint(init_value))
		} else {
			modules.SelfStabilizingBinaryConsensus(1, int(init_value))
		}

		done := make(chan interface{}) // To keep the server running
		<-done

	} else {
		log.Fatal("go run main.go <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> " +
			"<Corruption> <Debug> <Optimization> <Initial value>")
	}
}
