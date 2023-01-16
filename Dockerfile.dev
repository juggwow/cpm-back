FROM golang:1.19 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .

RUN go mod download
RUN swag init
RUN go build -o /goapi
# RUN echo "/app -> $(ls -a /app)"

FROM gcr.io/distroless/base-debian10

COPY --from=builder /goapi /goapi

USER nonroot:nonroot

EXPOSE 80

CMD [ "/goapi" ]

