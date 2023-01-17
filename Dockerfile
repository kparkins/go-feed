FROM golang:1.18-alpine as build

WORKDIR /app/
COPY . /app/
RUN CGO_ENABLED=0 go build -o /bin/go-message

FROM scratch
COPY  --from=build /bin/go-message /bin/go-message
ENTRYPOINT ["/bin/go-message"]
