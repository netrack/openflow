FROM debian:sid

# Install Open vSwtich and mininet to create custom topologies.
RUN apt-get update \
    && apt-get install -y \
        openvswitch-switch \
        net-tools \
        mininet \
    && apt-get clean

# Create a run directory for Open vSwitch daemons.
RUN mkdir -p /var/run/openvswitch
COPY binaries/docker-entrypoint.sh /usr/bin/docker-entrypoint

ENTRYPOINT ["docker-entrypoint"]
CMD ["ovs-vswitchd", "--pidfile"]
