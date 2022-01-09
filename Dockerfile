##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY * ./
COPY decisionbot/ decisionbot/

RUN go mod download

# Build app
RUN CGO_ENABLED=0 go build -o /build/decisionbot main.go

# Create passwd file with user nobody
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

##
## Deploy
##
FROM scratch

# Copy app
COPY --from=build /build/decisionbot decisionbot

# Copy passwd with user nobody
COPY --from=build /etc_passwd /etc/passwd

# Copy SSL certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nobody

ENTRYPOINT ["/decisionbot"]