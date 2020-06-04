FROM golang:latest
WORKDIR /go/src/github.com/drags/tagger-backend/
COPY *.go ./
RUN go get github.com/drags/pinboard
RUN go get github.com/gorilla/handlers
RUN go get github.com/rs/cors
RUN CGO_ENABLED=0 go build

FROM alpine:latest
COPY --from=0 /go/src/github.com/drags/tagger-backend/tagger-backend /usr/local/bin/
WORKDIR /root/
CMD ["/usr/local/bin/tagger-backend"]
