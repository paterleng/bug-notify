FROM golang:1.18 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY config/ config/

RUN CGO_ENABLED=0 GOOS=linux  GO111MODULE=on go build -a -o manager main.go


FROM alpine:3.11.2
WORKDIR /
COPY --from=builder /workspace/manager .
CMD ["/manager"]