#!/bin/bash

N=12
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
SELF_STABILIZING=1
CORRUPTION=1
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 40 80 100 150 200 250 400 480)

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

