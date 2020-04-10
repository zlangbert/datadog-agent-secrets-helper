# builder image
FROM golang:1.14 as builder

WORKDIR /workspace

COPY . .
RUN make build

# final image
FROM busybox:latest
LABEL maintainer="Zach Langbert <zach.langbert@gmail.com>"

COPY --from=builder /workspace/build/datadog-agent-secrets-helper /bin/

RUN chmod 0500 /bin/datadog-agent-secrets-helper

ENTRYPOINT ["/bin/datadog-agent-secrets-helper"]