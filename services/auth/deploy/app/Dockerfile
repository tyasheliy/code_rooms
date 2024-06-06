FROM golang:1.22-alpine as build
WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -o ./build/app ./cmd/app/main.go

FROM scratch

COPY --from=build /build/build/app /usr/bin/app

ENTRYPOINT ["/usr/bin/app"]