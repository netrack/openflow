package main

import (
	"bytes"
	"log"

	"github.com/netrack/net/pkg"
	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp.v13"
)

func main() {
	var id uint16

	of.HandleFunc(of.T_HELLO, func(rw of.ResponseWriter, r *of.Request) {
		log.Println("RECV ofp_hello:", r.Header)
		rw.Header().Set(of.TypeHeaderKey, of.T_HELLO)
		rw.Header().Set(of.VersionHeaderKey, ofp.VERSION)
		rw.WriteHeader()

		log.Println("SEND ofp_hello:", rw.Header())
	})

	of.HandleFunc(of.T_PACKET_IN, func(rw of.ResponseWriter, r *of.Request) {
		var p ofp.PacketIn
		var eth pkg.EthernetII

		plen, err := p.ReadFrom(r.Body)
		log.Println("#1", err)

		ethlen, err := eth.ReadFrom(r.Body)
		log.Println("#2", err)

		log.Println("RECV ofp_packet_in:", r.Header)
		log.Println("RECV ofp_packet_in:", eth.HWDst, eth.HWSrc)

		pout := ofp.PacketOut{
			BufferID: ofp.NO_BUFFER,
			InPort:   1,
			Actions:  ofp.Actions{ofp.ActionOutput{ofp.P_IN_PORT, 0}},
		}

		rw.Header().Set(of.TypeHeaderKey, of.T_PACKET_OUT)
		rw.Header().Set(of.VersionHeaderKey, ofp.VERSION)

		_, err = rw.Write(pout.Bytes())
		log.Println("#3", err)

		if eth.EtherType == pkg.PROTO_ARP {
			var arp pkg.ARP

			_, err = arp.ReadFrom(r.Body)
			log.Println("#4", err)
			log.Println("RECV ARP:", arp)

			eth := pkg.EthernetII{
				HWDst:     eth.HWSrc,
				HWSrc:     []byte{0, 0, 0, 0, 0, 254},
				EtherType: pkg.PROTO_ARP,
			}

			_, err = eth.WriteTo(rw)
			log.Println("#5.1", err)

			arp = pkg.ARP{
				HWType:    pkg.ARPT_ETHERNET,
				ProtoType: pkg.PROTO_IPV4,
				Operation: pkg.ARPOT_REPLY,
				HWSrc:     []byte{0, 0, 0, 0, 0, 254},
				ProtoSrc:  arp.ProtoDst,
				HWDst:     arp.HWSrc,
				ProtoDst:  arp.ProtoSrc,
			}

			_, err = arp.WriteTo(rw)
			log.Println("#5.2", err)
		}

		if eth.EtherType == pkg.PROTO_IPV4 {
			var ip pkg.IPv4
			var icmp pkg.ICMP
			var echo pkg.ICMPEcho

			iplen, err := ip.ReadFrom(r.Body)
			log.Println("#6", err)
			log.Println("RECV IPv4:", ip)

			icmplen, err := icmp.ReadFrom(r.Body)
			log.Println("#6.1", icmp)

			echo.Data = make([]byte, r.ContentLength-plen-ethlen-iplen-icmplen-4)
			echlen, err := echo.ReadFrom(r.Body)
			log.Println("#6.2", err)

			log.Println("LENGTH:", r.ContentLength, plen, ethlen, iplen, icmplen, echlen)

			eth = pkg.EthernetII{
				HWDst:     eth.HWSrc,
				HWSrc:     eth.HWDst,
				EtherType: pkg.PROTO_IPV4,
			}

			_, err = eth.WriteTo(rw)
			log.Println("#6.2.1", err)

			ip = pkg.IPv4{
				TTL:   255,
				ID:    id,
				Len:   uint16(iplen + icmplen + echlen),
				Proto: pkg.IPV4_PROTO_ICMP,
				Src:   []byte{10, 0, 0, 2},
				Dst:   []byte{10, 0, 0, 1},
			}

			id++
			_, err = ip.WriteTo(rw)
			log.Println("#6.3", err)

			var buf bytes.Buffer
			icmp = pkg.ICMP{Type: pkg.ICMPT_ECHO_REPLY}

			_, err = icmp.WriteTo(&buf)
			log.Println("#6.4", err)

			_, err = echo.WriteTo(&buf)
			log.Println("#6.5", err)

			b := buf.Bytes()
			copy(b[2:4], pkg.Checksum(b))

			rw.Write(b)
		}

		err = rw.WriteHeader()
		log.Println("#5", err)

		log.Println("SEND ofp_of_mod:", rw.Header())
	})

	of.HandleFunc(of.T_ECHO_REQUEST, func(rw of.ResponseWriter, r *of.Request) {
		log.Println("RECV ofp_echo_request:", r.Header)
		rw.Header().Set(of.TypeHeaderKey, of.T_ECHO_REPLY)
		rw.Header().Set(of.VersionHeaderKey, ofp.VERSION)
		rw.WriteHeader()
		log.Println("SEND ofp_echo_reply:", rw.Header())
	})

	log.Println("started listening...")
	of.ListenAndServe()
}
