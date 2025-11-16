module github.com/dizzrt/ellie/cmd/protoc-gen-ellie-go

go 1.25

require google.golang.org/protobuf v1.36.10

replace google.golang.org/protobuf => github.com/dizzrt/protobuf-go v1.36.10-ellie.1
