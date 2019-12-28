FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD aimpanel-master /
ADD swagger.json /
ADD master/redoc.html /
ADD master/swagger.html /
RUN chmod +x /aimpanel-master
ENTRYPOINT ["/aimpanel-master"]
EXPOSE 3000