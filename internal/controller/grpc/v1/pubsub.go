package v1

import (
	"context"

	v1 "github.com/andreyxaxa/PubSub_gRPC_Service/docs/proto/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *V1) Subscribe(req *v1.SubscribeRequest, stream grpc.ServerStreamingServer[v1.Event]) error {
	sub, err := r.sp.Subscribe(req.GetKey(), func(msg interface{}) {
		data, ok := msg.(string)
		if !ok {
			r.l.Error(nil, "invalid message type")
			return
		}
		stream.Send(&v1.Event{Data: data})
	})
	if err != nil {
		r.l.Error(err, "grpc - v1 - Subscribe")

		return status.Errorf(codes.Internal, "Subscribe failed: %v", err)
	}
	defer sub.Unsubscribe()

	<-stream.Context().Done()
	return nil
}

func (r *V1) Publish(ctx context.Context, req *v1.PublishRequest) (*emptypb.Empty, error) {
	err := r.sp.Publish(req.GetKey(), req.GetData())
	if err != nil {
		r.l.Error(err, "grpc - v1 - Publish")

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Publish failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
