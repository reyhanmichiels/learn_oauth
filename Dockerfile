FROM golang:1.21.0-alpine3.18 as stage
WORKDIR /app/
COPY . .
RUN go mod download
RUN go build -o binary ./cmd

FROM alpine:3.18
WORKDIR /app/
COPY --from=stage /app/ .
CMD /app/binary