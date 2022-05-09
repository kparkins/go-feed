package message_test

import (
	"go-message/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedPubSubSimple(t *testing.T) {
	f := message.NewPublisher[int]()
	sub := f.Subscribe()
	f.Publish(1)
	<-sub.Updated()
	assert.Equal(t, 1, sub.Value())
}
