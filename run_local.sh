#!/bin/bash

N=12
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
SELF_STABILIZING=1
CORRUPTION=0
DEBUG=0
OPTIMIZATION=1

RECEIVE_PROCESSING_TIME=(0 0 0 0 30 0 80 0 150 0 250 0 480 0 670 0 1200)
RECEIVE_PROCESSING_TIME_OPT=(0 0 0 0 20 0 50 0 70 0 120 0 180 0 300 0 400)

# Create directories
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

# Generate keys
go run main.go generate_keys $N

# Set processing time based on optimization flag
if [ $OPTIMIZATION -eq 0 ]; then
	PROCESSING_TIME=${RECEIVE_PROCESSING_TIME[$N]}
else
	PROCESSING_TIME=${RECEIVE_PROCESSING_TIME_OPT[$N]}
fi

# Run
for (( ID=0; ID<$N; ID++ )); do
	# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Optimization> <Initial value>
	go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $PROCESSING_TIME $OPTIMIZATION 0 &
done

