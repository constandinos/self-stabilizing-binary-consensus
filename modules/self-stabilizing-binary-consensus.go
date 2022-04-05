package modules

import (
	"bytes"
	"encoding/gob"
	"self-stabilizing-binary-consensus/config"
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/messenger"
	"self-stabilizing-binary-consensus/types"
	"self-stabilizing-binary-consensus/variables"
	"strconv"
	"sync"
	"time"
)

var (
	est       [][][]int
	aux       [][][]int
	mutex_est [][]sync.RWMutex
	mutex_aux [][]sync.RWMutex
	N         int
	M         int
	F         int
	ID        int
	r         int
)

/* Operations */

// propose(v)
func SelfStabilizingBinaryConsensus(v int) {
	// Initialization
	N = variables.N
	M = variables.M
	F = variables.F
	ID = variables.ID

	// Create mutex
	mutex_est = make([][]sync.RWMutex, M+2)
	mutex_aux = make([][]sync.RWMutex, M+2)
	for i := 0; i < M+2; i++ {
		mutex_est[i] = make([]sync.RWMutex, N)
		mutex_aux[i] = make([]sync.RWMutex, N)
		for j := 0; j < N; j++ {
			mutex_est[i][j] = sync.RWMutex{}
			mutex_aux[i][j] = sync.RWMutex{}
		}
	}

	/// initState := (0, [[∅, . . . , ∅], . . . , [∅, . . . , ∅]], [[⊥, . . . , ⊥], . . . , [⊥, . . . , ⊥]])
	/// (r, est, aux) ← initState

	// r = 0
	r = 0

	// est = [[∅, . . . , ∅], . . . , [∅, . . . , ∅]]
	// aux = [[⊥, . . . , ⊥], . . . , [⊥, . . . , ⊥]]
	est = make([][][]int, M+2)
	aux = make([][][]int, M+2)
	for i := 0; i < M+2; i++ {
		est[i] = make([][]int, N)
		aux[i] = make([][]int, N)
		for j := 0; j < N; j++ {
			est[i][j] = make([]int, 2)
			clear(est[i][j])
			aux[i][j] = make([]int, 2)
			clear(aux[i][j])
		}
	}

	// est[0][i] ← {v}
	append_val(est[0][ID], v)

	// do forever begin
	for {
		// if ((r, est, aux) != initState) then
		if (r != 0) || !check_empty(est) || !check_empty(aux) {
			// loop counter
			repeat := 0

			// r ← min{r+1, M+1};
			r = min(r+1, M+1)

			// [Unit test 3] Corruption of r: the variable r suddenly changes value
			/*if config.Corruptions[3] && r == 2 {
				logger.OutLogger.Println("CORRUPTION 3 r=" + strconv.Itoa(r) + "->3")
				r = 3
				config.Corruptions[3] = false
			}*/

			// Receive income messages in background communication
			go receive(r)

			// repeat
			for {
				// Increase loop counter
				repeat++

				// Debugging
				logger.OutLogger.Println("round="+strconv.Itoa(r), "repeat="+strconv.Itoa(repeat))

				// [Unit test 0] Corruption of est[0][ID]: the variable est[0][ID] loses its value
				if config.Corruptions[0] && r == 1 {
					mutex_est[0][ID].Lock()
					logger.OutLogger.Println("CORRUPTION 0 est[0][" + strconv.Itoa(ID) + "]=" + arr2set(est[0][ID]) + "->{}")
					clear(est[0][ID])
					mutex_est[0][ID].Unlock()
					config.Corruptions[0] = false
				}

				// [Unit test 1] Corruption of est[0][ID]: setting variable est[0][ID] with {0,1}
				if config.Corruptions[1] && r == 2 {
					mutex_est[0][ID].Lock()
					logger.OutLogger.Println("CORRUPTION 1 est[0][" + strconv.Itoa(ID) + "]=" + arr2set(est[0][ID]) + "->{0,1}")
					append_val(est[0][ID], 0)
					append_val(est[0][ID], 1)
					mutex_est[0][ID].Unlock()
					config.Corruptions[1] = false
				}

				// if (est[0][i] != {v})
				mutex_est[0][ID].Lock()
				if size(est[0][ID]) == 0 {
					// est[0][i] ← {v}
					append_val(est[0][ID], v)
					logger.OutLogger.Println("FIXED CORRUPTION 0 est[0][" + strconv.Itoa(ID) + "]={}->" + arr2set(est[0][ID]))
				} else if size(est[0][ID]) > 1 {
					// est[0][i] ← {w} : ∃w ∈ est[0][i];
					w := get_a_value(est[0][ID])
					set_val(est[0][ID], w)
					logger.OutLogger.Println("FIXED CORRUPTION 1 est[0][" + strconv.Itoa(ID) + "]={0,1}->" + arr2set(est[0][ID]))
				}
				mutex_est[0][ID].Unlock()

				// [Unit test 2] Corruption of est[r'][ID] and aux[r'][ID]: the variables est[r'][ID] and aux[r'][ID]
				// lose their value
				if config.Corruptions[2] && r == 2 {
					for rr := 1; rr <= r-1; rr++ {
						mutex_est[rr][ID].Lock()
						mutex_aux[rr][ID].Lock()
						logger.OutLogger.Println("CORRUPTION 2 est[" + strconv.Itoa(rr) + "][" + strconv.Itoa(ID) + "]=" +
							arr2set(est[rr][ID]) + "->{} aux[" + strconv.Itoa(rr) + "][" + strconv.Itoa(ID) + "]=" +
							arr2set(aux[rr][ID]) + "->{}")
						clear(est[rr][ID])
						clear(aux[rr][ID])
						mutex_est[rr][ID].Unlock()
						mutex_aux[rr][ID].Unlock()
					}
					config.Corruptions[2] = false
				}

				// foreach r' ∈ {1, . . . , r−1}
				for rr := 1; rr <= r-1; rr++ {
					mutex_est[rr][ID].Lock()
					mutex_aux[rr][ID].Lock()
					// if est[r'][i] = ∅ ∨ aux[r'][i] = ⊥
					if (size(est[rr][ID]) == 0) || (size(aux[rr][ID]) == 0) {
						logger.OutLogger.Println("FIXED CORRUPTION 2/3 est[" + strconv.Itoa(rr) + "][" + strconv.Itoa(ID) +
							"]=" + arr2set(est[rr][ID]) + "->" + arr2set(est[0][ID]) + " aux[" + strconv.Itoa(rr) + "][" +
							strconv.Itoa(ID) + "]=" + arr2set(aux[rr][ID]) + "->{" + strconv.Itoa(get_a_value(est[0][ID])) +
							"}")
						/// (est[r'][i], aux[r'][i]) ← (est[0][i], x) : x ∈ est[0][i];
						// est[r'][i] = est[0][i]
						set(est[rr][ID], est[0][ID])
						// aux[r'][i] = x : x ∈ est[0][i]
						x := get_a_value(est[0][ID])
						set_val(aux[rr][ID], x)
					}
					mutex_est[rr][ID].Unlock()
					mutex_aux[rr][ID].Unlock()
				}

				// if ((∃w ∈ binValues(r, 2t+1) ∧ (aux[r][i] = ⊥ ∨ aux[r][i] ¬∈ binValues(r, 2t+1)))
				binValues := bin_values(r, 2*F+1)
				w := get_a_value(binValues)

				mutex_aux[r][ID].Lock()
				if (w != -1) && ((size(aux[r][ID]) == 0) || !contains(binValues, get_a_value(aux[r][ID]))) {
					// aux[r][i] ← w
					set_val(aux[r][ID], w)
					// Debugging
					//logger.OutLogger.Println("aux[" + strconv.Itoa(r) + "][" + strconv.Itoa(ID) + "]=" + strconv.Itoa(w))
				}
				mutex_aux[r][ID].Unlock()

				// foreach p j ∈ P do send EST(True, r, est[r−1][i] ∪ binValues(r, t+1), aux[r][i])
				binValues = bin_values(r, F+1)
				new_est := union(est[r-1][ID], binValues)
				// [Unit test 4] Corruption of SEND message: the variable r suddenly changes value
				if config.Corruptions[4] && size(aux[r][ID]) == 1 {
					logger.OutLogger.Println("CORRUPTION 4 SEND true, " + strconv.Itoa(r) + ", " + arr2set(new_est) + ", " +
						arr2set(new_est) + "->{}")
					send("EST", types.NewSSBCMessage(true, r, new_est[0], new_est[1], 0, 0))
					config.Corruptions[4] = false
				} else {
					send("EST", types.NewSSBCMessage(true, r, new_est[0], new_est[1], aux[r][ID][0], aux[r][ID][1]))
				}
				set(est[r][ID], new_est)

				// time.Sleep(1000 * time.Millisecond)
				time.Sleep(1 * time.Second)

				// until infoResult() != ∅
				infoResults := info_results()
				if size(infoResults) > 0 {
					break
				}
			}

			// tryToDecide(infoResult())
			try_to_decide(info_results())

			// if (∃w ∈ binValues(M +1, t+1)) then decide(w)
			binValues := bin_values(M+1, F+1)
			w := get_a_value(binValues)
			if w != -1 {
				decide(w)
			}
		}
	}
}

