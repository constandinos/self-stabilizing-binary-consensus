#!/bin/bash

N=$1
M=$2
CLIENTS=$3
REMOTE=$4
BYZANTINE_SCENARIO=$5
SELF_STABILIZING=$6
CORRUPTION=$7
DEBUG=$8
OPTIMIZATION=$9

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

# Run
for (( ID=0; ID<$N; ID++ )); do

	# go run main.go <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Optimization> <Initial value>
	
	# Initial value = 0
	go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $OPTIMIZATION 0 &
	
	# Initial value = 1
	# go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $OPTIMIZATION 1 &
	
	# Initial value = ID%2
	# go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $OPTIMIZATION $(($ID%2)) &
	
done

