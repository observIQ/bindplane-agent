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

# Golang has multi arch manifests for amd64 and arm64
FROM golang:1.17 as build
WORKDIR /collector
COPY . /collector
ARG JMX_JAR_VERSION=v1.7.0
ARG GITHUB_TOKEN
RUN \
    make install-tools && \
    goreleaser build --single-target --skip-validate --rm-dist

# Find built executable, there is only one, and copy it to working dir
RUN find /collector/dist -name observiq-otel-collector -exec cp {} . \;

RUN curl -L \
    --output /opt/opentelemetry-java-contrib-jmx-metrics.jar \
    "https://github.com/open-telemetry/opentelemetry-java-contrib/releases/download/${JMX_JAR_VERSION}/opentelemetry-jmx-metrics.jar"

# Official OpenJDK has multi arch manifests for amd64 and arm64
# Java is required for JMX receiver
# Contrib's integration tests use openjdk 1.8.0
# https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/jmxreceiver/testdata/Dockerfile.cassandra
FROM openjdk:8u312-slim-buster as openjdk


FROM gcr.io/observiq-container-images/stanza-base:v1.1.0
WORKDIR /

# configure java runtime
COPY --from=openjdk /usr/local/openjdk-8 /usr/local/openjdk-8
ENV JAVA_HOME=/usr/local/openjdk-8
ENV PATH=$PATH:/usr/local/openjdk-8/bin

# config directory
RUN mkdir -p /etc/otel

# copy binary
COPY --from=build /collector/observiq-otel-collector /collector/

# copy jmx receiver dependency
COPY --from=build /opt/opentelemetry-java-contrib-jmx-metrics.jar /opt/opentelemetry-java-contrib-jmx-metrics.jar

# User should mount /etc/otel/config.yaml at runtime using docker volumes / k8s configmap
ENTRYPOINT [ "/collector" ]
CMD ["--config", "/etc/otel/config.yaml"]