// result() do {if (est[M +1][i] = {v}) then return v else if (r ≥ M ∧ infoResult() != ∅)
// then return Ψ else return ⊥
func result() int {
	// if (est[M+1][i] = {v})
	if size(est[M+1][ID]) == 1 {
		// return v
		return get_a_value(est[M+1][ID])
		// else if (r ≥ M ∧ infoResult() != ∅)
	} else if (r >= M) && (size(info_results()) > 0) {
		// return Ψ (transient error symbol)
		return -1
	} else {
		// return ⊥ (no value was decided)
		return -2
	}
}

/* Macros */

// binValues creates a set of values that appeared at least x times in est for each j processor
// binValues(r, x) return {y ∈ {0, 1} : ∃s ⊆ P : |{p j ∈ s : y ∈ est[r][j]}| ≥ x};
func bin_values(rr int, x int) []int {
	counter := [2]int{0, 0}
	for j := 0; j < N; j++ {
		counter[0] += est[rr][j][0]
		counter[1] += est[rr][j][1]
		//logger.OutLogger.Print(arr2set(est[rr][j]))
	}
	s := make([]int, 2)
	clear(s)
	if counter[0] >= x {
		append_val(s, 0)
	}
	if counter[1] >= x {
		append_val(s, 1)
	}
	// Debugging
	logger.OutLogger.Println("binValues(" + strconv.Itoa(rr) + "," + strconv.Itoa(x) + ")=" + arr2set(s))
	return s
}

