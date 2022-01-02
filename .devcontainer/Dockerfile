FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu-21.04

ARG GOLANG_VERSION="1.17.5"
ARG LOCALE="en_GB.UTF-8"

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get update \
    && apt-get -y install --no-install-recommends \
        build-essential fish nodejs npm \
    && apt-get autoremove -y \
    && apt-get clean autoclean -y \
    && rm -r /var/cache/* /var/lib/apt/lists/*

RUN locale-gen $LOCALE \
    && echo LANG=$LOCALE > /etc/default/locale

RUN npm install snyk@latest -g

USER vscode

ENV HOME="/home/vscode"
ENV GOENV_ROOT="$HOME/.goenv"
ENV GOROOT="$GOENV_ROOT/versions/$GOLANG_VERSION"
ENV GOPATH="$HOME/go/$GOLANG_VERSION"
ENV LC_ALL=$LOCALE
ENV PATH="$GOENV_ROOT/bin:$HOME/.local/bin:$PATH"

WORKDIR $HOME

RUN git clone https://github.com/syndbg/goenv.git $HOME/.goenv \
    && eval "$(goenv init -)" \
    && goenv install $GOLANG_VERSION \
    && goenv rehash \
    && goenv global $GOLANG_VERSION \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@2f13672765fe \
    && go install github.com/uudashr/gopkgs/cmd/gopkgs@v2 \
    && go install github.com/ramya-rao-a/go-outline@latest \
    && go get -u github.com/haya14busa/goplay/cmd/goplay \
    && go install github.com/fatih/gomodifytags@latest \
    && go install github.com/josharian/impl@latest \
    && go get -u github.com/cweill/gotests/... \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && git clone https://github.com/magefile/mage /tmp/mage \
    && cd /tmp/mage \
    && go run bootstrap.go \
    && cd - \
    && rm -rf tmp/mage

RUN echo "eval \"\$($GOENV_ROOT/bin/goenv init -)\"" >> $HOME/.profile \
    && echo "eval \"\$($GOENV_ROOT/bin/goenv init -)\"" >> $HOME/.bashrc \
    && echo "eval \"\$($GOENV_ROOT/bin/goenv init -)\"" >> $HOME/.zshrc
