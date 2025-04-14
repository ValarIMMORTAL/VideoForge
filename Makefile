DB_URL=postgresql://sj_admin:123@localhost:5431/videoforge?sslmode=disable

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migratedownall:
	migrate -path db/migration -database "$(DB_URL)" -verbose drop
sqlc:
	sqlc generate

protoc:
	rm -f pb/*.go
	protoc \
	--proto_path=proto \
	--proto_path=proto/validate \
	--go_out=pb --go_opt=paths=source_relative \
	--go_opt=Muser.proto=github.com/pule1234/VideoForge/pb \
   	--go_opt=Mvideo.proto=github.com/pule1234/VideoForge/pb \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--validate_out=paths=source_relative,lang=go:pb \
	proto/*.proto