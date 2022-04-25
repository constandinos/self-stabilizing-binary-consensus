#!/bin/bash

N=5
M=6
CLIENTS=1
REMOTE=1
BYZANTINE_SCENARIO=0
SELF_STABILIZING=1
CORRUPTION=0
DEBUG=1
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

CLUSTER_SIZE=5
MACHINE_ID=$1

if [ "$MACHINE_ID" == "" ]; then
	echo "./run_cluster.sh <machine_id>"
	exit
fi

if [ $MACHINE_ID == 0 ]; then
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
fi

echo $MACHINE_ID >> sync.txt
counter=$(cat sync.txt | wc -l)
while [ $counter -lt $CLUSTER_SIZE ]; do
	counter=$(cat sync.txt | wc -l)
done

ID=$MACHINE_ID
while [ $ID -lt $N ]; do
	# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Initial value>
	go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
	ID=$(( $ID + $CLUSTER_SIZE ))
done

if [ $MACHINE_ID == 0 ]; then
	rm sync.txt
fi

