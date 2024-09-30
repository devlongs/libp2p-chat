package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/2000"),
	)
	if err != nil {
		panic(err)
	}

	// Set a stream handler on the host
	node.SetStreamHandler("/chat/1.0.0", handleStream)

	// Get the full multiaddresses with peer ID
	peerInfo := &peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	addrs, err := peer.AddrInfoToP2pAddrs(peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Node 1 is listening on:")
	for _, addr := range addrs {
		fmt.Println(addr)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	if err := node.Close(); err != nil {
		panic(err)
	}
}

// handleStream is a stream handler function
func handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)

	go writeData(rw)
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
