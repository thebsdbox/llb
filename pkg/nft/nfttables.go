package nft

import (
	"fmt"
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
)

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}

// CreateLoadBalancer -
func CreateLoadBalancer() {
	c := &nftables.Conn{}

	// Remove existing rules
	c.FlushRuleset()

	// Create nat nftable
	nat := c.AddTable(&nftables.Table{
		Family: nftables.TableFamilyIPv4,
		Name:   "nat",
	})

	// Add prerouting chain
	prerouting := c.AddChain(&nftables.Chain{
		Name:     "prerouting",
		Hooknum:  nftables.ChainHookPrerouting,
		Priority: nftables.ChainPriorityFilter,
		Table:    nat,
		Type:     nftables.ChainTypeNAT,
	})

	// Add postrouting chain
	postrouting := c.AddChain(&nftables.Chain{
		Name:     "postrouting",
		Hooknum:  nftables.ChainHookPostrouting,
		Priority: nftables.ChainPriorityNATSource,
		Table:    nat,
		Type:     nftables.ChainTypeNAT,
	})

	// Add rule for masquerading of return packets
	c.AddRule(&nftables.Rule{
		Table: nat,
		Chain: postrouting,
		Exprs: []expr.Any{
			&expr.Masq{},
		},
	})

	// Create VIP 192.168.1.5 that load balances to 192.168.1.110

	c.AddRule(&nftables.Rule{
		Table: nat,
		Chain: prerouting,
		Exprs: []expr.Any{
			// [ payload load 4b @ network header + 16 => reg 1 ]
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       16, // TODO
				Len:          4,  // TODO
			},

			//  [ cmp eq reg 1 0x0501a8c0 ]

			&expr.Cmp{
				Op:       expr.CmpOpEq,
				Register: 1,
				Data:     net.ParseIP("192.168.1.5").To4(),
			},
			//  [ meta load l4proto => reg 1 ]

			&expr.Meta{Key: expr.MetaKeyL4PROTO, Register: 1},
			//  [ cmp eq reg 1 0x00000006 ]

			&expr.Cmp{
				Op:       expr.CmpOpEq,
				Register: 1,
				Data:     []byte{unix.IPPROTO_TCP},
			},
			//  [ payload load 2b @ transport header + 2 => reg 1 ]

			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseTransportHeader,
				Offset:       2, // TODO
				Len:          2, // TODO
			},

			// [ cmp eq reg 1 0x0000e60f ]
			&expr.Cmp{
				Op:       expr.CmpOpEq,
				Register: 1,
				Data:     binaryutil.BigEndian.PutUint16(6443),
			},

			// [ immediate reg 1 0x0217a8c0 ]
			&expr.Immediate{
				Register: 1,
				Data:     net.ParseIP("192.168.1.110").To4(),
			},
			// [ immediate reg 2 0x0000f00f ]
			&expr.Immediate{
				Register: 2,
				Data:     binaryutil.BigEndian.PutUint16(6443),
			},

			// [ nat dnat ip addr_min reg 1 addr_max reg 0 proto_min reg 2 proto_max reg 0 ]
			&expr.NAT{
				Type:        expr.NATTypeDestNAT,
				Family:      unix.NFPROTO_IPV4,
				RegAddrMin:  1,
				RegProtoMin: 2,
			},
		},
	})

	// Flush all rules/chains and tables
	if e := c.Flush(); e != nil {
		fmt.Printf("%v\n", e)
	}

	// Print Tables debug
	t, e := c.ListTables()
	if e != nil {
		fmt.Printf("%v\n", e)
	}
	for x := range t {
		fmt.Printf("%v\n", t[x].Name)
	}
}
