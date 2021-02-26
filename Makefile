genG:
	  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative calc/calcpb/*.proto
clean:
	rm greet/calcpb/*.go
runs:
	go run calc/calc_server/server.go
runc:
	go run calc/calc_client/client.go