FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -buildvcs=false -o ./server ./cmd/server

# Copy existing TLS certificates
RUN cp -r cmd/server/certs ./certs

EXPOSE 50051

CMD ["./server"]
