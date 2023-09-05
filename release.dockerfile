# Builder image
FROM golang:latest as builder

WORKDIR /home/go/src/go-gatefuse

COPY . /home/go/src/go-gatefuse/

RUN go build -ldflags='-s -w' .

# Release image
FROM voidlinux/voidlinux:latest

WORKDIR /app

COPY --from=builder /home/go/src/go-gatefuse/go-gatefuse /app/go-gatefuse
COPY --from=builder /home/go/src/go-gatefuse/static /app/static
COPY --from=builder /home/go/src/go-gatefuse/templates /app/templates

RUN xbps-install -Syu xbps nginx nginx-mod-stream unbound && \
    ln -s /usr/lib/nginx/modules /etc/nginx/modules && ./go-gatefuse -init

COPY docker/etc/nginx/nginx.conf /etc/nginx/nginx.conf
# Load installed module for this container only
RUN sed -r -i 's|events|load_module modules/ngx_stream_module.so;\n\nevents|' /etc/nginx/nginx.conf

EXPOSE 3000

ENTRYPOINT [ "./go-gatefuse" ]