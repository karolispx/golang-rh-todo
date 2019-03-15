FROM centos:latest

# gcc for cgo
RUN yum install -y git g++ \
	gcc \
	libc6-dev \
	make \
	telnet \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.11
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 b3fcf280ff86558e0559e185b601c9eade0fd24c900b4c63cd14d1d38613e499

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
COPY bin/* /go/

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 0755 "$GOPATH"
WORKDIR $GOPATH

USER 1001

ENTRYPOINT [ "./uid_entrypoint.sh" ]

# This will change depending on each golang-rh-todo entry point
CMD ["./golang-rh-todo"]
