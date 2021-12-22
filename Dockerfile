# Builder Stage
FROM golang:latest as builder
WORKDIR /go/src/angorasix.com/media
COPY . .
ARG BUILD
ENV GO111MODULE=on
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