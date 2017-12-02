FROM openshift/base-centos7 

RUN yum install -y golang && \
    yum clean all

ENV GOLANG_VERSION 1.9 \
    GOPATH /go \
    PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY go-wrapper /usr/local/bin/

RUN mkdir -p /go/src/github.com/openshift/shinodev
WORKDIR /go/src/github.com/openshift/shinodev

COPY . /go/src/github.com/openshift/shinodev
RUN go-wrapper download && go-wrapper install

EXPOSE 8080

USER 1001

CMD ["go-wrapper", "run"]
