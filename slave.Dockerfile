FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD aimpanel-slave /
RUN chmod +x /aimpanel-slave
ENTRYPOINT ["/aimpanel-slave"]
EXPOSE 3000