FROM golang:alpine

LABEL maintainer="Debashish Kundu"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o akatsuki

EXPOSE 8081

CMD [ "./akatsuki" ]