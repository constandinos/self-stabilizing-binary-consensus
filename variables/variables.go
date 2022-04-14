package variables

import "sync"

var (
	// ID - This processor's id.
	ID int

	// N - Number of processors
	N int

	// M - Number of iterations
	M int

	// F - Number of faulty processors
	F int

	// Byzantine - If the processor is byzantine or not
	Byzantine bool

	// Clients - Size of Clients Set
	Clients int

	// Remote - If we are running locally or remotely
	Remote bool

	// DEFAULT - The default value that is used in the algorithms
	DEFAULT []byte

	// Server metrics regarding the experiment evaluation
	MsgComplexity int
	MsgSize       int64
	MsgMutex      sync.RWMutex

	TotalRCVMessages int
	//Sleep_time       int
)

// Initialize - Variables initializer method
func Initialize(id int, n int, m int, c int, rem int) {
	ID = id
	N = n
	M = m
	F = (N - 1) / 3

	Clients = c

	if rem == 1 {
		Remote = true
	} else {
		Remote = false
	}

	DEFAULT = []byte("")

	MsgComplexity = 0
	MsgSize = 0
	MsgMutex = sync.RWMutex{}
	TotalRCVMessages = 0
}
