FROM golang:latest 

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY lego start.sh static/ tmpl/ /usr/local/bin/

RUN mkdir -p /go/src/github.com/dedgarsites/dedgar
WORKDIR /go/src/github.com/dedgarsites/dedgar

COPY . /go/src/github.com/dedgarsites/dedgar
RUN go-wrapper download && \
    go-wrapper install

EXPOSE 8443

USER 1001

CMD ["/usr/local/bin/start.sh"]
