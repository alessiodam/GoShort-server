FROM golang:1.22.5 AS builder
LABEL author="TKB Studios"
LABEL maintainer="TKB Studios"
LABEL personal-contact="tkbstudios@mail.tkbstudios.com"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goshort-server .

FROM alpine:latest
LABEL author="TKB Studios"
LABEL maintainer="TKB Studios"
LABEL personal-contact="tkbstudios@mail.tkbstudios.com"

WORKDIR /app

COPY --from=builder /app/goshort-server .
COPY --from=builder /app/entrypoint.sh .

RUN chmod +x /app/entrypoint.sh

EXPOSE 8000

CMD ["/app/entrypoint.sh"]
