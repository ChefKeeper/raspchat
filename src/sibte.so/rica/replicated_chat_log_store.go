package rica

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"sibte.so/rascluster"
	"time"
    "fmt"
)

type replicatedMessage struct {
	Id      uint64        `json:"id"`
	Group   string        `json:"group"`
	Message interface{} `json:"message"`
}

type replicatedChatLogStore struct {
	stateMachine rascluster.ClusterStateMachine
	logStore     ChatLogStore
}

func NewReplicatedChatLogStore(logStore ChatLogStore, stateMachine rascluster.ClusterStateMachine) ChatLogStore {
	store := &replicatedChatLogStore{
		logStore:     logStore,
		stateMachine: stateMachine,
	}

	stateMachine.OnUpdate(store)
	return store
}

func (s *replicatedChatLogStore) Save(group string, id uint64, msg IEventMessage) error {
	replicatedMsg := replicatedMessage{
		Id:      id,
		Group:   group,
		Message: msg,
	}

	if err := s.stateMachine.ApplyMessage(s.serialize(replicatedMsg), 1*time.Second); err != nil {
		return err
	}

	return nil
}

func (c *replicatedChatLogStore) GetMessagesFor(group string, start_id string, offset uint, limit uint) ([]IEventMessage, error) {
	return c.logStore.GetMessagesFor(group, start_id, offset, limit)
}

func (c *replicatedChatLogStore) GetMessage(id uint64) (IEventMessage, error) {
	return c.logStore.GetMessage(id)
}

func (c *replicatedChatLogStore) Cleanup(group string) {
	c.logStore.Cleanup(group)
}

func (c *replicatedChatLogStore) OnStateUpdated(msg []byte) error {
	if msg == nil {
		return errors.New("Nil message")
	}

	if rMsg := c.deserialize(msg); rMsg != nil {
        evtMsg := ConvertToIEventMessage(rMsg.Message)
        if evtMsg == nil {
            return errors.New(fmt.Sprintf("Invalid message type %v", rMsg.Message))
        }

		log.Println("Applying message group:", rMsg.Group, "id:", rMsg.Id, "message:", rMsg.Message)
		if err := c.logStore.Save(rMsg.Group, rMsg.Id, evtMsg); err != nil {
			return err
		}

		log.Println("Applied message group:", rMsg.Group, "id:", rMsg.Id)
		return nil
	}

	return errors.New("Unable to decode message")
}

func (c *replicatedChatLogStore) serialize(v replicatedMessage) []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(v); err != nil {
		log.Println("Unable to encode message...", err)
		return nil
	}

	return buffer.Bytes()
}

func (c *replicatedChatLogStore) deserialize(b []byte) *replicatedMessage {
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)

	decM := &replicatedMessage{}
	if err := dec.Decode(decM); err != nil {
		log.Println("Error decoding message...", err)
		return nil
	}

	return decM
}
