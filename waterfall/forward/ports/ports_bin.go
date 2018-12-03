// ports_bin starts a gRPC server that manages port forwarding sessions in the host.
package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"github.com/waterfall/forward/ports"
	waterfall_grpc "github.com/waterfall/proto/waterfall_go_grpc"
	"google.golang.org/grpc"
)

var (
	// For qemu connections addr is the working dir of the emulator
	addr = flag.String("addr", "", "Address to listen for port forwarding requests. <unix|tcp>:addr")
	waterfallAddr = flag.String("waterfall_addr", "", "Address of the waterfall server. <unix|tcp>:addr")
)

func init() {
	flag.Parse()
}

func main() {
	log.Println("Starting port forwarding server ...")

	if *addr == "" || *waterfallAddr == "" {
		log.Fatalf("Need to specify -addr and -waterfall_addr.")
	}

	pts := strings.SplitN(*addr, ":", 2)
	if len(pts) != 2 {
		log.Fatalf("failed to parse address %s", addr)
	}

	lis, err := net.Listen(pts[0], pts[1])
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	conn, err := grpc.Dial(*waterfallAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to establish connection to waterfall server: %v", err)
	}
	defer conn.Close()

	grpcServer := grpc.NewServer()
	waterfall_grpc.RegisterPortForwarderServer(grpcServer, ports.NewServer(waterfall_grpc.NewWaterfallClient(conn)))

	log.Println("Forwarding ports ...")
	grpcServer.Serve(lis)
}
