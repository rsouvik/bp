FROM golang:alpine
#FROM golang:latest

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

LABEL maintainer="xxx"

RUN mkdir -p /data

WORKDIR /build

ARG DATA_DIR=/build/data
ENV GOPATH  "/build/data"
ENV PATH "${PATH}:/build/data/bin"

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o bp .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/bp .

VOLUME ${DATA_DIR}

RUN echo "export PATH=$PATH" > /etc/environment
RUN echo "export GOPATH=$GOPATH" > /etc/environment

# Export necessary port
EXPOSE 9999

# Command to run when starting the container
#CMD ["/dist/iotProxySvc"]
#ENTRYPOINT ["/dist/bp"]
CMD ["/dist/bp","ipfs_cids.csv"]