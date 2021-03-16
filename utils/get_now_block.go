package utils

import (
	"context"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"time"
)

func GetNowBlock(uri string) (int64, error) {
	//g := NewGrpcClient(uri)
	g := client.NewGrpcClient(uri)
	err := g.Start()
	if err != nil {
		return 0, err
	}
	GrpcTimeout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), GrpcTimeout)
	defer cancel()

	rt, err := g.Client.GetNowBlock(ctx, new(api.EmptyMessage))
	if err != nil {
		return 0, err
	}
	return rt.BlockHeader.RawData.Number, nil
}
