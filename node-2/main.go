package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// The multiaddress of Node 1
	if len(os.Args) < 2 {
		fmt.Println("Please provide the multiaddress of the peer to connect to")
		fmt.Println("Usage: go run main.go <multiaddress>")
		os.Exit(1)
	}

	peerAddr := os.Args[1]

	// Turn the multiaddress into a peer ID and target address
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		panic(err)
	}

	peerinfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	// Connect to the peer
	fmt.Println("Connecting to:", peerinfo.ID)
	if err := node.Connect(context.Background(), *peerinfo); err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	fmt.Println("Connected to:", peerinfo.ID)

	// Open a stream with the peer
	s, err := node.NewStream(context.Background(), peerinfo.ID, "/chat/1.0.0")
	if err != nil {
		fmt.Println("Stream open failed:", err)
		return
	}
	fmt.Println("Stream opened with:", peerinfo.ID)

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go writeData(rw)
	go readData(rw)

	select {}
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer:", err)
			return
		}
		if str != "" {
			fmt.Printf("Received message: %s", str)
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}

		_, err = rw.WriteString(sendData)
		if err != nil {
			fmt.Println("Error writing to buffer:", err)
			return
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer:", err)
			return
		}
	}
}
