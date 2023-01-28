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

func TestFinishedFeedIsFinished(t *testing.T) {
	pub := message.NewPublisher[int]()
	sub := pub.Subscribe()
	pub.Finish()
	<-sub.Updated()
	assert.Equal(t, true, sub.Finished())
}

func TestFeedFinishNoValue(t *testing.T) {
	pub := message.NewPublisher[int]()
	sub := pub.Subscribe()
	pub.Finish()
	<-sub.Updated()
	assert.Equal(t, true, sub.Finished())
	assert.Equal(t, 0, sub.Value())
	assert.Equal(t, false, sub.Next())
}

func TestFeedUnsubscribe(t *testing.T) {
	pub := message.NewPublisher[int]()
	defer pub.Finish()
	sub := pub.Subscribe()
	sub.Unsubscribe()
	assert.Equal(t, true, sub.Finished())
	assert.Equal(t, 0, sub.Value())
	assert.Equal(t, false, sub.Next())
	<-sub.Updated()
}

func TestPublishToFinishedFeed(t *testing.T) {
	pub := message.NewPublisher[int]()
	sub := pub.Subscribe()
	pub.Finish()
	assert.True(t, sub.Finished())
	assert.NotNil(t, pub.Publish(1))
}
