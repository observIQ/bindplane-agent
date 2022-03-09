describe file('/opt/observiq-collector') do
    its('mode') { should cmp '0755' }
    its('owner') { should eq 'observiq-collector' }
    its('group') { should eq 'observiq-collector' }
    its('type') { should cmp 'directory' }
end

[
    '/opt/observiq-collector/plugins',
].each do |dir|
    describe file(dir) do
        its('mode') { should cmp '0750' }
        its('owner') { should eq 'observiq-collector' }
        its('group') { should eq 'observiq-collector' }
        its('type') { should cmp 'directory' }
    end
end

[
    '/opt/observiq-collector/observiq-collector',
    '/opt/observiq-collector/opentelemetry-java-contrib-jmx-metrics.jar'
].each do |bin|
    describe file(bin) do
        its('mode') { should cmp '0755' }
        its('owner') { should eq 'observiq-collector' }
        its('group') { should eq 'observiq-collector' }
        its('type') { should cmp 'file' }
    end
end

[
    '/opt/observiq-collector/config.yaml',
].each do |config|
    describe file(config) do
        its('mode') { should cmp '0660' }
        its('owner') { should eq 'observiq-collector' }
        its('group') { should eq 'observiq-collector' }
        its('type') { should cmp 'file' }
    end
end

source_dir = 'release_deps/plugins/'
Find.find(source_dir) do |path|
    plugin_file = path.delete_prefix(source_dir)
    path = "/opt/observiq-collector/plugins/#{plugin_file}"
    if path == '/opt/observiq-collector/plugins/' then
        describe file(path) do
            its('mode') { should cmp '0750' }
            its('owner') { should eq 'observiq-collector' }
            its('group') { should eq 'observiq-collector' }
        end
    else
        describe file(path) do
            its('mode') { should cmp '0640' }
            its('owner') { should eq 'observiq-collector' }
            its('group') { should eq 'observiq-collector' }
        end
    end
end
