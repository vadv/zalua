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
BuildRequires:  make

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

%install
%{__mkdir} %{buildroot}%{restream_zabbix_bin_dir}
%{__mkdir} -p %{buildroot}/%{_sysconfdir}/logrotate.d
%{__mkdir} -p %{buildroot}%{_sysconfdir}/zalua/plugins
%{__install} -m 0644 examples/config.lua %{buildroot}%{_sysconfdir}/zalua/config.lua
%{__install} -m 0644 -p %{SOURCE1} %{buildroot}/%{_sysconfdir}/logrotate.d/zalua
cp -rva examples/plugins/* %{buildroot}%{_sysconfdir}/zalua/plugins/
install bin/%{bin_name} %{buildroot}%{restream_zabbix_bin_dir}

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_zabbix_bin_dir}/%{bin_name}
%{_sysconfdir}/zalua/config.lua
%{_sysconfdir}/zalua/plugins/*
%doc README.md
