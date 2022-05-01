#!/bin/bash

N=16
M=6
CLIENTS=1
REMOTE=0
BYZANTINE_SCENARIO=0
CORRUPTION=0
SELF_STABILIZING=1
DEBUG=0
OPTIMIZATION=1

#RECEIVE_PROCESSING_TIME=(0 0 0 0 30 0 80 0 150 0 250 0 480 0 670 0 1200)
#RECEIVE_PROCESSING_TIME_OPT=(0 0 0 0 20 0 50 0 70 0 120 0 180 0 300 0 400)

PROCESSING_TIME=(1000 1200 1400)

# Create directories
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

# Generate keys
go run main.go generate_keys $N

for t in ${PROCESSING_TIME[@]}; do
	# Run
	for (( ID=0; ID<$N; ID++ )); do
		# <ID> <N> <M> <Clients> <Remote> <Byzantine scenario> <Self-Stabilizing> <Corruption> <Debug> <Receive processing time> <Optimization> <Initial value>
		go run main.go $ID $N $M $CLIENTS $REMOTE $BYZANTINE_SCENARIO $SELF_STABILIZING $CORRUPTION $DEBUG $t $OPTIMIZATION 0 &
	done
	sleep $N
	sh ./kill.sh
	echo -n $t" " >> temp.txt
	grep "stats" logs/out/*.log | awk '{print $5, $6}' | awk '($1=="false"){time+=$2;count+=1} END{print time/count}' >> temp.txt
	rm logs/out/*.log
done

cat temp.txt | sort -n -k2
rm temp.txt

