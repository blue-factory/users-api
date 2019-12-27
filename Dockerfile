# Base container for compile service
FROM golang:alpine AS builder

# Install dependencies
RUN apk add make

# Go to builder workdir
WORKDIR /go/src/github.com/microapis/users-api/

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

# Go to workdir
WORKDIR /src/users-api

# Install dependencies
RUN apk add --update ca-certificates wget

# Copy binaries
COPY --from=builder /go/src/github.com/microapis/users-api/bin/users-api /usr/bin/users-api

# # Copy goose migration tool and add permission
# RUN wget https://raw.githubusercontent.com/microapis/lib/master/bin/goose -o /usr/bin/goose
# RUN chmod 777 /usr/bin/goose

# # Copy wait-db util and add permission
# RUN wget https://raw.githubusercontent.com/microapis/lib/master/bin/wait-db -o /usr/bin/wait-db
# RUN chmod 777 /usr/bin/wait-db

COPY bin/goose /usr/bin/goose
COPY bin/wait-db /usr/bin/wait-db

# Copy all database migrations
COPY database/migrations/* /src/users-api/migrations/

# Expose service port
EXPOSE 5020

# Run service
CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/users-api/migrations/ && goose postgres ${POSTGRES_DSN} up && users-api"]