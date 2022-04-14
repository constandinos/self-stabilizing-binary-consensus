#!/bin/bash

M=5
CLIENTS=1
REMOTE=0

SELF_STABILIZATION=1
CORRUPTION_SCENARIO=1

BYZ_STR=("0-Normal" "1-Idle" "2-Inverse" "3-HH" "4-Random")
#SLEEP_TIME=(0 0 0 0 1 1 2 2 3 4 4 12 12)

mkdir results

for (( BYZANTINE_SCENARIO=0; BYZANTINE_SCENARIO<=4; BYZANTINE_SCENARIO++ ))
do
	for (( N=4; N<=12; N++ ))
	do
		for (( i=0; i<10; i++ ))
		do
			go install self-stabilizing-binary-consensus
			self-stabilizing-binary-consensus generate_keys $N
			for (( ID=0; ID<$N; ID++ ))
			do
				self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $CORRUPTION_SCENARIO $SELF_STABILIZATION 1 &
			done
			sleep $(( $N-2 ))
			# sleep ${SLEEP_TIME[$N]}
			sh ./kill.sh
			grep "stats" logs/out/*.log | awk '{print $5, $6, $7}' | awk '($1=="false"){time+=$2;msg+=$3;count+=1} END{print time/count,msg/count}' >> logs/out/temp.txt
			rm logs/out/*.log
		done
		sort -n -k2 logs/out/temp.txt | awk 'BEGIN{i=0} {t[i]=$1; m[i]=$2; i++} END{for(i=1; i<NR-1; i++){time+=t[i]; msg+=m[i]; count+=1} print time/count,msg/count}' >> results/"${BYZ_STR[$BYZANTINE_SCENARIO]}".txt
		rm logs/out/temp.txt
	done
	
	#mkdir logs/out/"${BYZ_STR[$BYZANTINE_SCENARIO]}"
	#mv logs/out/*.log logs/out/"${BYZ_STR[$BYZANTINE_SCENARIO]}"
done

