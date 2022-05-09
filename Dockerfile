# begin build container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as build

RUN microdnf install -y golang

ENV GOBIN=/bin \
    GOPATH=/go

RUN /usr/bin/go install github.com/dedgarsites/dedgar@master


# begin run container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as run

ADD scripts/ /usr/local/bin/
ADD static/ /usr/local/bin/static
ADD tmpl/ /usr/local/bin/tmpl

RUN mkdir /cert

COPY --from=build /bin/dedgar /usr/local/bin

EXPOSE 8443

CMD /usr/local/bin/start.sh
