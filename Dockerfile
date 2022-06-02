# FROM alpine

# COPY drone-discord /bin/
# RUN apk -Uuv add ca-certificates

# ENTRYPOINT ["/bin/drone-discord"]

FROM plugins/base:linux-amd64

LABEL maintainer="Yury Kotov <beorc@httplab.ru>" \
  org.label-schema.name="Drone Discord" \
  org.label-schema.vendor="HttpLab" \
  org.label-schema.schema-version="1.0"

COPY drone-discord /bin/

ENTRYPOINT ["/bin/drone-discord"]
