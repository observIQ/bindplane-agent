describe service('observiq-otel-collector') do
    it { should_not be_enabled }
    it { should_not be_installed }
    it { should_not be_running }
end
