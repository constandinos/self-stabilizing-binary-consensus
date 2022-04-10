#!/bin/bash

N=12
M=5
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION_SCENARIO=0
SELF_STABILIZATION=1

go install self-stabilizing-binary-consensus

self-stabilizing-binary-consensus generate_keys $N


for (( ID=0; ID<$N; ID++ ))
do
	# id, n, m, clients, remote, byzantine_scenario, corrupt_scenario, binValue
	
	#self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $SELF_STABILIZATION $(($ID%2)) &
	self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $SELF_STABILIZATION 1 &
done


