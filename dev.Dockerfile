FROM golang:1.22-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o /build/bin ./cmd/." -command=/build/bin