package modules

import (
	"bytes"
	"encoding/gob"
	"math"
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/messenger"
	"self-stabilizing-binary-consensus/types"
	"self-stabilizing-binary-consensus/variables"
	"strconv"
	"sync"
	"time"
)

var (
	binValues         []int
	mutex             = sync.RWMutex{}
	bc_decision_timer time.Time
	bc_decided        bool = false
)

// BinaryConsensus - The method that is called to initiate the BC module
func BinaryConsensus(bcid int, v uint) {

	// start timer for aglorithm decision time
	bc_decision_timer = time.Now()

	// (est, r) ← (v, 0);
	est := v
	r := 0

	for {
		// initialization
		binValues = make([]int, 2)
		binValues[0] = 0
		binValues[1] = 0

		// r ← r + 1;
		r += 1
		if variables.Debug {
			logger.OutLogger.Println("round=" + strconv.Itoa(r))
		}
		id := ComputeUniqueIdentifier(bcid, r)

		// bvBroadcast EST[r](est);
		go BvBroadcast(id, est)

		// wait(binValues[r] != ∅);
		for {
			mutex.Lock()
			if size(binValues) > 0 {
				mutex.Unlock()
				break
			}
			mutex.Unlock()
		}

		// initializations
		aux_counter := [2]int{0, 0}
		vals := make([]int, 2)
		vals[0] = 0
		vals[1] = 0

		// broadcast AUX[r](w) where w ∈ binValues[r];
		w := get_a_value(binValues)
		broadcast("AUX", types.NewBcMessage(id, uint(w)))
		aux_counter[w] += 1

		// wait ∃ a set of binary values, vals, and a set of (n−t) messages AUX[r](x), such
		// that vals is the set union of the values, x, carried by these (n−t) messages ∧
		// vals ⊆ binValues[r];

		// receive aux messages
		if _, in := messenger.BcChannel[id]; !in {
			messenger.BcChannel[id] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}
		for message := range messenger.BcChannel[id] {
			j := message.From
			vJ := message.BcMessage.Value

			if variables.Debug {
				logger.OutLogger.Println("RECEIVE AUX", "j="+strconv.Itoa(j), "v="+strconv.Itoa(int(vJ)))
			}

			aux_counter[vJ] += 1

			if contains(binValues, 0) && (aux_counter[0] > 0) && contains(binValues, 1) && (aux_counter[1] > 0) {
				if (aux_counter[0] + aux_counter[1]) >= (variables.N - variables.F) {
					vals[0] = 1
					vals[1] = 1
					break
				}
			} else if contains(binValues, 0) && (aux_counter[0] > 0) {
				if aux_counter[0] >= (variables.N - variables.F) {
					vals[0] = 1
					break
				}
			} else if contains(binValues, 1) && (aux_counter[1] > 0) {
				if aux_counter[1] >= (variables.N - variables.F) {
					vals[1] = 1
					break
				}
			}
		}

		if variables.Debug {
			logger.OutLogger.Println("vals=" + arr2set(vals))
		}

		// s[r] ← randomBit();
		s := uint(randomBit(r))

		if variables.Debug {
			logger.OutLogger.Println("randomBit(" + strconv.Itoa(r) + ")=" + strconv.Itoa(int(s)))
		}

		// if (vals = {v}) then
		if size(vals) == 1 {
			// if (v = s[r]) then
			v := uint(get_a_value(vals))
			if v == s {
				// decide(v) if not yet done
				if !bc_decided {
					decide_bc(id, v)
					bc_decided = true
				}
			}
			// est ← v;
			est = v
		} else {
			// else est ← s[r];
			est = s
		}

		// expected termination round
		if r == variables.M {
			return
		}
	}
}

func BvBroadcast(identifier int, v uint) {
	// initializations
	bval_counter := [2]int{0, 0}

	// do broadcast bVAL(v)
	broadcasted := [2]bool{false, false}
	broadcast("EST", types.NewBcMessage(identifier, v))
	broadcasted[v] = true
	bval_counter[v] += 1

	// receive bVAL messages
	if _, in := messenger.BvbChannel[identifier]; !in {
		messenger.BvbChannel[identifier] = make(chan struct {
			BcMessage types.BcMessage
			From      int
		})
	}
	for message := range messenger.BvbChannel[identifier] {
		j := message.From
		tag := message.BcMessage.Tag
		vJ := message.BcMessage.Value

		if variables.Debug {
			logger.OutLogger.Println("RECEIVE EST", "j="+strconv.Itoa(j), "v="+strconv.Itoa(int(vJ)))
		}

		bval_counter[vJ] += 1

		// if (bVAL(vJ) received from (t + 1) different processors and bVAL(vJ) not yet broadcast)
		if (bval_counter[vJ] >= (variables.F + 1)) && (!broadcasted[vJ]) {
			// broadcast bVAL(vJ)
			broadcast("EST", types.NewBcMessage(tag, vJ))
			broadcasted[vJ] = true
			bval_counter[vJ] += 1
		}

		// if (bVAL(vJ) received from (2t + 1) different processors)
		if bval_counter[vJ] >= (2*variables.F + 1) {
			// binValues ← binValues ∪ {vJ}
			mutex.Lock()
			binValues[vJ] = 1
			mutex.Unlock()
			if variables.Debug {
				logger.OutLogger.Println("binValues=" + arr2set(binValues))
			}
		}
	}
}

func broadcast(tag string, bcMessage types.BcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(bcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	if tag == "EST" {
		messenger.Broadcast(types.NewMessage(w.Bytes(), "BVB"))
	} else if tag == "AUX" {
		messenger.Broadcast(types.NewMessage(w.Bytes(), "BC"))
	}

	if variables.Debug {
		logger.OutLogger.Println("BROADCAST", tag, "v="+strconv.Itoa(int(bcMessage.Value)))
	}
}

func decide_bc(id int, v uint) {
	//BCAnswer[id] <- v
	duration := float64(time.Since(bc_decision_timer).Seconds())
	duration = math.Round(duration*100) / 100
	logger.OutLogger.Println("stats<byzantine,decision_time,messages,decision>:", variables.Byzantine, duration,
		variables.ReceivingMessages, v)
	if variables.Debug {
		logger.OutLogger.Println("decision=" + strconv.Itoa(int(v)))
	}
}

/* -------------------------------- Helper Functions -------------------------------- */

// ComputeUniqueIdentifier - Creates a unique num from (bcid,round) pair (Cantor's pairing func)
func ComputeUniqueIdentifier(a int, b int) int {
	res := (a * a) + (3 * a) + (2 * a * b) + b + (b * b)
	res = res / 2
	return res
}
