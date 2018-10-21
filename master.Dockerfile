FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD bin/aimpanel-master /
RUN chmod +x /aimpanel-master
ENTRYPOINT ["/aimpanel-master"]
EXPOSE 3000