#!/bin/bash

M=6
CLIENTS=1
REMOTE=0
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

BYZ_STR=("0-Normal" "1-Idle" "2-Inverse" "3-HH" "4-Random")

CLUSTER_SIZE=5
MACHINE_ID=$1

# Check if machine_id is an argument
if [ "$MACHINE_ID" == "" ]; then
	echo "./run_cluster.sh <machine_id>"
	exit
fi

# Create directories
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

	if [ ! -d "./results" ]; then
		mkdir results
	else
		rm -rf ./results
		mkdir results
	fi
fi

# Byzantine Scenario
for (( BYZANTINE_SCENARIO=0; BYZANTINE_SCENARIO<=4; BYZANTINE_SCENARIO++ )); do
	# Network size
	for (( N=4; N<=12; N++ )); do
		# Experements
		for (( i=0; i<10; i++ )); do
			# Create keys
			if [ $MACHINE_ID == 0 ]; then
				echo "BYZANTINE_SCENARIO:"$BYZANTINE_SCENARIO "N:"$N "i:"$i
				go run main.go generate_keys $N
				sleep 4
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
				# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Initial value>
				go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
				ID=$(( $ID + $CLUSTER_SIZE ))
			done
			if [ $MACHINE_ID == 0 ]; then
				rm sync.txt
			fi
			# Wait
			sleep $(( $N - 2 ))
			# Stop processes
			./kill.sh
			sleep 2
			# Get local results
			grep "stats" logs/out/*.log
			grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt
			rm logs/out/*.log
		done
		# Get global results
		sort -n -k2 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
		rm logs/out/temp.txt
	done
done

