FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN mkdir /var/kubeStone
COPY . .
RUN mv  install/* /var/kubeStone/
RUN rm -rf front-end

RUN go build -o kubestone .

EXPOSE 8888

CMD ["./kubestone"]

