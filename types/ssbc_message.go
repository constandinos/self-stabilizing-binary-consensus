package types

import (
	"bytes"
	"encoding/gob"
	"self-stabilizing-binary-consensus/logger"
)

// SSBCMessage - Self-Stabilizing Binary Consensus EST message struct
type SSBCMessage struct {
	Identifier int
	Flag       bool // aJ
	Round      int  // rJ
	Est_0      int  // vJ[0]
	Est_1      int  // vJ[1]
	Aux_0      int  // uJ[0]
	Aux_1      int  // uJ[1]
}

// NewEstMessage - Creates a new est message
func NewSSBCMessage(identifier int, flag bool, round int, est_0 int, est_1 int, aux_0 int, aux_1 int) SSBCMessage {
	return SSBCMessage{
		Identifier: identifier,
		Flag:       flag,
		Round:      round,
		Est_0:      est_0,
		Est_1:      est_1,
		Aux_0:      aux_0,
		Aux_1:      aux_1}
}

// GobEncode - Binary Consensus message encoder
func (estm SSBCMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(estm.Identifier)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Flag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Round)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Est_0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Est_1)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Aux_0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(estm.Aux_1)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Binary Consensus message decoder
func (estm *SSBCMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&estm.Identifier)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Flag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Round)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Est_0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Est_1)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Aux_0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&estm.Aux_1)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
