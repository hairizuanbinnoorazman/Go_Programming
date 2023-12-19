FROM golang:1.18 as builder
WORKDIR /helloworld
ADD . .
RUN CGO_ENABLED=0 go build -o app .

FROM rockylinux:8 as rpm-builder
RUN dnf install -y gcc rpm-build rpm-devel rpmlint make python3.11 bash diffutils patch rpmdevtools
WORKDIR /helloworld
COPY basic.spec .
RUN rpmdev-setuptree
COPY --from=builder /helloworld/app /root/rpmbuild/BUILD/app
COPY ./deployment/bin/app.service /root/rpmbuild/BUILD/app.service
RUN rpmbuild -ba basic.spec

FROM scratch
COPY --from=rpm-builder /root/rpmbuild/RPMS /