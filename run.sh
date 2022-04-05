#!/bin/bash

N=4
M=5
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION_SCENARIO=0

go install self-stabilizing-binary-consensus

self-stabilizing-binary-consensus generate_keys $N

initValue=(0 0 1 1)

for (( ID=0; ID<$N; ID++ ))
do
	#self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM ${initValue[ID]} &
	
	# id, n, m, clients, remote, byzantine_scenario, corrupt_scenario, binValue
	self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $(($ID%2)) &
	
	#self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM 0 &
done

