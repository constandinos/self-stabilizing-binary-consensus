#!/bin/bash

N=4
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION=0
SELF_STABILIZING=1
DEBUG=0
# RECEIVE_PROCESSING_TIME=(0 0 0 0 30 40 100 150 200 250 350 450 500)

TIME=(20 25 30 35 40)

if [ ! -d "./logs" ]; then
	mkdir -p logs/{error,out}
fi

if [ ! -d "./keys" ]; then
	mkdir keys
fi

for t in ${TIME[@]}; do
	go install self-stabilizing-binary-consensus
	self-stabilizing-binary-consensus generate_keys $N	
	for (( ID=0; ID<$N; ID++ )); do
		self-stabilizing-binary-consensus $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $t $(( $ID%2 )) &
	done
	sleep 1
	sh ./kill.sh
	echo -n $t" " >> temp.txt
	grep "stats" logs/out/*.log | awk '{print $5, $6}' | awk '($1=="false"){time+=$2;count+=1} END{print time/count}' >> temp.txt
	rm logs/out/*.log
done

cat temp.txt | sort -n -k2
rm temp.txt

