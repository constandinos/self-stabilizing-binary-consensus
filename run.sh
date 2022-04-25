#!/bin/bash

N=12
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=4
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

if [ ! -d "./logs" ]; then
	mkdir -p logs/{error,out}
else
	rm -rf ./logs
	mkdir -p logs/{error,out}
fi

if [ ! -d "./keys" ]; then
	mkdir keys
else
	rm -rf ./keys
	mkdir keys
fi

go run main.go generate_keys $N

for (( ID=0; ID<$N; ID++ )); do
	# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Initial value>
	go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
done

