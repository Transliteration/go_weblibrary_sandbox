FROM golang:1.18 AS builder
RUN apt-get update -y && apt-get upgrade -y
WORKDIR /grpc_service
# COPY ./services/gad-manager/ ./services/gad-manager/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
# COPY . .
COPY ./database ./database
COPY ./server ./server
COPY ./grpc_server ./grpc_server

RUN cd server && CGO_ENABLED=0 go build -a -installsuffix cgo

FROM alpine:latest AS runner
RUN apk -U upgrade
WORKDIR /grpc_service
COPY --from=builder /grpc_service/server/ ./
EXPOSE 50051
EXPOSE 3306
ENTRYPOINT [ "./server" ]