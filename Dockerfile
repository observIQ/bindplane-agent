# Copyright  observIQ, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Build stage builds the release version of observiq-otel-collector
# and downloads the opentelemetry-jmx-metrics.jar used by JMX receivers
#
FROM golang:1.17 as build
WORKDIR /collector
COPY . /collector
ARG JMX_JAR_VERSION=v1.7.0

COPY dist dist

RUN cp "dist/collector_linux_$(go env GOARCH)/observiq-otel-collector" .

RUN curl -L \
    --output /opt/opentelemetry-java-contrib-jmx-metrics.jar \
    "https://github.com/open-telemetry/opentelemetry-java-contrib/releases/download/${JMX_JAR_VERSION}/opentelemetry-jmx-metrics.jar"


# OpenJDK stage provides the Java runtime used by JMX receivers.
# Contrib's integration tests use openjdk 1.8.0
# https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/jmxreceiver/testdata/Dockerfile.cassandra
#
FROM openjdk:8u312-slim-buster as openjdk


# Final Stage
#
FROM gcr.io/observiq-container-images/stanza-base:v1.2.0
WORKDIR /

COPY --from=openjdk /usr/local/openjdk-8 /usr/local/openjdk-8
ENV JAVA_HOME=/usr/local/openjdk-8
ENV PATH=$PATH:/usr/local/openjdk-8/bin

COPY --from=build /collector/observiq-otel-collector /collector/
COPY --from=build /opt/opentelemetry-java-contrib-jmx-metrics.jar /opt/opentelemetry-java-contrib-jmx-metrics.jar

RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home \
    --uid 10005 \
    otel

USER otel

# User should mount /etc/otel/config.yaml at runtime using docker volumes / k8s configmap
ENTRYPOINT [ "/collector/observiq-otel-collector" ]
CMD ["--config", "/etc/otel/config.yaml"]
