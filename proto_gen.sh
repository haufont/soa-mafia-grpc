protoc -I=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./internal/role/role.proto          \
    ./internal/player/player.proto      \
    ./internal/event/event.proto        \
    ./internal/state/state.proto        \
    ./api/api.proto
