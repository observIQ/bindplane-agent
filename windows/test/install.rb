collector_home="C:/Program Files/observiq-otel-collector"

[
    "#{collector_home}/plugins"
].each do |dir|
    describe file(dir) do
        it { should exist }
        it { should be_directory }
    end
end

[
    "#{collector_home}/observiq-otel-collector.exe",
    "#{collector_home}/config.yaml",
    "#{collector_home}/plugins/aerospike.yaml",
    "#{collector_home}/plugins/microsoft_iis.yaml",
    "#{collector_home}/plugins/zookeeper.yaml"
].each do |file|
    describe file(file) do
        it { should exist }
        it { should be_file }
    end
end

describe service('observiq-otel-collector') do
    it { should be_installed }
    it { should be_enabled }
    it { should be_running }
end
