# builder image
FROM golang:1.16-alpine as builder
WORKDIR /app
RUN apk -U add git curl openssh-client
RUN curl -sLO https://github.com/argoproj/argo/releases/download/v3.0.7/argo-linux-amd64.gz && \
    # Unzip
    gunzip argo-linux-amd64.gz && \
    # Make binary executable
    chmod +x argo-linux-amd64 && \
    # Move binary to path
    mv ./argo-linux-amd64 /bin/argo
COPY webhook .
RUN CGO_ENABLED=0 go build -ldflags " -s -w -extldflags '-static'" -o /bin/webhook
# final image
FROM alpine:3.13
WORKDIR /hook

COPY --from=builder /bin/webhook /hook/webhook
COPY --from=builder /bin/argo /hook/argo
COPY --from=builder /app/argo.yml /hook/argo.yml
CMD ["/hook/webhook"]
