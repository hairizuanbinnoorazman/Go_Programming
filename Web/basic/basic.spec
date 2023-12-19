Name:       basic
Version:    0
Release:    1
Summary:    RPM package to contain basic Golang app
License:    FIXME

%description
RPM package to encapsulate basic golang application

%prep
# we have no source, so nothing here

%build
# Built using Golang docker image

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}/etc/systemd/system/
install -m 755 app %{buildroot}%{_bindir}/app
install -m 755 app.service %{buildroot}/etc/systemd/system/basic.service

%files
%{_bindir}/app
/etc/systemd/system/basic.service

%pre
getent group app >/dev/null 2>&1 || groupadd app
getent passwd app >/dev/null 2>&1 || useradd -G app app

%post
chown app:app %{_bindir}/app
systemctl daemon-reload
systemctl enable basic.service
systemctl start basic.service

%preun
systemctl stop basic.service
systemctl disable basic.service
systemctl daemon-reload
# systemctl reset-failed - not sure if needed here

%postun
userdel app
groupdel app

%changelog
# let's skip this for now