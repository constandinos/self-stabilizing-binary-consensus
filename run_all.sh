#!/bin/bash

M=5
CLIENTS=1
REMOTE=0

SELF_STABILIZATION=1
CORRUPTION_SCENARIO=1

ITERATIONS=10
MAX_SIZE=12

BYZ_STR=("1-Normal" "2-Idle" "3-Inverse" "4-HH" "5-Random")

mkdir results

for (( BYZANTINE_SCENARIO=0; BYZANTINE_SCENARIO<=4; BYZANTINE_SCENARIO++ ))
do
	for (( N=4; N<=$MAX_SIZE; N++ ))
	do
		for (( i=0; i<$ITERATIONS; i++ ))
		do
			go install self-stabilizing-binary-consensus
			self-stabilizing-binary-consensus generate_keys $N
			for (( ID=0; ID<$N; ID++ ))
			do
				self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $SELF_STABILIZATION 1 &
			done
			sleep $(( $N*3 ))
			sh ./kill.sh
		done
	done
	
	grep "statistics" logs/out/*.log | awk '{print $5, $6, $7}' | sort -n | awk '($2=="false"){sum[$1]+=$3; count[$1]+=1}END{for (i in sum){print i, sum[i]/count[i]}}' | sort -n | sort -n > results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
	mkdir logs/out/"${BYZ_STR[$BYZANTINE_SCENARIO]}"
	mv logs/out/*.log logs/out/"${BYZ_STR[$BYZANTINE_SCENARIO]}"
done

