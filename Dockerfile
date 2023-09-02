FROM nginx:latest
COPY --from=golang:latest /usr/local/go /usr/local/go
# Configure Nginx
COPY docker/etc/nginx/nginx.conf /etc/nginx/nginx.conf
RUN mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.disabled
# Install dev packages
RUN apt update && apt install -y gcc g++ git

ENV PATH=${PATH}:/usr/local/go/bin:/opt/go/bin EDITOR=nano

WORKDIR /opt/go

RUN go env -w GOPATH=/opt/go && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest