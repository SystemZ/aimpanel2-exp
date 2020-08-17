# minimal OS for minimal attack surface and easy to use TLS and timezones
FROM alpine:latest

# run all installs and filesystem adds as root
USER root

# install OS dependencies
# add non-root group and user account called "go"
RUN apk add --no-cache ca-certificates bash \
 && mkdir /exp \
 && groupadd -r go \
 && useradd --no-log-init -r -g go go

# extras for main binary
ADD ../swagger.json /exp/swagger.json
ADD ../master/redoc.html /exp/redoc.html
ADD ../master/swagger.html /exp/swagger.html

# frontend dir to be served by main binary
ADD aimpanel-master-frontend /exp/frontend

# slave binary as a main or backup source of update and deployment of new hosts
ADD aimpanel-slave /exp/slave

# add main binary
ADD aimpanel-master /exp/master
# set proper permissions for binary and extras
RUN chmod +x /exp/master \
 && chown -r go:root /exp

# list files added to container as debug
RUN ls -alh /exp/

# run main binary as non-root in production
USER go
ENTRYPOINT ["/exp/master"]