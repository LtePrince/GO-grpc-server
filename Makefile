# localhost:16686
jaeger-up:
	sudo docker-compose -f ./docker/docker-compose.yml up -d

jaeger-down:
	sudo docker-compose -f ./docker/docker-compose.yml down

# protofile generation
# export PATH="/home/adolph/Envirenment/protobuffers/protoc-31.0-linux-x86_64/bin:$PATH"
# protoc --go_out=. --go-grpc_out=. pkg/api/user.proto

run:
	go run ./cmd/server/main.go

client:
	go run ./client_test/client.go

cleandb:
	psql -U adolph -d grpctest -c "DROP TABLE users;"
	psql -U adolph -d grpctest -f ./pkg/storage/migrations/users.sql

test:
	go test ./client_test -v

pprof:
	mkdir -p ./report
	rm -rf ./report/*
	curl http://localhost:6060/debug/pprof/profile?seconds=30 -o ./report/profile.out
	go tool pprof -svg ./report/profile.out > ./report/profile.svg

.PHONY: jaeger-up jaeger-down run client cleandb test