#!/bin/bash

M=6
CLIENTS=1
REMOTE=1
SELF_STABILIZING=1
CORRUPTION=0
DEBUG=0
OPTIMIZATION=1

RECEIVE_PROCESSING_TIME=(0 0 0 0 30 0 80 0 150 0 250 0 480 0 670 0 1200)
RECEIVE_PROCESSING_TIME_OPT=(0 0 0 0 20 0 50 0 70 0 120 0 180 0 300 0 400)
SLEEP_TIME=(0 0 0 0 40 0 40 0 60 0 60 0 60 0 80 0 80)
BYZ_STR=("0-Normal" "1-Idle" "2-Inverse" "3-HH" "4-Random")

CLUSTER_SIZE=5
MACHINE_ID=$1
ROUND=0

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
for BYZANTINE_SCENARIO in {0..4}; do
	# Network size
	for N in {4..16..2}; do
		# Set processing time based on optimization flag
		if [ $OPTIMIZATION -eq 0 ]; then
			PROCESSING_TIME=${RECEIVE_PROCESSING_TIME[$N]}
		else
			PROCESSING_TIME=${RECEIVE_PROCESSING_TIME_OPT[$N]}
		fi
		
		# Experements
		for REPEAT in {1..10}; do
		
			# Generate keys
			if [ $MACHINE_ID == 0 ]; then
				echo "BYZANTINE_SCENARIO:"$BYZANTINE_SCENARIO "N:"$N "REPEAT:"$REPEAT
				go run main.go generate_keys $N
				sleep 4
			fi
		
			# Synchronization
			printf $ROUND"\n" >> sync_$MACHINE_ID.txt
			echo "Just printf "$ROUND" @ "sync_$MACHINE_ID.txt
			counter=$(cat sync_*.txt | grep -w $ROUND | wc -l)
			while [ $counter -lt $CLUSTER_SIZE ]; do
				counter=$(cat sync_*.txt | grep -w $ROUND | wc -l)
			done
			ROUND=$(( $ROUND + 1 ))
			
			# Run
			ID=$MACHINE_ID
			while [ $ID -lt $N ]; do
				# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Optimization> <Initial value>
				echo go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $PROCESSING_TIME $OPTIMIZATION 0
				go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $PROCESSING_TIME $OPTIMIZATION 0 &
				ID=$(( $ID + $CLUSTER_SIZE ))
			done
			
			# Wait
			echo "sleep 60"
			sleep 60
			
			# Stop processes
			./kill.sh
			echo "kill processes"
			sleep 2
			
			# Get local results
			if [ $MACHINE_ID == 0 ]; then
				grep "stats" logs/out/*.log 2> /dev/null
				grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt 2> /dev/null
				rm logs/out/*.log			
			fi
		done
		# Get global results
		if [ $MACHINE_ID == 0 ]; then
			sort -n -k1 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt 2> /dev/null
			rm logs/out/temp.txt
		fi
	done
done

rm sync_$MACHINE_ID.txt

