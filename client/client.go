package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/google/uuid"
	"io/ioutil"
	"jtunnel-go/proto"
	"net"
	"os"
	"strconv"
)

var registered = false

func main() {

	parser := argparse.NewParser("jtunnel", "Expose local server over the internet")
	subdomain := parser.String("s", "subdomain", &argparse.Options{Required: true, Help: "Subdomain eg. test.jtunnel.net"})
	localport := parser.Int("p", "localport", &argparse.Options{Required: true, Help: "Local server port to expose"})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}
	conn, err := createControlConnection(*subdomain)
	if err != nil {
		fmt.Println("Unable to connect to manoj.jtunnel.net:8585", err.Error())
	}
	handleData(conn, *localport, *subdomain)

}

func makeRequest(data []byte, localPort int) ([]byte, error) {
	fmt.Println(string(data))
	server, _ := net.ResolveTCPAddr("tcp", "localhost:"+strconv.Itoa(localPort))
	client, _ := net.ResolveTCPAddr("tcp", ":")
	conn, err := net.DialTCP("tcp", client, server)
	if err != nil {
		fmt.Println("Error is ", err.Error())
		return nil, err
	} else {
		fmt.Println("Making Local Request")
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Received Error", err.Error())
			return nil, err
		}
		fmt.Println("Reading Local Response")
		result, _ := ioutil.ReadAll(conn)
		fmt.Println("Read Local Response")
		//fmt.Println(string(result))
		return result, nil
	}
}

func createControlConnection(subdomain string) (*net.TCPConn, error) {
	server, _ := net.ResolveTCPAddr("tcp", "localhost:9999" /*subdomain+".jtunnel.net:8585"*/)
	client, _ := net.ResolveTCPAddr("tcp", ":")
	conn, err := net.DialTCP("tcp", client, server)
	return conn, err
}

func handleData(conn *net.TCPConn, localPort int, subdomain string) {
	for {
		if !registered {
			err := proto.SendMessage(proto.NewMessage("localhost:8080" /*subdomain+".jtunnel.net"*/, uuid.New().String(), "register", make([]byte, 0)), conn)
			if err != nil {
				fmt.Println("Unable to handleData ", err.Error())
			}
			registered = true
		}
		message, err := proto.ReceiveMessage(conn)
		if err != nil {
			fmt.Println("Unable to receive message", err.Error())
		}
		resp, err := makeRequest(message.Data, localPort)
		if err != nil {
			fmt.Println("Unable to make local request", err.Error())
			return
		}
		message = proto.NewMessage("localhost:8080" /*subdomain+".jtunnel.net"*/, message.MessageId, "response", resp)
		err = proto.SendMessage(message, conn)
		if err != nil {
			fmt.Println("Unable to send response to tunnel server ", err.Error())
		}

	}

}
