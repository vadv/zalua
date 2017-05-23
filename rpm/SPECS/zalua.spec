%define version unknown
%define bin_name zalua
%define debug_package %{nil}

Name:           %{bin_name}
Version:        %{version}
Release:        1%{?dist}
Summary:        ZaLua: zabbix metric aggregator
License:        BSD
URL:            http://git.itv.restr.im/infra/%{bin_name}
Source:         %{bin_name}-%{version}.tar.gz
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
mkdir -p %{buildroot}%{restream_zabbix_bin_dir}
install -m 0755 -d %{buildroot}%{_sysconfdir}/zalua
install -m 0644 examples/config.lua %{buildroot}%{_sysconfdir}/zalua/config.lua
install -m 0755 -d %{buildroot}%{_sysconfdir}/zalua/plugins
cp examples/plugins/* %{buildroot}%{_sysconfdir}/zalua/plugins/
install bin/%{bin_name} %{buildroot}%{restream_zabbix_bin_dir}

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{restream_zabbix_bin_dir}/%{bin_name}
%{_sysconfdir}/zalua/config.lua
%{_sysconfdir}/zalua/plugins/*
%doc README.md
