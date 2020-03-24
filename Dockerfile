# Base container for compile service
FROM golang:alpine AS builder

# Define service name
ARG SVC=users-api

# Install dependencies
RUN apk add make

# Go to builder workdir
WORKDIR /go/src/github.com/microapis/${SVC}/

# Copy go modules files
COPY go.mod .
COPY go.sum .

# Install dependencies
RUN go mod download

# Copy all source code
COPY . .


# Compile service
RUN make linux

#####################################################################
#####################################################################

# Base container for run service
FROM alpine

# Define service name
ARG SVC=users-api

# Go to workdir
WORKDIR /src/${SVC}

# Install dependencies
RUN apk add --update ca-certificates wget

# Copy binaries
COPY --from=builder /go/src/github.com/microapis/${SVC}/bin/${SVC} /usr/bin/${SVC}

COPY bin/goose /usr/bin/goose
COPY bin/wait-db /usr/bin/wait-db

# Copy all database migrations
COPY database/migrations/* /src/${SVC}/migrations/

# Expose service port
EXPOSE 5020

# Run service
CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/$SVC/migrations/ && goose postgres ${POSTGRES_DSN} up && $SVC"]