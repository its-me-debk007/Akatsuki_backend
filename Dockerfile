FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o akatsuki

EXPOSE 8080

CMD [ "./akatsuki" ]