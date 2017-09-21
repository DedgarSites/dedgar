FROM openshift/base-centos7 

RUN yum install -y golang && \
    yum clean all

ENV GOLANG_VERSION 1.9

RUN mkdir -p /go/src/github.com/openshift/shinodev
WORKDIR /go/src/github.com/openshift/shinodev

COPY . /go/src/github.com/openshift/shinodev
RUN go-wrapper download && go-wrapper install

ENV GOPATH /go 

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY go-wrapper /usr/local/bin/

USER 1001

CMD ["go-wrapper", "run"]
