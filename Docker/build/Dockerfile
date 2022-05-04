FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o showtime main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/showtime /build/showtime

CMD [ "./showtime" ]
