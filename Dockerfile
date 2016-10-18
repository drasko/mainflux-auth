FROM golang:1.7-alpine
MAINTAINER Mainflux

# copy the local package files into the container's workspace
COPY . /go/src/github.com/mainflux/mainflux-auth-server

# build the service inside the container
RUN go install github.com/mainflux/mainflux-auth-server

# specify the service entrypoint
CMD ["/go/bin/mainflux-auth-server"]

# document exposed ports
EXPOSE 8180
