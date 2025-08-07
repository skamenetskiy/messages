module github.com/skamenetskiy/messages

go 1.24.5

replace github.com/skamenetskiy/messages/api => ./api

require (
	github.com/goccy/go-yaml v1.18.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/lib/pq v1.10.9
	github.com/skamenetskiy/messages/api v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.16.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
