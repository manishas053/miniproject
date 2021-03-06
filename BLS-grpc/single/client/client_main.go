package main

import (
	"log"

	pb "google.golang.org/grpc/examples/bls/bls"
	"github.com/Nik-U/pbc"
	"crypto/sha256"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewSignServiceClient(conn)
	params := pbc.GenerateA(160, 512)
  pairing := params.NewPairing()
 	g := pairing.NewG2().Rand()

	r, err := c.SignString(context.Background(), &pb.SignRequest{

		SharedParams : params.String(),
  	SharedG : g.Bytes(),

	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("Signature: ", r.Signature)
	h := pairing.NewG1().SetFromStringHash(r.Data, sha256.New())
	signature := pairing.NewG1().SetBytes(r.Signature)
	public := pairing.NewG2().SetBytes(r.Publickey)
	temp1 := pairing.NewGT().Pair(h, public)
	temp2 := pairing.NewGT().Pair(signature, g)

	log.Println("")
	log.Println("e(v, h): ", temp1)
	log.Println("")
	log.Println("e(g, σ): ", temp2)

	if !temp1.Equals(temp2) {
      		log.Println("*BUG* Signature check failed *BUG*")
  } else {
      		log.Println("Signature verified correctly")
  }
}
