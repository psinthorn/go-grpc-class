#!/bin/bash

protoc services/greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc services/calculator/proto/calculator.proto --go_out=plugins=grpc:.

# Start calculator server
go run services/calculator/server/server.go


# Start greeting server
go run services/greet/server/server.go