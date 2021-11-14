FROM golang:alpine3.14 AS build

COPY . /root/hznuoj-exporter
WORKDIR /root

RUN cd hznuoj-exporter \
	&& go build hznuoj_exporter \
	&& cp hznuoj_exporter /root/hznuoj_exporter \
	&& rm -rf /root/hznuoj-exporter

FROM alpine:3.14

COPY --from=build /root/hznuoj_exporter /hznuoj_exporter

EXPOSE 9800
ENTRYPOINT  [ "/hznuoj_exporter" ]
