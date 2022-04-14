#!/bin/bash

N=12
M=5
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION_SCENARIO=0
SELF_STABILIZATION=1

SLEEP_TIME=(130)

for t in ${SLEEP_TIME[@]}; do
	go install self-stabilizing-binary-consensus
	self-stabilizing-binary-consensus generate_keys $N	
	for (( ID=0; ID<$N; ID++ )); do
		self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $SELF_STABILIZATION 1 &
	done
	sleep 4
	sh ./kill.sh
done

grep "stats" logs/out/*.log | awk '{time[$8]+=$6; count[$8]+=1} END{for (i in time){print i, time[i]/count[i]}}' | sort -n -k2
rm logs/out/*.log
