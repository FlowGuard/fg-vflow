//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    packet.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

package packet

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

// The header protocol describes the format of the sampled header
const (
	headerProtocolEthernet   uint32 = 1
	headerProtocolTokenbus   uint32 = 2
	headerProtocolTokenring  uint32 = 3
	headerProtocolFddi       uint32 = 4
	headerProtocolFrameRelay uint32 = 5
	headerProtocolX25        uint32 = 6
	headerProtocolPpp        uint32 = 7
	headerProtocolSmds       uint32 = 8
	headerProtocolAal5       uint32 = 9
	headerProtocolAal5Ip     uint32 = 10
	headerProtocolIPv4       uint32 = 11
	headerProtocolIPv6       uint32 = 12
	headerProtocolMpls       uint32 = 13
	headerProtocolPos        uint32 = 14
)

// Packet represents layer 2,3,4 available info
type Packet struct {
	L2   Datalink
	L3   interface{}
	L4   interface{}
	data []byte
}

var (
	errUnknownEtherType      = errors.New("unknown ether type")
	errUnknownHeaderProtocol = errors.New("unknown header protocol")
)

// NewPacket constructs a packet object
func NewPacket() Packet {
	return Packet{}
}

// Decoder decodes packet's layers
func (p *Packet) Decoder(data []byte, protocol uint32) (*Packet, error) {
	var (
		err error
	)

	p.data = data

	switch protocol {
	case headerProtocolEthernet:
		err = p.decodeEthernetHeader()
		return p, err
	case headerProtocolIPv4:
		err = p.decodeIPv4Header()
		if err != nil {
			log.Error("Unable to decode IPv4 header")
			return p, err
		}
	case headerProtocolIPv6:
		err = p.decodeIPv6Header()
		if err != nil {
			log.Error("Unable to decode IPv6 header")
			return p, err
		}
	default:
		log.Errorf("Unsupported sflow protocol header %v", protocol)
		return p, errUnknownHeaderProtocol
	}

	err = p.decodeNextLayer()
	if err != nil {
		return p, err
	}

	return p, nil
}

func (p *Packet) decodeEthernetHeader() error {
	var (
		err error
	)

	err = p.decodeEthernet()
	if err != nil {
		log.Error("Unable to decode Ethernet header")
		return err
	}

	switch p.L2.EtherType {
	case EtherTypeIPv4:

		err = p.decodeIPv4Header()
		if err != nil {
			log.Error("Unable to decode IPv4 header")
			return err
		}

		err = p.decodeNextLayer()
		if err != nil {
			log.Error("Unable to decode next layer")
			return err
		}

	case EtherTypeIPv6:

		err = p.decodeIPv6Header()
		if err != nil {
			log.Error("Unable to decode IPv6 header")
			return err
		}

		err = p.decodeNextLayer()
		if err != nil {
			log.Error("Unable to decode next layer header")
			return err
		}

	default:
		log.Errorf("Unsupported sflow protocol header %v", p.L2.EtherType)
		return errUnknownEtherType
	}

	return nil
}
