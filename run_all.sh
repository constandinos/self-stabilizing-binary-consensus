#!/bin/bash

M=6
CLIENTS=1
REMOTE=0
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

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

for (( BYZANTINE_SCENARIO=0; BYZANTINE_SCENARIO<=4; BYZANTINE_SCENARIO++ )); do
	for (( N=4; N<=12; N++ )); do
		for (( i=0; i<10; i++ )); do
			echo "BYZANTINE_SCENARIO:"$BYZANTINE_SCENARIO "N:"$N "i:"$i
			go run main.go generate_keys $N
			sleep 4
			for (( ID=0; ID<$N; ID++ )); do
				go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
			done
			sleep $(( $N - 2 ))
			#sleep 2
			./kill.sh
			sleep 2
			grep "stats" logs/out/*.log
			grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt
			rm logs/out/*.log
		done
		sort -n -k2 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
		rm logs/out/temp.txt
	done
done

