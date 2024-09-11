FROM golang:alpine AS builder 
WORKDIR /apps

COPY go.mod .
COPY go.sum .
COPY . .

RUN go get -u -d -v
RUN go build -o randombot

FROM apline:latest
WORKDIR /app

COPY --from=builder /apps/randombot ./
CMD [ "./randombot" ]