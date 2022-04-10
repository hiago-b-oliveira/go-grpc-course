#!/bin/bash
protoc blog/blogpb/blog.proto --go_out=. --go-grpc_out=.