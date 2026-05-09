module storage-service

go 1.26.1

require (
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.81.0
	storagepb v0.0.0-00010101000000-000000000000
)

replace storagepb => ../proto/storagepb

require (
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260504160031-60b97b32f348 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
