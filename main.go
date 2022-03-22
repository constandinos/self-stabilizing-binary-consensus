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

// Initializer - Method that initializes all required processes
func initializer(id int, n int, clients int, scenario int, rem int) {
	variables.Initialize(id, n, clients, rem)
	logger.InitializeLogger("./logs/out/", "./logs/error/")

	config.InitializeScenario(scenario)
	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}

	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F, " | Clients:",
		variables.Clients, " | Scenario:", config.Scenario, " | Remote:", variables.Remote, "\n\n",
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

	} else if len(args) == 6 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])
		binValue, _ := strconv.Atoi(args[5])

		initializer(id, n, clients, scenario, remote)

		//modules.BvBroadcast(1, 0)
		logger.OutLogger.Println("Initial estimate value: ", uint(binValue))
		//modules.BinaryConsensus(1, uint(binValue))
		modules.SelfStabilizingMultivaluedConsensus(int(binValue))

		done := make(chan interface{}) // To keep the server running
		<-done

	} else {
		log.Fatal("Arguments should be '<ID> <N> <Clients> <Scenario> <Remote>'")
	}
}
