#!/bin/bash

N=4
CLIENTS=1
REM=0
SCEN=3

go install self-stabilizing-binary-consensus

self-stabilizing-binary-consensus generate_keys $N

initValue=(0 0 1 1)

for (( ID=0; ID<$N; ID++ ))
do
	#self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM ${initValue[ID]} &
	self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM $(($ID%2)) &
	#self-stabilizing-binary-consensus $ID $N $CLIENTS $SCEN $REM 0 &
done

