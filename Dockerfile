FROM docker.io/golang:1.19 AS build

WORKDIR /go/task-manager

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg/ pkg/
COPY vendor/ vendor/

RUN CGO_ENABLED=0 go test --cover ./...

RUN CGO_ENABLED=0 go build -o bin/task-manager

FROM gcr.io/distroless/static

COPY --from=build /go/task-manager/bin/task-manager /usr/local/bin/task-manager

WORKDIR /usr/local/bin

CMD ["task-manager"]

USER nonroot:nonroot