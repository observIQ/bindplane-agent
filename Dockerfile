# Golang has multi arch manifests for amd64 and arm64
FROM golang:1.17 as build
WORKDIR /collector
COPY . /collector
ARG JMX_JAR_VERSION=v1.7.0
RUN make collector
RUN curl -L \
    --output /opt/opentelemetry-java-contrib-jmx-metrics.jar \
    "https://github.com/open-telemetry/opentelemetry-java-contrib/releases/download/${JMX_JAR_VERSION}/opentelemetry-jmx-metrics.jar"

FROM gcr.io/observiq-container-images/stanza-base:v1.1.0
WORKDIR /

# config directory
RUN mkdir -p /etc/otel

# copy binary with unpredictable due to dynamic GOOS / GOARCH
COPY --from=build /collector/build/* /

# copy jmx receiver dependency
COPY --from=build /opt/opentelemetry-java-contrib-jmx-metrics.jar /opt/opentelemetry-java-contrib-jmx-metrics.jar

# User should mount /etc/otel/config.yaml at runtime using docker volumes / k8s configmap
ENTRYPOINT [ "/collector" ]
CMD ["--config", "/etc/otel/config.yaml"]