#!/bin/sh

ovs_db_socket=/var/run/openvswitch/db.sock

sysctl -w net.ipv4.ip_forward=1
sysctl -w net.ipv6.conf.all.forwarding=1

# Initialize configuration database.
ovsdb-tool create /etc/openvswitch/conf.db \
    /usr/share/openvswitch/vswitch.ovsschema

# Start the OVS database server in the detached mode.
ovsdb-server \
    --detach --pidfile \
    --remote=punix:${ovs_db_socket} \
    --remote=db:Open_vSwitch,Open_vSwitch,manager_options

# Initialized OVS database.
ovs-vsctl --db=unix:${ovs_db_socket} --no-wait init
exec "$@"
