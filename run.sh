#!/bin/bash

N=10
CLIENTS=1
REM=0
SCEN=0

go install self-stabilizing-binary-consensus

self-stabilizing-binary-consensus generate_keys $N

for (( ID=0; ID<$N; ID++ ))
do
	self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM $(($ID%2)) &
done

