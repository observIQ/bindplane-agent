#!/bin/sh
set -e

[ -f go-msi.exe ] || curl -f -L -o go-msi.exe https://github.com/observIQ/go-msi/releases/download/v2.0.0/go-msi.exe
[ -f ./cinc-auditor.msi ] || curl -f -L -o cinc-auditor.msi http://downloads.cinc.sh/files/stable/cinc-auditor/4.17.7/windows/2012r2/cinc-auditor-4.17.7-1-x64.msi

[ -f ./wix-binaries.zip ] || curl -f -L -o wix-binaries.zip https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip

mkdir -p wix
[ -d wix/sdk ] || unzip -o wix-binaries.zip -d wix

[ -f ./opentelemetry-java-contrib-jmx-metrics.jar ] || curl -f -L -o ./opentelemetry-java-contrib-jmx-metrics.jar \
    "https://github.com/open-telemetry/opentelemetry-java-contrib/releases/download/$(cat ../../JAVA_CONTRIB_VERSION)/opentelemetry-jmx-metrics.jar"

if [ ! -d "./stanza-plugins" ]; then
    git clone git@github.com:observIQ/stanza-plugins.git
fi

cd stanza-plugins
git fetch --all
git checkout "$(cat ../../../PLUGINS_VERSION)"
cd ..

cp -r ./stanza-plugins/plugins .

cp ../../config/example.yaml ./config.yaml
