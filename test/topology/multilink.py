#!/usr/bin/python

import mininet
import mininet.net
import mininet.topo
import mininet.node
import mininet.log


class Topo(mininet.topo.Topo):
    """Star topology of one switch and 3 hosts."""

    def __init__(self, **kwargs):
        super(Topo, self).__init__(**kwargs)

        s1 = self.addSwitch("s1", protocols="OpenFlow13")

        h1 = self.addHost("h1", ip="10.0.1.1/24")
        h2 = self.addHost("h2", ip="10.0.1.2/24")
        h3 = self.addHost("h3", ip="10.0.1.3/24")

        self.addLink(s1, h1, port1=1)
        self.addLink(s1, h2, port1=2)
        self.addLink(s1, h3, port1=3)


def main():
    mininet.log.setLogLevel("info")

    net = mininet.net.Mininet(topo=Topo())
    net.addController(ip="127.0.0.1",
        controller=mininet.node.RemoteController)

    set_default_route(net)
    net.start()


if __name__ == "__main__":
    main()
