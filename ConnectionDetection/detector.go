package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"

)

var packetChan = make(chan gopacket.Packet, 100)

func main() {
	var phyId int
	// Get a list of all available network devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Print information about each Ethernet device
	for i, device := range devices {
		for _, addr := range device.Addresses {
			if addr.IP.To4() == nil {
				// Print only Ethernet devices
				fmt.Printf("[%d] Name: %s\n", i, device.Name)
				fmt.Printf("Description: %s\n", device.Description)
				fmt.Println()
			}
		}
	}

	// Get input to select ethernet device
	fmt.Print("Enter Network PHY Identifier: ")
	_, err = fmt.Scan(&phyId)
	if err != nil {
		fmt.Println("Invalid input. Please enter a valid number.")
		return
	}

	// Specify the network interface to capture packets from
	device := devices[phyId].Name
	fmt.Println(devices[phyId].Description, "selected.")

	// Open the device for packet capture
	handle, err := pcap.OpenLive(device, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set a filter to capture only broadcast packets
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
	}

	// Start packet capture in a separate goroutine
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	go capturePackets(packetSource)
	go treatPacketChan()
	// Wait for termination signal (Ctrl+C)
	waitForTerminationSignal()
}

func capturePackets(packetSource *gopacket.PacketSource) {

	// Start the packet capture loop
	for packet := range packetSource.Packets() {
		// Increment the packet count
		packetChan <- packet
	}

}

func treatPacketChan(){
	for packet := range packetChan{
		packet1 := packet.Layer(layers.LayerTypeIPv4)
		if packet1 != nil{
			fmt.Printf("packet: %v %v\n", packet1.LayerContents(),len(packet1.LayerContents()))
			
		}else{
			fmt.Println("tuputamadre")
		}
		
	}
}

func waitForTerminationSignal() {
	// Create a channel to receive termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for the termination signal
	<-sigChan
	fmt.Println("\nTermination signal received. Exiting...")
	os.Exit(0)
}
