package v1

import (
	v1 "github.com/andreyxaxa/PubSub_gRPC_Service/docs/proto/pubsub/v1"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/logger"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub"
	pbgrpc "google.golang.org/grpc"
)

func NewPubSubRouter(app *pbgrpc.Server, sp subpub.SubPub, l logger.Interface) {
	r := &V1{
		sp: sp,
		l:  l,
	}

	{
		v1.RegisterPubSubServer(app, r)
	}
}
