FROM golang:alpine
#FROM golang:latest

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

LABEL maintainer="xxx"

WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.csv ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /scrape

# Export necessary port
EXPOSE 9999

# Command to run when starting the container

CMD ["/scrape","ipfs_cids.csv"]