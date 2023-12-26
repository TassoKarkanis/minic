# syntax=docker/dockerfile:1.3.1

###########################################################################
# minic-builder
###########################################################################

FROM ubuntu:22.04 as minic-builder

# install some dependencies
# Note: to suppress interactive menu of tzdata, DEBIAN_FRONTEND must be set
RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean && \
    apt-get update && \
    apt-get install -y \
		antlr4 \
		gcc \
		git \
		make \
		nasm \
		wget

# Download and install Go
RUN wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz && \
tar -xf go1.19.1.linux-amd64.tar.gz

# Configure Go
ENV GOROOT=/go
ENV GOCACHE=/root/go/cache
ENV PATH=$GOROOT/bin:$PATH
# Diable buildvcs flag to build on CircleCI [See https://github.com/ko-build/ko/issues/672 and ]
ENV GOFLAGS="-buildvcs=false"


###########################################################################
# minic-devcontainer
###########################################################################

FROM minic-builder AS minic-devcontainer

# install VS Code (code-server)
RUN curl -fsSL https://code-server.dev/install.sh | sh

# # install VS Code extensions
# RUN code-server --install-extension redhat.vscode-yaml \
#                 --install-extension golang.Go
# 
# # make git work in /w
# RUN git config --global --add safe.directory /w
