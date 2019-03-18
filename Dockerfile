# builder image
FROM golang as builder

WORKDIR /workspace

COPY . .
RUN make build

# final image
FROM busybox:latest
LABEL maintainer="Zach Langbert <zach.langbert@gmail.com>"

COPY --from=builder /workspace/build/datadog-secrets-provider-aws-secretsmanager /bin/

RUN chmod 0700 /bin/datadog-secrets-provider-aws-secretsmanager

ENTRYPOINT ["/bin/datadog-secrets-provider-aws-secretsmanager"]