// info_results
// infoResult() do {if (∃s ⊆ P : n−t ≤ |s| ∧ (∀p j ∈ s : aux [r][j] ∈ binValues(r, 2t+1))) then
// return {aux [r][j]} p j ∈s else return ∅;}
func info_results() []int {
	counter := [2]int{0, 0}
	binValues := bin_values(r, 2*F+1)
	for j := 0; j < N; j++ {
		w := get_a_value(aux[r][j])
		if w != -1 && contains(binValues, w) {
			counter[w]++
		}
	}
	s := make([]int, 2)
	clear(s)
	if counter[0] >= (N - F) {
		append_val(s, 0)
	}
	if counter[1] >= (N - F) {
		append_val(s, 1)
	}
	// Debugging
	logger.OutLogger.Println("infoResults()=" + arr2set(s))
	return s
}

/* Functions */

// decide(x)
func decide(x int) {
	// foreach r' ∈ {r, . . . , M+1} do
	for rr := r; rr <= M+1; rr++ {
		// if (est[r'][i] = ∅ ∨ aux[r'][i] = ⊥)
		mutex_est[rr][ID].Lock()
		mutex_aux[rr][ID].Lock()
		if (size(est[rr][ID]) == 0) || (size(aux[rr][ID]) == 0) {
			/// (est[r'][i], aux[r'][i]) ← ({x}, x)
			// est[r'][i] = {x}
			set_val(est[rr][ID], x)
			//  aux[r'][i] = x
			set_val(aux[rr][ID], x)
		}
		mutex_est[rr][ID].Unlock()
		mutex_aux[rr][ID].Unlock()
	}

	// r <- M + 1
	r = M + 1

	// [Unit test 5] Corruption of est[M+1][ID] and aux[M+1][ID]: the variables est[M+1][ID] and aux[M+1][ID]
	if config.Corruptions[5] {
		logger.OutLogger.Println("CORRUPTION 5 est[" + strconv.Itoa(M+1) + "][" + strconv.Itoa(ID) + "]=" +
			arr2set(est[M+1][ID]) + "->{} aux[" + strconv.Itoa(M+1) + "][" + strconv.Itoa(ID) + "]=" +
			arr2set(aux[M+1][ID]) + "->{}")
		mutex_est[M+1][ID].Lock()
		mutex_aux[M+1][ID].Lock()
		clear(est[M+1][ID])
		clear(aux[M+1][ID])
		mutex_est[M+1][ID].Unlock()
		mutex_aux[M+1][ID].Unlock()
		config.Corruptions[5] = false
	}

	// Debugging
	logger.OutLogger.Println("decision=" + strconv.Itoa(x))
}

// tryToDecide(values)
func try_to_decide(values []int) {
	randomBit := random_bit(r)
	logger.OutLogger.Println("randomBit(" + strconv.Itoa(r) + ")=" + strconv.Itoa(randomBit))
	// if (values != {v})
	if size(values) != 1 {
		// est[r][i] ← {randomBit(r)}
		mutex_est[r][ID].Lock()
		set_val(est[r][ID], randomBit)
		mutex_est[r][ID].Unlock()
	} else {
		// est[r][i] ← {v}
		mutex_est[r][ID].Lock()
		set(est[r][ID], values)
		mutex_est[r][ID].Unlock()
		// if (v = randomBit(r))
		if values[randomBit] == 1 {
			decide(randomBit)
		}
	}
}

/* Auxiliary functions */

// clear remove all elements from the list l.
func clear(s []int) {
	s[0] = 0
	s[1] = 0
}

