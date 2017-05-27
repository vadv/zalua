%define version unknown
%define bin_name zalua
%define debug_package %{nil}

Name:           %{bin_name}
Version:        %{version}
Release:        1%{?dist}
Summary:        ZaLua: zabbix metric aggregator with plugin in lua
License:        BSD
URL:            http://git.itv.restr.im/infra/%{bin_name}
Source:         %{bin_name}-%{version}.tar.gz
Source1:        zalua-logrotate.in
Requires:       zabbix-agent

%define restream_dir /opt/restream/
%define restream_zabbix_bin_dir %{restream_dir}/zabbix/bin

%description
This package provides zabbix monitoring with lua.

%prep
%setup

%build
make

%post
rm -f /tmp/%{bin_name}-mon.sock
%post
if [ -x /usr/sbin/setcap ]; then
    setcap cap_net_raw+ep %{restream_zabbix_bin_dir}/%{bin_name}
else
    chmod 1777 %{restream_zabbix_bin_dir}/%{bin_name}
fi

%install
# bin
%{__mkdir} -p %{buildroot}%{restream_zabbix_bin_dir}
install -m 0755 bin/%{bin_name} %{buildroot}%{restream_zabbix_bin_dir}
# logrotate
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/logrotate.d
%{__install} -m 0644 %{SOURCE1} %{buildroot}/%{_sysconfdir}/logrotate.d/%{bin_name}
# plugins
%{__mkdir} -p %{buildroot}%{_sysconfdir}/%{bin_name}/plugins
cp -v config/plugins/*.lua %{buildroot}%{_sysconfdir}/%{bin_name}/plugins/
%{__install} -m 0644 config/init.lua %{buildroot}%{_sysconfdir}/%{bin_name}/init.lua
# zabbix
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/zabbix/zabbix.d
%{__install} -m 0644 config/zabbix.conf %{buildroot}/%{_sysconfdir}/zabbix/zabbix.d/%{bin_name}_default.conf

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_zabbix_bin_dir}/%{bin_name}
%{_sysconfdir}/%{bin_name}/init.lua
%{_sysconfdir}/%{bin_name}/plugins/*
%{_sysconfdir}/logrotate.d/%{bin_name}
%{_sysconfdir}/zabbix/zabbix.d/%{bin_name}_default.conf
%doc README.md
