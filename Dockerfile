FROM alpine:latest

LABEL org.opencontainers.image.source=https://github.com/Bedrock-Technology/lambda

WORKDIR /app

COPY ./lambda /app/lambda

ENTRYPOINT ["/app/lambda"]
