FROM alpine:3.6

RUN apk add --no-cache ca-certificates

ADD /bin/cloudscale-cloud-controller-manager /bin/

CMD ["/bin/cloudscale-cloud-controller-manager"]
