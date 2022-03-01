package messenger

import (
	"bytes"
	"encoding/gob"
	"self-stabilizing-binary-consensus/config"
	"self-stabilizing-binary-consensus/logger"
	"self-stabilizing-binary-consensus/threshenc"
	"self-stabilizing-binary-consensus/types"
	"self-stabilizing-binary-consensus/variables"
	"time"

	"github.com/pebbe/zmq4"
)

// Sockets
var (
	// Context to initialize sockets
	Context *zmq4.Context

	// SendSockets - Send messages to other servers
	SendSockets map[int]*zmq4.Socket

	// ReceiveSockets - Receive messages from other servers
	ReceiveSockets map[int]*zmq4.Socket

	// ServerSockets - Get the client requests
	ServerSockets map[int]*zmq4.Socket

	// ResponseSockets - Send responses to clients
	ResponseSockets map[int]*zmq4.Socket
)

// Channels for messages
var (
	// MessageChannel - Channel to put the messages that need to be transmitted in
	MessageChannel = make(map[int]chan types.Message)

	// BvbChannel - Channel to put the BVB messages in
	BvbChannel = make(map[int]chan struct {
		BcMessage types.BcMessage
		From      int
	})

	// BcChannel - Channel to put the BC messages in
	BcChannel = make(map[int]chan struct {
		BcMessage types.BcMessage
		From      int
	})

	SSBCChannel = make(map[int]chan struct {
		SSBCMessage types.SSBCMessage
		From        int
	})

	// RequestChannel - Channel to put the client requests in
	RequestChannel = make(chan []byte, 100)
)

