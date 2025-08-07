FROM golang:1.24-alpine AS dependencies
RUN apk add make
LABEL authors="skamenetskiy"
WORKDIR /app
ADD api api
ADD go.mod .
ADD go.sum .
RUN go mod download

FROM dependencies AS builder
ADD . .
RUN make build-static
RUN ls -lah bin/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/messages /messages
CMD ["/messages"]