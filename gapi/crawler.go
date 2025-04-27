package gapi

import (
	"context"
	"github.com/pule1234/VideoForge/internal/crawler"
	"github.com/pule1234/VideoForge/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) Crawler(ctx context.Context, req *pb.CrawlerRequest) (*emptypb.Empty, error) {
	dc, err := crawler.NewDyCrawler(server.config.DouYingQueueName, server.store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create DyCrawler: %v", err)
	}

	err = dc.Start(req.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not start DyCrawler: %v", err)
	}
	return nil, nil
}
