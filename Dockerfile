FROM alpine:latest

ADD bin/qwwebhook /qwwebhook
#ENTRYPOINT ["/qwwebhook"]
CMD ["/qwwebhook"]