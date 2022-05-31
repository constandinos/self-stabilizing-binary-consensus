package variables

import (
	"math/rand"
	"sync"
)

var (
	// This processor's id.
	ID int

	// Number of processors
	N int

	// Number of iterations
	M int

	// Number of faulty processors
	F int

	// If the processor is byzantine or not
	Byzantine bool

	// Size of Clients Set
	Clients int

	// If we are running locally or remotely
	Remote bool

	// If the logger will print all debug messages
	Debug bool

	// Optimization flag
	Optimization bool

	// Counter for receiving messages
	ReceivingMessages int

	// Random number generator
	RandomGenerator            *rand.Rand
	RandomGeneratorCorruptions *rand.Rand

	/* From Vasilis implentation */
	// DEFAULT - The default value that is used in the algorithms
	DEFAULT []byte

	// Server metrics regarding the experiment evaluation
	MsgComplexity int
	MsgSize       int64
	MsgMutex      sync.RWMutex
)

// Initialize - Variables initializer method
func Initialize(id int, n int, m int, clients int, remote int, debug int, optimization int) {
	ID = id

	N = n
	F = (N - 1) / 3

	M = m

	Clients = clients

	if remote == 1 {
		Remote = true
	} else {
		Remote = false
	}

	if debug == 1 {
		Debug = true
	} else {
		Debug = false
	}

	if optimization == 1 {
		Optimization = true
	} else {
		Optimization = false
	}

	ReceivingMessages = 0

	RandomGenerator = rand.New(rand.NewSource(int64(ID)))
	RandomGeneratorCorruptions = rand.New(rand.NewSource(int64(ID)))

	/* From Vasilis implentation */
	DEFAULT = []byte("")

	MsgComplexity = 0
	MsgSize = 0
	MsgMutex = sync.RWMutex{}
}
