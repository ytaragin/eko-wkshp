FROM golang:1.18


#COPY ./tasks /usr/local/go/src/tasks/tasks


WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
#COPY tasks/*.go ./tasks/
RUN go build -v -o /usr/local/bin ./...

CMD ["protection"]