package grpc

import (
	v1 "github.com/andreyxaxa/PubSub_gRPC_Service/internal/controller/grpc/v1"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/logger"
	"github.com/andreyxaxa/PubSub_gRPC_Service/pkg/subpub"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, sp subpub.SubPub, l logger.Interface) {
	{
		v1.NewPubSubRouter(app, sp, l)
	}

	reflection.Register(app)
}
