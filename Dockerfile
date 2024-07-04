FROM ubuntu:22.04

ENV PATH="$PATH:/usr/local/go/bin"

# The SHA256 checksum used to verify the go archive can be found at https://go.dev/dl/

ENV GO_FILENAME=go1.22.5.linux-amd64.tar.gz
ENV GO_FILEHASH=904b924d435eaea086515bc63235b192ea441bd8c9b198c507e85009e6e4c7f0

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl make git zip \
  && curl -L https://go.dev/dl/${GO_FILENAME} >/tmp/${GO_FILENAME} \
  && echo ${GO_FILEHASH} /tmp/${GO_FILENAME} | sha256sum --check \
  && tar -C /usr/local -xzf /tmp/${GO_FILENAME}
