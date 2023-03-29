package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/grpc-storage/pkg/api/balancer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
)

type commands struct {
	cl balancer.BalancerClient
}
type action string

var ErrWrongInput = errors.New("wrong input")
var ErrUnknownCMD = errors.New("unknown cmd")

const (
	get = "get"
	set = "set"
)

func (c *commands) do(act action, args ...string) error {
	switch act {
	case get:
		if len(args) < 1 {
			return ErrWrongInput
		}
		return c.get(args[0])
	case set:
		if len(args) < 2 {
			return ErrWrongInput
		}
		return c.set(args[0], strings.Join(args[1:], " "))
	}

	return ErrUnknownCMD
}

func (c *commands) get(key string) error {
	res, err := c.cl.Get(context.Background(), &balancer.BalancerGetRequest{Key: key})
	if err != nil {
		return err
	}

	log.Println(res.Value)
	return nil
}

func (c *commands) set(key, value string) error {
	res, err := c.cl.Set(context.Background(), &balancer.BalancerSetRequest{Key: key, Value: value})
	if err != nil {
		return err
	}

	log.Println(res.String())
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:800", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var cmds commands

	cmds.cl = balancer.NewBalancerClient(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		scanner.Scan()
		line := scanner.Text()
		split := strings.Split(line, " ")

		err = cmds.do(action(strings.ToLower(split[0])), split[1:]...)
		if err != nil {
			log.Println("Error! :", err)
		}
	}
}
