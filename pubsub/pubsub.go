package pubsub

import (
	"context"
	"encoding/json"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"

	"github.com/likecoin/likecoin-chain-tx-indexer/logger"
)

const (
	CmdProjectID = "pubsub-project-id"
	CmdTopic     = "pubsub-topic"
	CmdChainID   = "pubsub-chain-id"
)

var (
	client        *pubsub.Client
	topic         *pubsub.Topic
	pubsubChainID string
	lock          = &sync.Mutex{}
)

func InitPubsubFromCmd(cmd *cobra.Command) error {
	projectID, err := cmd.Flags().GetString(CmdProjectID)
	if err != nil {
		return err
	}
	topicID, err := cmd.Flags().GetString(CmdTopic)
	if err != nil {
		return err
	}
	if topicID == "" {
		logger.L.Infow("Pubsub topic is empty, pubsub disabled")
		return nil
	}
	chainID, err := cmd.Flags().GetString(CmdChainID)
	if err != nil {
		return err
	}
	logger.L.Infow("Pubsub enabled", "project_id", projectID, "topic", topicID, "chain_id", chainID)
	return InitPubsub(projectID, topicID, chainID)
}

func InitPubsub(projectID, topicID, chainID string) error {
	lock.Lock()
	defer lock.Unlock()
	if client == nil {
		c, err := pubsub.NewClient(context.Background(), projectID)
		if err != nil {
			return err
		}
		client = c
		topic = client.Topic(topicID)
		pubsubChainID = chainID
	}
	return nil
}

type PubsubPayload struct {
	Event   string      `json:"event"`
	Source  string      `json:"source"`
	ChainID string      `json:"chain_id"`
	Payload interface{} `json:"payload"`
}

func Publish(action string, payload interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	if topic == nil {
		return nil
	}
	encodedBz, err := json.Marshal(PubsubPayload{
		Event:   action,
		Source:  "indexer",
		ChainID: pubsubChainID,
		Payload: payload,
	})
	if err != nil {
		logger.L.Errorw("Fail to marshal pubsub data", "action", action, "payload", payload, "error", err)
		return err
	}
	res := topic.Publish(context.Background(), &pubsub.Message{
		Data: encodedBz,
	})
	logger.L.Debugw("Publishing message", "data", string(encodedBz))
	go func() {
		id, err := res.Get(context.Background())
		logger.L.Debugw("Topic publish done", "id", id, "error", err)
	}()
	return nil
}

func ConfigCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String(CmdChainID, "", "Chain ID marked in pubsub payloads")
	cmd.PersistentFlags().String(CmdProjectID, "", "Pubsub project ID")
	cmd.PersistentFlags().String(CmdTopic, "", "Pubsub topic (empty means disable pubsub)")
}
