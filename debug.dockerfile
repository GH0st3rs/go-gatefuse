FROM voidlinux/voidlinux
COPY --from=golang:latest /usr/local/go /usr/local/go
# Install pretties
RUN xbps-install -Syu xbps ncurses-term zsh zsh-autosuggestions zsh-completions grml-zsh-config && chsh -s /bin/zsh
# Set environment variables
ENV TERM=xterm-256color SHELL=/bin/zsh PATH=${PATH}:/usr/local/go/bin:/opt/go/bin

RUN xbps-install -Sy nginx nginx-mod-stream git tree gcc

# Configure Nginx
COPY docker/etc/nginx/nginx.conf /etc/nginx/nginx.conf

WORKDIR /opt/go

RUN go env -w GOPATH=/opt/go && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest