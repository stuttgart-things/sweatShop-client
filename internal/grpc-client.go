/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	//"google.golang.org/grpc/credentials"
	revisionrun "github.com/stuttgart-things/sweatShop-server/revisionrun"

	"google.golang.org/grpc/credentials"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Client struct {
	yasClient revisionrun.YachtApplicationServiceClient
	timeout   time.Duration
}

func NewClient(conn grpc.ClientConnInterface, timeout time.Duration) Client {
	return Client{
		yasClient: revisionrun.NewYachtApplicationServiceClient(conn),
		timeout:   timeout,
	}
}

func (c Client) CreateRevisionRun(ctx context.Context, json io.Reader) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(c.timeout))
	defer cancel()

	req := revisionrun.CreateRevisionRunRequest{}
	if err := jsonpb.Unmarshal(json, &req); err != nil {
		return fmt.Errorf("client create revisionrun: unmarshal: %w", err)
	}

	res, err := c.yasClient.CreateRevisionRun(ctx, &req)

	fmt.Println(res)

	if err != nil {
		if er, ok := status.FromError(err); ok {
			return fmt.Errorf("client create revisionrun: code: %s - msg: %s", er.Code(), er.Message())
		}
		return fmt.Errorf("client create revisionrun: %w", err)
	}

	log.Println("RESULT:", res.Result)
	log.Println("RESPONSE:", res)

	return nil
}

func ConnectSecure(address, file string) {

	log.Println("client started connecting to.. " + address)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithBlock())

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	json, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	yasClient := NewClient(conn, time.Second)
	err = yasClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(json))

	log.Println("ERR:", err)

}

func ConnectInsecure(address, file string) {

	log.Println("client started connecting to.. " + address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	json, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	yasClient := NewClient(conn, time.Second)
	err = yasClient.CreateRevisionRun(context.Background(), bytes.NewBuffer(json))

	log.Println("ERR:", err)

}
