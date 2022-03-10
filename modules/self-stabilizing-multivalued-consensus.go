package modules

import (
	"fmt"
	"self-stabilizing-binary-consensus/logger"
	"time"
)

func SelfStabilizingMultivaluedConsensus(binVal int) {
	// Call self-stabilizing binary consensus
	go SelfStabilizingBinaryConsensus(binVal)
	time.Sleep(1 * time.Second)

	for {
		bin_consensus := result()
		// Decide
		if bin_consensus >= 0 {
			logger.OutLogger.Println("Successful decide")
			fmt.Println(bin_consensus)
			break
			// Transient error
		} else if bin_consensus == -1 {
			logger.OutLogger.Println("Transient error")
			break
			// No value was decided
		} else if bin_consensus == -2 {
			time.Sleep(1 * time.Second)
		}
	}
}
