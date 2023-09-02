# Builder image
FROM golang:latest as builder

WORKDIR /home/go/src/go-gatefuse

COPY . /home/go/src/go-gatefuse/

RUN go build -ldflags='-s -w' .

# Release image
FROM voidlinux/voidlinux:latest

WORKDIR /app

RUN xbps-install -Syu xbps nginx unbound

COPY --from=builder /home/go/src/go-gatefuse/go-gatefuse /app/go-gatefuse
COPY --from=builder /home/go/src/go-gatefuse/static /app/static
COPY --from=builder /home/go/src/go-gatefuse/templates /app/templates

EXPOSE 3000

ENTRYPOINT [ "./go-gatefuse" ]