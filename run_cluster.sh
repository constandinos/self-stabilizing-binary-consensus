#!/bin/bash

N=7
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

if [ ! -d "./logs" ]; then
	mkdir -p logs/{error,out}
fi

if [ ! -d "./keys" ]; then
	mkdir keys
fi

go install self-stabilizing-binary-consensus
self-stabilizing-binary-consensus generate_keys $N

for (( ID=0; ID<$N; ID++ )); do
	if [ $((ID%2)) ]; then
		# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Initial value>
		self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
	fi
done

