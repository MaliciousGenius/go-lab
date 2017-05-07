FROM golang:alpine

RUN mkdir -p /opt/app
ADD hello /opt/app/
WORKDIR /opt/app
EXPOSE 8000
ENTRYPOINT /opt/app/hello