FROM ubuntu:22.04

ENV PATH="$PATH:/usr/local/go/bin"

# The SHA256 checksum used to verify the go archive can be found at https://go.dev/dl/

ENV GO_FILENAME=go1.22.3.linux-amd64.tar.gz
ENV GO_FILEHASH=8920ea521bad8f6b7bc377b4824982e011c19af27df88a815e3586ea895f1b36

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl make git zip \
  && curl -L https://go.dev/dl/${GO_FILENAME} >/tmp/${GO_FILENAME} \
  && echo ${GO_FILEHASH} /tmp/${GO_FILENAME} | sha256sum --check \
  && tar -C /usr/local -xzf /tmp/${GO_FILENAME}
