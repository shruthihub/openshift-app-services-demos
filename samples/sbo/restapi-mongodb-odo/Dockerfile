FROM quay.io/redhat-developer/servicebinding-operator:builder-golang-1.16 AS builder

USER root

WORKDIR /workspace
COPY / /workspace/

# Build
RUN go build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.access.redhat.com/ubi8-minimal
WORKDIR /
COPY --from=builder /workspace/go-rest-mongodb .
COPY --from=builder /workspace/config.yml .
USER 65532:65532

ENTRYPOINT ["/go-rest-mongodb"]
