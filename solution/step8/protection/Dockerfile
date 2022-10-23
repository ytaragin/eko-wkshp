FROM golang:1.18 AS builder
WORKDIR /workdir
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o protection
FROM scratch
COPY --from=builder /workdir/protection /
CMD ["/protection"]

