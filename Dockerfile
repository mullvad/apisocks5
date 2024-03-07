FROM ubuntu:22.04

ENV PATH="$PATH:/usr/local/go/bin"

# The SHA256 checksum used to verify the go archive can be found at https://go.dev/dl/

ENV GO_FILENAME=go1.22.2.linux-amd64.tar.gz
ENV GO_FILEHASH=5901c52b7a78002aeff14a21f93e0f064f74ce1360fce51c6ee68cd471216a17

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl make git zip \
  && curl -L https://go.dev/dl/${GO_FILENAME} >/tmp/${GO_FILENAME} \
  && echo ${GO_FILEHASH} /tmp/${GO_FILENAME} | sha256sum --check \
  && tar -C /usr/local -xzf /tmp/${GO_FILENAME}
