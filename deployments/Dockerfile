FROM golang:1.24.4 AS builder

WORKDIR /app
COPY . .

RUN rm -rf .git

RUN GOARCH=wasm GOOS=js go build -o web/app.wasm cmd/frontend/main.go 
RUN CGO_ENABLED=0 GOOSE=linux GO111MODULE=on go build -a -installsuffix nocgo -o accounter cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/ ./
CMD ["./accounter"]

