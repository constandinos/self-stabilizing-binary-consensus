#!/bin/bash

N=16
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION=0
SELF_STABILIZING=1
DEBUG=0
# RECEIVE_PROCESSING_TIME=(0 0 0 0 30 45 90 100 150 210 270 310 370)

TIME=(1000 1200 1400)

if [ ! -d "./logs" ]; then
	mkdir -p logs/{error,out}
fi

if [ ! -d "./keys" ]; then
	mkdir keys
fi

go run main.go generate_keys $N

for t in ${TIME[@]}; do
	for (( ID=0; ID<$N; ID++ )); do
		go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $t 0 &
	done
	sleep $N
	sh ./kill.sh
	echo -n $t" " >> temp.txt
	grep "stats" logs/out/*.log | awk '{print $5, $6}' | awk '($1=="false"){time+=$2;count+=1} END{print time/count}' >> temp.txt
	rm logs/out/*.log
done

cat temp.txt | sort -n -k2
rm temp.txt

