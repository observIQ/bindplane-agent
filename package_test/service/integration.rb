describe systemd_service('observiq-collector') do
    it { should be_installed }
    it { should be_enabled }
    it { should be_running }

    describe file('/usr/lib/systemd/system/observiq-collector.service') do
        its('mode') { should cmp '0644' }
        its('owner') { should eq 'root' }
        its('group') { should eq 'root' }
        its('type') { should cmp 'file' }
    end
end
