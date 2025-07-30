#####################################
#          Builder image            #
#####################################
FROM golang:1.24 AS builder

WORKDIR /message-processor

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./cmd ./cmd
COPY ./configs ./configs
COPY ./internal ./internal
COPY ./schemas ./schemas

# Download modules and delete the ssh key
RUN go mod download

WORKDIR /message-processor/cmd/message-processor

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o message-processor

#####################################
#        Deployable image           #
#####################################
FROM alpine:3.19

RUN addgroup -S message-processor && adduser -S message-processor -G message-processor -h /home/message-processor
USER message-processor

WORKDIR /home/message-processor

# Copy the Pre-built binary file from the previous stage
COPY --from=builder message-processor/cmd/message-processor/message-processor .
COPY --from=builder message-processor/schemas ./schemas

CMD ./message-processor
