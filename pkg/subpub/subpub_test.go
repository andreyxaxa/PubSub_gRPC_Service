package subpub_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub"
	subpuberr "github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub/errors"
	"github.com/stretchr/testify/assert"
)

func TestPublishSubscribe(t *testing.T) {
	sp := subpub.NewSubPub()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	_, err := sp.Subscribe("WEATHER", func(msg interface{}) {
		if msg != "cloudy" {
			t.Errorf("unexpected message: %v", msg)
		}
		wg.Done()
	})

	assert.NoError(t, err)

	err = sp.Publish("WEATHER", "cloudy")
	assert.NoError(t, err)

	waitWithTimeout(wg, t)
}

func TestUnsubscribe(t *testing.T) {
	sp := subpub.NewSubPub()

	received := false

	sub, err := sp.Subscribe("NEWS", func(_ interface{}) {
		received = true
	})
	assert.NoError(t, err)

	sub.Unsubscribe()

	err = sp.Publish("NEWS", "breaking news")
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	assert.NotEqual(t, received, true)
}

func TestClose(t *testing.T) {
	sp := subpub.NewSubPub()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	_, err := sp.Subscribe("NEWS", func(_ interface{}) {
		wg.Done()
	})
	assert.NoError(t, err)

	sp.Publish("NEWS", "breaking news")
	waitWithTimeout(wg, t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = sp.Close(ctx)
	assert.NoError(t, err)

	err = sp.Publish("WEATHER", "sunny")
	assert.EqualError(t, err, subpuberr.ErrSubPubClosed.Error())
}

func waitWithTimeout(wg *sync.WaitGroup, t *testing.T) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for message")
	}
}
