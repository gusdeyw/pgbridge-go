FROM golang:alpine AS builder
RUN apk add --no-cache build-base
WORKDIR /app
ENV config=docker
COPY . .
RUN go mod download && go mod tidy
RUN go build -o main ./main.go

FROM golang:alpine AS runner
WORKDIR /app
ENV config=docker
COPY --from=builder /app ./
RUN apk add --no-cache build-base
RUN go install -mod=mod github.com/githubnemo/CompileDaemon@latest
EXPOSE 5000
ENTRYPOINT ["/go/bin/CompileDaemon", "--build=go build -o main ./main.go", "--command=./main"]
