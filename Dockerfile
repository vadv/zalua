FROM centos:6

RUN echo -e "[wandisco-Git]\nname=CentOS-6 - Wandisco Git\nbaseurl=http://opensource.wandisco.com/centos/6/git/\$basearch/\nenabled=1\ngpgcheck=0" > /etc/yum.repos.d/wandisco-git.repo
RUN yum install -y epel-release
RUN yum install -y wget make git gzip rpm-build nc
RUN yum groupinstall -y "Development tools"
RUN wget https://dl.google.com/go/go1.9.2.src.tar.gz -O /tmp/go.linux-amd64.tar.gz
RUN tar xvf /tmp/go.linux-amd64.tar.gz -C /usr/local
RUN rm -f /tmp/go.linux-amd64.tar.gz
RUN ln -s /usr/local/go/bin/go* /usr/local/bin/
