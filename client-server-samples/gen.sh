#!/bin/bash
protoc   --go-grpc_out=. --go_out=.  greet/greetpb/greet.proto
protoc   --go-grpc_out=. --go_out=.  calculator/calculatorpb/calculator.proto