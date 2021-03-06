FROM golang:1.15.6-alpine3.12 as builder
RUN mkdir /app
COPY ./main.go /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

#scratch
FROM scratch
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]


# FROM golang:1.15.6-alpine3.12 as builder
# RUN mkdir /app
# COPY ./main.go /app/
# WORKDIR /app
# RUN go build -o server .
#
# #scratch  alpine:latest
# FROM golang:1.15.6-alpine3.12
# WORKDIR /root/
# COPY --from=builder /app/server .
# CMD ["./server"]
