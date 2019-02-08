
FROM fedora:27

ARG CYBERPROBE_VERSION

RUN dnf install -y libgo

COPY fedora-cyberprobe-${CYBERPROBE_VERSION}-1.fc27.x86_64.rpm /cyberprobe.rpm
RUN dnf install -y /cyberprobe.rpm
RUN rm -f /fedora-cyberprobe-${CYBERPROBE_VERSION}-1.fc26.x86_64.rpm
COPY cyberprobe.cfg /etc/cyberprobe.cfg

COPY socks-proxy /usr/local/bin/
COPY start /

CMD /start

EXPOSE 8080

