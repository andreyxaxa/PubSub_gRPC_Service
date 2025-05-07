package v1

import (
	v1 "github.com/andreyxaxa/PubSub_gRPC_Service/docs/proto/pubsub/v1"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/logger"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub"
)

type V1 struct {
	v1.UnimplementedPubSubServer

	// Dependency Injection
	sp subpub.SubPub
	l  logger.Interface
}
