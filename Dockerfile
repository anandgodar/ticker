FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN cd /go/src
RUN mkdir -p /go/src/projects
RUN cp -R /app/* /go/src/projects
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main /go/src/projects/cmd/main.go
CMD ["./main"]