// append appends a new value val in the set s.
func append_val(s []int, val int) {
	s[val] = 1
}

// size returns the number of elements of set s.
func size(s []int) int {
	return s[0] + s[1]
}

// check_empty checks if a given array arr is empty.
func check_empty(arr [][][]int) bool {
	for i := 0; i < M+2; i++ {
		for j := 0; j < N; j++ {
			if size(arr[i][j]) > 0 {
				return false
			}
		}
	}
	return true
}

// min finds the minimun number of x and y.
func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// get_a_value returns a value from the set s.
func get_a_value(s []int) int {
	if s[0] == 1 {
		return 0
	} else if s[1] == 1 {
		return 1
	} else {
		return -1
	}
}

// set overwrite the values of set dst with the values of set src.
func set(dst []int, src []int) {
	dst[0] = src[0]
	dst[1] = src[1]
}

// set_val clears the set s and append the value w.
func set_val(s []int, w int) {
	clear(s)
	append_val(s, w)
}

// contains checks if the set s contains the value w.
func contains(s []int, w int) bool {
	if s[w] == 1 {
		return true
	} else {
		return false
	}
}

// union joins two sets.
func union(s1 []int, s2 []int) []int {
	s := make([]int, 2)
	clear(s)
	if (s1[0] == 1) || (s2[0] == 1) {
		s[0] = 1
	}
	if (s1[1] == 1) || (s2[1] == 1) {
		s[1] = 1
	}
	return s
}

// random_bit generate a psedo-random number
func random_bit(rr int) int {
	random_num := [7]int{0, 0, 0, 1, 0, 1, 1}
	return random_num[rr]
}

// arr2set create a string with a set
func arr2set(arr []int) string {
	if size(arr) == 0 {
		return "{}"
	} else if size(arr) == 1 {
		if arr[0] == 1 {
			return "{0}"
		} else {
			return "{1}"
		}
	} else if size(arr) == 2 {
		return "{0 1}"
	}
	return ""
}

/* Communication */

// send sends a message to pj ∈ P processors
func send(tag string, estMessage types.SSBCMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(estMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	messenger.Broadcast(types.NewMessage(w.Bytes(), tag))
}

// receive receives a message from pj ∈ P processors
func receive(id int) {
	if _, in := messenger.SSBCChannel[id]; !in {
		messenger.SSBCChannel[id] = make(chan struct {
			SSBCMessage types.SSBCMessage
			From        int
		})
	}

	for message := range messenger.SSBCChannel[id] {
		j := message.From                  // From
		aJ := message.SSBCMessage.Flag     // Flag
		rJ := message.SSBCMessage.Round    // Round
		est_0 := message.SSBCMessage.Est_0 // est[0]
		est_1 := message.SSBCMessage.Est_1 // est[1]
		vJ := make([]int, 2)
		clear(vJ)
		vJ[0] = est_0
		vJ[1] = est_1
		aux_0 := message.SSBCMessage.Aux_0 // aux[0]
		aux_1 := message.SSBCMessage.Aux_1 // aux[1]
		uJ := make([]int, 2)
		clear(uJ)
		uJ[0] = aux_0
		uJ[1] = aux_1

		// est[rJ][j] ← est[rJ][j] ∪ vJ
		mutex_est[rJ][j].Lock()
		set(est[rJ][j], union(est[rJ][j], vJ))
		mutex_est[rJ][j].Unlock()

		// aux[rJ][j] ← uJ
		mutex_aux[rJ][j].Lock()
		set(aux[rJ][j], uJ)
		mutex_aux[rJ][j].Unlock()

		// Debugging
		/*logger.OutLogger.Println("RECEIVED j="+strconv.Itoa(j), "flag="+strconv.FormatBool(aJ), "r="+strconv.Itoa(rJ),
		"est="+arr2set(vJ), "aux="+arr2set(uJ))*/
		//logger.OutLogger.Println("est["+strconv.Itoa(rJ)+"]["+strconv.Itoa(j)+"]="+arr2set(est[rJ][j]),
		//	"aux["+strconv.Itoa(rJ)+"]["+strconv.Itoa(j)+"]="+arr2set(aux[rJ][j]))

		// if aJ then
		if aJ {
			// send EST(False, rJ , est[rJ−1][i], aux[r][i]) to pj
			send("EST", types.NewSSBCMessage(false, rJ, est[rJ-1][ID][0], est[rJ-1][ID][1], aux[rJ][ID][0], aux[rJ][ID][1]))
			/* Debugging
			logger.OutLogger.Println("SEND flag=false r="+strconv.Itoa(rJ), "est="+arr2set(est[rJ-1][ID]), "aux="+arr2set(aux[rJ][ID]))*/
		}
	}
}
