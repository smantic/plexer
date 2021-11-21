FROM golang:alpine as build

RUN apk --no-cache add ca-certificates
WORKDIR /app 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /plexer

FROM scratch 
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /plexer / 

ENTRYPOINT [ "/plexer" ]
