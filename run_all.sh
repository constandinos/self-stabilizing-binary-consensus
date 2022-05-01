#!/bin/bash

M=6
CLIENTS=1
REMOTE=0
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
#RECEIVE_PROCESSING_TIME=(0 0 0 0 30 0 80 0 150 0 250 0 480 0 670 0 1200) #SEND ALL
RECEIVE_PROCESSING_TIME=(0 0 0 0 20 0 50 0 70 0 120 0 180 0 300 0 400) #SEND%2

BYZ_STR=("0-Normal" "1-Idle" "2-Inverse" "3-HH" "4-Random")

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

for BYZANTINE_SCENARIO in {0..0}; do
	for N in {12..16..2}; do
		for REPEAT in {1..10}; do
			echo "BYZANTINE_SCENARIO:"$BYZANTINE_SCENARIO "N:"$N "REPEAT:"$REPEAT
			
			# Create keys
			go run main.go generate_keys $N
			sleep 4
			
			# Run
			for (( ID=0; ID<$N; ID++ )); do
				echo go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0
				go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
			done
			
			# Wait
			echo "sleep" $N
			sleep $N
			
			# Stop processes
			./kill.sh
			echo "kill processes"
			sleep 2
			
			# Get local results
			grep "stats" logs/out/*.log
			grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}'
			grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt
			rm logs/out/*.log			
		done
		# Get global results
		sort -n -k1 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
		rm logs/out/temp.txt
	done
done
