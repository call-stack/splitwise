FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build cmd/splitwise_main.go


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/splitwise_main /splitwise_main
USER nonroot:nonroot

ENTRYPOINT ["/splitwise_main"]