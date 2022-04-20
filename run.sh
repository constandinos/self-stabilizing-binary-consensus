#!/bin/bash

N=12
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
SELF_STABILIZING=1
CORRUPTION=0
DEBUG=1
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 40 100 150 200 250 350 450 500)

go install self-stabilizing-binary-consensus
self-stabilizing-binary-consensus generate_keys $N

for (( ID=0; ID<$N; ID++ ))
do
	# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Initial value>
	self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} $(( $ID%2 )) &
done

#counter=$(ps -u const | grep "self" | wc -l)
#while [ $counter -gt 0 ]; do
#	sleep 1
#	counter=$(ps -u const | grep "self" | wc -l)
#done
#echo "end"
#grep "stats" logs/out/*.log

