FROM golang:alpine as builder

WORKDIR /app 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /plexer

FROM scratch 

COPY --from=builder /plexer /plexer

ENTRYPOINT [ "/plexer" ]
