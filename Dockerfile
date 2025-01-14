# s6 overlay builder
FROM azinchen/nordvpn AS nordvpn

FROM golang:alpine AS builder

COPY . /src/go

WORKDIR /src/go

RUN mkdir -p /src/go/root/usr/bin ; \
    mkdir -p /src/go/root/etc/nordvpn ; \
    go generate ./... ; \
    go build -o /src/go/root/usr/bin cmd/go-nordvpn.go

COPY --from=nordvpn /etc/nordvpn/template.ovpn /src/go/root/etc/nordvpn

RUN  sed -e 's/__PROTOCOL__/{{.Protocol}}/' \
        -e 's/__IP__/{{.IP}}/' \
        -e 's/__PORT__/{{.Port}}/' \
        -e 's/__HOSTNAME__/{{.Hostname}}/' \
        -e 's/__X509_NAME__/{{.Hostname}}/' \
        root/etc/nordvpn/template.ovpn \
        >root/etc/nordvpn/template.ovpn.tmpl

FROM nordvpn

COPY --from=builder /src/go/root/ /

ENTRYPOINT ["/init"]
