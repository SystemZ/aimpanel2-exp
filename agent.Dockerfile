FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD aimpanel-agent /
RUN chmod +x /aimpanel-agent
ENTRYPOINT ["/aimpanel-agent"]
EXPOSE 3000