FROM alpine:latest
WORKDIR /
COPY aprsweb-native /
EXPOSE 5000
CMD ["/aprsweb-native"]
