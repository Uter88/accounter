killall accounter -q
GOARCH=wasm GOOS=js go build -o web/app.wasm cmd/frontend/main.go 
go build -o accounter cmd/backend/main.go
./accounter -mode=dev