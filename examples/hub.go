package main

import (
	"log"

	"github.com/netrack/net/l2"
	"github.com/netrack/openflow"
	"github.com/netrack/openflow/ofp.v13"
)

func main() {
	of.HandleFunc(of.T_HELLO, func(rw of.ResponseWriter, r *of.Request) {
		log.Println("RECV ofp_hello:", r.Header)
		rw.Header().Set(of.TypeHeaderKey, of.T_HELLO)
		rw.Header().Set(of.VersionHeaderKey, ofp.VERSION)
		rw.WriteHeader()

		log.Println("SEND ofp_hello:", rw.Header())
	})

	of.HandleFunc(of.T_PACKET_IN, func(rw of.ResponseWriter, r *of.Request) {
		var p ofp.PacketIn
		var e l2.EthernetII

		p.ReadFrom(r.Body)
		e.ReadFrom(r.Body)

		log.Println("RECV ofp_packet_in:", r.Header)
		log.Println("RECV ofp_packet_in:", e.HWDst, e.HWSrc)

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

		rw.Header().Set(of.TypeHeaderKey, of.T_FLOW_MOD)
		rw.Header().Set(of.VersionHeaderKey, ofp.VERSION)

		rw.Write(fmod.Bytes())
		rw.WriteHeader()
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
