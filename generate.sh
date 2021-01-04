#!/bin/bash

protoc unary/greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc unary/calculator/proto/calculator.proto --go_out=plugins:.