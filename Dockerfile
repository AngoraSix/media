# Builder Stage
FROM golang:latest as builder
WORKDIR /go/src/angorasix.com/media
COPY . .
ARG BUILD
ENV GO111MODULE=on
RUN go install goa.design/goa/v3/cmd/goa@v3
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go get -u goa.design/goa/v3
RUN go get -u goa.design/goa/v3/...
RUN goa gen angorasix.com/media/design
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o media-svc -ldflags "-X main.Build=$BUILD" cmd/main/*
# Runner Stage
FROM centurylink/ca-certs as runner
WORKDIR /root/
COPY --from=builder /go/src/angorasix.com/media .
ENTRYPOINT [ "./media-svc" ]
EXPOSE 80
CMD []