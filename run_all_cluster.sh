#!/bin/bash

M=6
CLIENTS=1
REMOTE=0
SELF_STABILIZING=0
CORRUPTION=0
DEBUG=0
RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

CLUSTER_SIZE=5
CLUSTER_ID=$1

BYZ_STR=("0-Normal" "1-Idle" "2-Inverse" "3-HH" "4-Random")

if [ ! -d "./logs" ]; then
	mkdir -p logs/{error,out}
fi

if [ ! -d "./keys" ]; then
	mkdir keys
fi

if [ ! -d "./results" ]; then
	mkdir results
fi

# go install self-stabilizing-binary-consensus
# self-stabilizing-binary-consensus generate_keys $N

for (( BYZANTINE_SCENARIO=0; BYZANTINE_SCENARIO<=4; BYZANTINE_SCENARIO++ )); do
	for (( N=4; N<=12; N++ )); do
		for (( i=0; i<10; i++ )); do
			echo "BYZANTINE_SCENARIO:"$BYZANTINE_SCENARIO "N:"$N "i:"$i
			for (( ID=0; ID<$N; ID++ )); do
				if [ $((ID%$CLUSTER_SIZE)) == $CLUSTER_ID ]; then
					self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG ${RECEIVE_PROCESSING_TIME[$N]} 0 &
				fi
			done
			sleep 2
			sh ./kill.sh
			if [ $CLUSTER_ID == 0 ]; then
				grep "stats" logs/out/*.log
				grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt
				rm logs/out/*.log
			fi
		done
		if [ $CLUSTER_ID == 0 ]; then
			sort -n -k2 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
			rm logs/out/temp.txt
		fi
	done
done

