package modules

import (
	"bytes"
	"encoding/gob"
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/messenger"
	"self-stabilizing-binary-consensus/types"
	"strconv"
)

func TestSS(identifier int) {
	if _, in := messenger.SSBCChannel[identifier]; !in {
		messenger.SSBCChannel[identifier] = make(chan struct {
			SSBCMessage types.SSBCMessage
			From        int
		})
	}

	go receive2(identifier)

	r := 0
	for {
		r++
		send2("EST", types.NewSSBCMessage(identifier, true, r, 1, 1, 0, 0))
		logger.OutLogger.Println("SEND ", r)
	}
}

func send2(tag string, estMessage types.SSBCMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(estMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	messenger.Broadcast(types.NewMessage(w.Bytes(), tag))
}

func receive2(identifier int) {
	for message := range messenger.SSBCChannel[identifier] {
		j := message.From // From
		//aJ := message.SSBCMessage.Flag     // Flag
		rJ := message.SSBCMessage.Round // Round
		//est_0 := message.SSBCMessage.Est_0 // est[0]
		//est_1 := message.SSBCMessage.Est_1 // est[1]
		//aux_0 := message.SSBCMessage.Aux_0 // aux[0]
		//aux_1 := message.SSBCMessage.Aux_1 // aux[1]
		logger.OutLogger.Println("RECEIVE r=" + strconv.Itoa(rJ) + " j=" + strconv.Itoa(j))
	}
}