// InitializeMessenger - Initializes the 0MQ sockets (between Servers and Clients)
func InitializeMessenger() {
	Context, err := zmq4.NewContext()
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	// Initialization of a socket pair to communicate with each one of the other servers
	ReceiveSockets = make(map[int]*zmq4.Socket)
	SendSockets = make(map[int]*zmq4.Socket)
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}

		// ReceiveSockets initialization to get information from other servers
		ReceiveSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var receiveAddr string
		if !variables.Remote {
			receiveAddr = config.GetRepAddressLocal(i)
		} else {
			receiveAddr = config.GetRepAddress(i)
		}
		err = ReceiveSockets[i].Bind(receiveAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Receive from Server", i, "on", receiveAddr)

		// SendSockets initialization to send information to other servers
		SendSockets[i], err = Context.NewSocket(zmq4.REQ)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var sendAddr string
		if !variables.Remote {
			sendAddr = config.GetReqAddressLocal(i)
		} else {
			sendAddr = config.GetReqAddress(i)
		}
		err = SendSockets[i].Connect(sendAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Send to Server", i, "on", sendAddr)

		// Init message channel
		MessageChannel[i] = make(chan types.Message)
	}

	logger.OutLogger.Println("-----------------------------------------")

	// Initialization of a socket pair to communicate with each one of the clients
	ServerSockets = make(map[int]*zmq4.Socket, variables.Clients)
	ResponseSockets = make(map[int]*zmq4.Socket, variables.Clients)
	for i := 0; i < variables.Clients; i++ {

		// ServerSockets initialization to get clients requests
		ServerSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var serverAddr string
		if !variables.Remote {
			serverAddr = config.GetServerAddressLocal(i)
		} else {
			serverAddr = config.GetServerAddress(i)
		}
		err = ServerSockets[i].Bind(serverAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Requests from Client", i, "on", serverAddr)

		// ResponseSockets initialization to publish the response back to the clients
		ResponseSockets[i], err = Context.NewSocket(zmq4.PUB)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var responseAddr string
		if !variables.Remote {
			responseAddr = config.GetResponseAddressLocal(i)
		} else {
			responseAddr = config.GetResponseAddress(i)
		}
		err = ResponseSockets[i].Bind(responseAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Response to Client", i, "on", responseAddr)
	}

	logger.OutLogger.Print("-----------------------------------------\n\n")
}

// Function to modify BC messages if byzantine
func modifyMessageBC(message types.Message, receiver int) types.Message {
	msg := new(types.BcMessage)
	buf := bytes.NewBuffer(message.Payload)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	switch msg.Tag % 3 {
	case 0:
		msg.Value = uint(receiver % 2)
	case 1:
		msg.Value = uint(0)
	case 2:
		msg.Value = uint(1)
	}

	logger.ErrLogger.Print(config.Scenario, ": (", message.Type, ") ", receiver, " --> [",
		msg.Tag, ",", msg.Value, "]\n")

	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err = encoder.Encode(msg)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	return types.NewMessage(w.Bytes(), message.Type)
}

// Function to modify messages and send only to the half the right message if byzantine
func modifyMessageHH(message types.Message, receiver int) types.Message {
	var newPayload []byte

	if (message.Type == "BVB") || (message.Type == "BC") {
		msg := new(types.BcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err := dec.Decode(&msg)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		msg.Value = uint(receiver % 2)

		logger.ErrLogger.Print(config.Scenario, ": (", message.Type, ") ", receiver, " --> [",
			msg.Tag, ",", msg.Value, "]\n")

		w := new(bytes.Buffer)
		encoder := gob.NewEncoder(w)
		err = encoder.Encode(msg)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		newPayload = w.Bytes()

	}

	return types.NewMessage(newPayload, message.Type)
}

// Broadcast - Broadcasts a message to all other servers
func Broadcast(message types.Message) {
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}

		// Modify message before sending it, in case of a special scenario
		if (config.Scenario == "BC_ATTACK") && (variables.Byzantine) &&
			((message.Type == "BVB") || (message.Type == "BC")) {
			message = modifyMessageBC(message, i)
		}

		if (config.Scenario == "HALF_&_HALF") && (variables.Byzantine) {
			message = modifyMessageHH(message, i)
		}

		timeout := time.NewTicker(5 * time.Second)
		select {
		case MessageChannel[i] <- message:
		case <-timeout.C:
		}
	}
}

// TransmitMessages - Transmits the messages to the other servers [started from main]
func TransmitMessages() {
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		go func(i int) { // Initializes them with a goroutine and waits forever
			for message := range MessageChannel[i] {
				w := new(bytes.Buffer)
				encoder := gob.NewEncoder(w)
				err := encoder.Encode(message)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				_, err = SendSockets[i].SendBytes(w.Bytes(), 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				_, err = SendSockets[i].Recv(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
				logger.OutLogger.Println("SENT", message.Type, "to", i)

				variables.MsgMutex.Lock()
				variables.MsgComplexity++
				variables.MsgSize += int64(len(w.Bytes()))
				variables.MsgMutex.Unlock()
			}
		}(i)
	}
}

// Subscribe - Handles the inputs from both clients and other servers [started from main]
func Subscribe() {
	// Gets messages from other servers and handles them
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		go func(i int) { // Initializes them with a goroutine and waits forever
			for {
				message, err := ReceiveSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				go HandleMessage(message)

				_, err = ReceiveSockets[i].Send("", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}

	// Gets requests from clients and handles them
	for i := 0; i < variables.Clients; i++ {
		go func(i int) { // Initialize them with a goroutine and waits forever
			for {
				message, err := ServerSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				go handleRequest(message, i)

				_, err = ServerSockets[i].Send("", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}
}

// Put client's message in RequestChannel to be handled
func handleRequest(message []byte, from int) {
	logger.OutLogger.Println("RECEIVED REQ from", from)
	RequestChannel <- message
}

// HandleMessage - Handles the messages from the other servers
func HandleMessage(msg []byte) {
	message := new(types.Message)
	buffer := bytes.NewBuffer([]byte(msg))
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	if !(threshenc.VerifyMessage(message.Payload, message.Signature, message.From)) {
		logger.OutLogger.Println("INVALID", message.Type, "from", message.From)
		return
	}

	logger.OutLogger.Println("RECEIVED", message.Type, "from", message.From)

	switch message.Type {
	case "BVB":
		bcMessage := new(types.BcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&bcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		tag := bcMessage.Tag
		if _, in := BvbChannel[tag]; !in {
			BvbChannel[tag] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}

		BvbChannel[tag] <- struct {
			BcMessage types.BcMessage
			From      int
		}{BcMessage: *bcMessage, From: message.From}

	case "BC":
		bcMessage := new(types.BcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&bcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		tag := bcMessage.Tag
		if _, in := BcChannel[tag]; !in {
			BcChannel[tag] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}

		BcChannel[tag] <- struct {
			BcMessage types.BcMessage
			From      int
		}{BcMessage: *bcMessage, From: message.From}

	case "EST":
		ssbcMessage := new(types.SSBCMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&ssbcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		round := ssbcMessage.Round
		if _, in := SSBCChannel[round]; !in {
			SSBCChannel[round] = make(chan struct {
				SSBCMessage types.SSBCMessage
				From        int
			})
		}

		SSBCChannel[round] <- struct {
			SSBCMessage types.SSBCMessage
			From        int
		}{SSBCMessage: *ssbcMessage, From: message.From}
	}
}
