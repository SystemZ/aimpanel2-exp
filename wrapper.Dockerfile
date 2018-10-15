FROM alpine:latest
RUN apk add --no-cache ca-certificates bash
ADD aimpanel-wrapper /
RUN chmod +x /aimpanel-wrapper
ENTRYPOINT ["/aimpanel-wrapper"]
EXPOSE 3000