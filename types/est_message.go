package types

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"self-stabilizing-binary-consensus/logger"
)

// ESTMessage - Binary Consensus EST message struct
type ESTMessage struct {
	AJ bool       // flag
	RJ int        // round
	VJ *list.List // est[i][j]
	UJ *list.List // aux[i][j]
}

// NewEstMessage - Creates a new est message
func NewESTMessage(aJ bool, rJ int, vJ *list.List, uJ *list.List) ESTMessage {
	return ESTMessage{AJ: aJ, RJ: rJ, VJ: vJ, UJ: uJ}
}

// GobEncode - Binary Consensus message encoder
func (estm ESTMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(estm.AJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.RJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.VJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.UJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Binary Consensus message decoder
func (estm *ESTMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&estm.AJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.RJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.VJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.UJ)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
