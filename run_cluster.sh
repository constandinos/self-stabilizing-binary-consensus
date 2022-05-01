#!/bin/bash

N=12
M=6
CLIENTS=1
REMOTE=1
BYZANTINE_SCENARIO=0
SELF_STABILIZING=1
CORRUPTION=0
DEBUG=0
OPTIMIZATION=1

RECEIVE_PROCESSING_TIME=(0 0 0 0 30 0 80 0 150 0 250 0 480 0 670 0 1200)
RECEIVE_PROCESSING_TIME_OPT=(0 0 0 0 20 0 50 0 70 0 120 0 180 0 300 0 400)

CLUSTER_SIZE=5
MACHINE_ID=$1

# Check if machine_id is an argument
if [ "$MACHINE_ID" == "" ]; then
	echo "./run_cluster.sh <machine_id>"
	exit
fi

if [ $MACHINE_ID == 0 ]; then
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
	sleep 4
fi

# Set processing time based on optimization flag
if [ $OPTIMIZATION -eq 0 ]; then
	PROCESSING_TIME=${RECEIVE_PROCESSING_TIME[$N]}
else
	PROCESSING_TIME=${RECEIVE_PROCESSING_TIME_OPT[$N]}
fi

# Synchronization
echo $MACHINE_ID >> sync.txt
counter=$(cat sync.txt | wc -l)
while [ $counter -lt $CLUSTER_SIZE ]; do
	counter=$(cat sync.txt | wc -l)
done

# Run
ID=$MACHINE_ID
while [ $ID -lt $N ]; do
	# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Optimization> <Initial value>
	go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $PROCESSING_TIME $OPTIMIZATION 0 &
	ID=$(( $ID + $CLUSTER_SIZE ))
done

if [ $MACHINE_ID == 0 ]; then
	rm sync.txt
fi

