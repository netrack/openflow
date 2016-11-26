FROM debian:sid

# Install Open vSwtich and mininet to create custom topologies.
RUN apt-get update \
    && apt-get install -y \
        openvswitch-switch \
        net-tools \
        mininet \
    && apt-get clean

# Create a run directory for Open vSwitch daemons.
RUN mkdir -p /var/run/openvswitch /usr/share/mininet

# Copy the entrypoint, which starts the OVS database server.
COPY binary/docker-entrypoint.sh /usr/bin/docker-entrypoint
COPY topology/* /usr/share/mininet/

ENTRYPOINT ["docker-entrypoint"]
CMD ["ovs-vswitchd", "--pidfile"]
