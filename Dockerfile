# minimal OS for minimal attack surface and easy to use TLS and timezones
FROM alpine:latest

# run all installs and filesystem adds as root
USER root

# install OS dependencies
# add non-root group and user account called "go"
RUN apk add --no-cache ca-certificates bash \
 && mkdir -p /exp/docs \
 # alpine way https://stackoverflow.com/questions/49955097/how-do-i-add-a-user-when-im-using-alpine-as-a-base-image
 && addgroup -S go && adduser -S go -G go
 # debian/ubuntu way
 #&& groupadd -r go && useradd --no-log-init -r -g go go

# extras for main binary
ADD swagger.json /exp/docs/swagger.json
ADD master/redoc.html /exp/docs/redoc.html
ADD master/swagger.html /exp/docs/swagger.html

# frontend dir to be served by main binary
ADD aimpanel-master-frontend /exp/frontend

# slave binary as a main or backup source of update and deployment of new hosts
ADD aimpanel-slave /exp/slave

# add main binary
ADD aimpanel-master /exp/master
# set proper permissions for binary and extras
RUN chmod +x /exp/master \
 && chown -R go:root /exp

# list files added to container as debug
RUN ls -alh /exp/

# run main binary as non-root in production
USER go
ENTRYPOINT ["/exp/master"]