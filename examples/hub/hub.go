package main

import (
	"log"

	"github.com/netrack/net/pkg"
	flow "github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp.v13"
)

func main() {
	flow.HandleFunc(flow.T_HELLO, func(rw flow.ResponseWriter, r *flow.Request) {
		log.Println("RECV ofp_hello:", r.Header)
		rw.Header().Set(flow.TypeHeaderKey, flow.T_HELLO)
		rw.Header().Set(flow.VersionHeaderKey, ofp.VERSION)
		rw.WriteHeader()

		log.Println("SEND ofp_hello:", rw.Header())
	})

	flow.HandleFunc(flow.T_PACKET_IN, func(rw flow.ResponseWriter, r *flow.Request) {
		var p ofp.PacketIn
		var eth pkg.EthernetII

		p.ReadFrom(r.Body)
		eth.Read(r.Body)

		log.Println("RECV ofp_packet_in:", r.Header)
		log.Println("RECV ofp_packet_in:", eth.HWDst, eth.HWSrc)

		instr := ofp.Instructions{ofp.InstructionActions{
			ofp.IT_APPLY_ACTIONS,
			ofp.Actions{ofp.ActionOutput{ofp.P_FLOOD, 0}},
		}}

		fmod := &ofp.FlowMod{
			Command:      ofp.FC_ADD,
			BufferID:     p.BufferID,
			Match:        p.Match,
			Instructions: instr,
		}

		rw.Header().Set(flow.TypeHeaderKey, flow.T_FLOW_MOD)
		rw.Header().Set(flow.VersionHeaderKey, ofp.VERSION)

		rw.Write(fmod.Bytes())
		log.Println("SEND ofp_flow_mod:", rw.Header())
	})

	flow.HandleFunc(flow.T_ECHO_REQUEST, func(rw flow.ResponseWriter, r *flow.Request) {
		log.Println("RECV ofp_echo_request:", r.Header)
		rw.Header().Set(flow.TypeHeaderKey, flow.T_ECHO_REPLY)
		rw.Header().Set(flow.VersionHeaderKey, ofp.VERSION)
		rw.WriteHeader()
		log.Println("SEND ofp_echo_reply:", rw.Header())
	})

	log.Println("started listening...")
	flow.ListenAndServe()
}
