FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD aimpanel-master /
RUN chmod +x /aimpanel-master
ENTRYPOINT ["/aimpanel-master"]
EXPOSE 3000