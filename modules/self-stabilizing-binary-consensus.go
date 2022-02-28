package modules

import (
	"container/list"
	"self-stabilizing-binary-consensus/variables"
	//"sync"
)

var (
	est map[int][]*list.List
	aux map[int][]*list.List
	//mutex = sync.RWMutex{}
)

func SelfStabilizingBinaryConsensus(bcid int, initVal uint) {
	// Initialization
	N := variables.N
	M := variables.M
	F := variables.F
	ID := variables.ID

	// initState := (0, [[∅, . . . , ∅], . . . , [∅, . . . , ∅]], [[⊥, . . . , ⊥], . . . , [⊥, . . . , ⊥]])
	// (r, est, aux) ← initState

	// r = 0
	r := 0

	// est = [[∅, . . . , ∅], . . . , [∅, . . . , ∅]]
	est = make(map[int][]*list.List, M+1)
	for i := 0; i < M+1; i++ {
		est[i] = make([]*list.List, N)
		for j := 0; j < N; j++ {
			est[i][j] = list.New()
		}
	}

	// aux = [[⊥, . . . , ⊥], . . . , [⊥, . . . , ⊥]]
	aux = make(map[int][]*list.List, M+1)
	for i := 0; i < M+1; i++ {
		aux[i] = make([]*list.List, N)
		for j := 0; j < N; j++ {
			aux[i][j] = list.New()
		}
	}

	// est[0][i] ← {v}
	append_val(est[0][ID], initVal)

	// do forever begin
	for {
		// if ((r, est, aux) != initState) then
		if (r != 0) || !check_empty(est, M+1, N) || !check_empty(aux, M+1, N) {
			// r ← min{r+1, M +1 };
			r = min(r+1, M+1)

			// repeat
			for {
				// if (est[0][i] != {v}) then est[0][i] ← {w} : ∃w ∈ est[0][i];
				if size(est[0][ID]) > 1 {
					// est[0][i] ← {w} : ∃w ∈ est[0][i];
					est[0][ID].Remove(est[0][ID].Back())
				}

				// foreach r' ∈ {1, . . . , r−1}
				for i := 1; i <= r-1; i++ {
					// if est[r'][i] = ∅ ∨ aux[r'][i] = ⊥
					if (size(est[i][ID]) == 0) || (size(aux[i][ID]) == 0) {
						// (est[r'][i], aux[r'][i]) ← (est[0][i], x) : x ∈ est[0][i];
						// est[r'][i] = est[0][i]
						clear(est[i][ID])
						for element := est[0][ID].Front(); element != nil; element = element.Next() {
							est[i][ID].PushBack(element.Value)
						}
						// aux[r'][i] = x : x ∈ est[0][i]
						clear(aux[i][ID])
						aux[i][ID].PushBack(est[0][ID].Front().Value)
					}
				}

				// if ((∃w ∈ binValues(r, 2t+1) ∧ (aux[r][i] = ⊥ ∨ aux[r][i] ¬∈ binValues(r, 2t+1)))
				binValues := bin_values(r, 2*F+1, ID, N, est)
				if (size(binValues) > 0) && ((size(aux[r][ID]) == 0) || !subset(aux[r][ID], binValues)) {
					clear(aux[r][ID])
					aux[r][ID].PushBack(binValues.Front().Value)
				}
			}

		}

	}

}

// append appends a new value val in the end of the list if val is not already in the list l.
func append_val(l *list.List, val uint) {
	for element := l.Front(); element != nil; element = element.Next() {
		if element.Value == val {
			return
		}
	}
	l.PushBack(val)
}

// size returns the number of elements of list l.
func size(l *list.List) int {
	return l.Len()
}

// check_empty checks if a given array arr is empty.
func check_empty(arr map[int][]*list.List, m int, n int) bool {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if size(arr[i][j]) > 0 {
				return false
			}
		}
	}
	return true
}

// min finds the minimun number of x and y
func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// clear remove all elements from the list l
func clear(l *list.List) {
	for element := l.Front(); element != nil; element = element.Next() {
		l.Remove(element)
	}
}

// exist checks if value v exist in the list l
func exist(l *list.List, v uint) bool {
	for element := l.Front(); element != nil; element = element.Next() {
		if element.Value == v {
			return true
		}
	}
	return false
}

// count counts the number of occurrences of variable v in the round r
func count(r int, v uint, id int, n int, a map[int][]*list.List) int {
	counter := 0
	for j := 0; j < n; j++ {
		if j != id {
			if exist(a[r][j], v) {
				counter++
			}
		}
	}
	return counter
}

// binValues creates a set of values that appeared at least x times
func bin_values(r int, x int, id int, n int, a map[int][]*list.List) *list.List {
	l := list.New()
	if count(r, 0, id, n, a) >= x {
		append_val(l, 0)
	}
	if count(r, 1, id, n, a) >= x {
		append_val(l, 1)
	}
	return l
}

// subset checks if list l1 is a subset of list l2
func subset(l1 *list.List, l2 *list.List) bool {
	for e1 := l1.Front(); e1 != nil; e1 = e1.Next() {
		flag := false
		for e2 := l2.Front(); e2 != nil; e2 = e2.Next() {
			if e1.Value == e2.Value {
				flag = true
			}
		}
		if !flag {
			return false
		}
	}
	return true
}
