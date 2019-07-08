// Copyright 2016 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"cloud.google.com/go/storage"
	pb "cloud.google.com/go/storage/internal/benchwrapper/proto"
	"google.golang.org/grpc"
)

const port = ":50051"

type server struct {
	c *storage.Client
}

func main() {
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterStorageServer(s, &server{
		c: c,
	})
	fmt.Printf("Running on %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (s *server) Read(ctx context.Context, in *pb.ObjectRead) (*pb.EmptyResponse, error) {
	b := s.c.Bucket(in.GetBucketName())
	o := b.Object(in.GetObjectName())
	r, err := o.NewReader(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	var ba []byte
	fmt.Println("Is reading")
	if _, err = r.Read(ba); err != nil {
		log.Fatal(err)
	}
	return &pb.EmptyResponse{}, nil
}

func (s *server) Write(ctx context.Context, in *pb.ObjectWrite) (*pb.EmptyResponse, error) {
	return &pb.EmptyResponse{}, nil
}
