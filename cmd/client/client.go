package main

import (
	"bufio"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	flag.Parse()

	content, _ := ioutil.ReadFile("port")
	text := string(content)
	address := "localhost:" + text

	var (
		addr = flag.String("addr", address, "the address to connect to")
	)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSumClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//------------------------------------------------------------------------------------------------------------------

	writeFile, err := os.OpenFile("./bin/output", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file to write to")
	}

	fileWriter := bufio.NewWriter(writeFile)

	file, err := os.Open("./bin/input")
	if err != nil {
		log.Fatalf("could not open file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	//for each line do this
	for scanner.Scan() {
		input := strings.Fields(scanner.Text())
		num1, _ := strconv.ParseInt(input[1], 10, 64)
		num2, _ := strconv.ParseInt(input[2], 10, 64)

		if input[0] == "add" {
			r, err := c.Addition(ctx, &pb.CalculationRequest{Int1: num1, Int2: num2})
			if err != nil {
				log.Fatalf("could not get sum: %v", err)
			}
			_, _ = fileWriter.WriteString(strconv.Itoa(int(r.GetResult())) + "\n")
			fileWriter.Flush()

		} else if input[0] == "sub" {
			r, err := c.Subtraction(ctx, &pb.CalculationRequest{Int1: num1, Int2: num2})
			if err != nil {
				log.Fatalf("could not get difference: %v", err)
			}
			_, _ = fileWriter.WriteString(strconv.Itoa(int(r.GetResult())) + "\n")
			fileWriter.Flush()

		} else if input[0] == "mult" {
			r, err := c.Multiplication(ctx, &pb.CalculationRequest{Int1: num1, Int2: num2})
			if err != nil {
				log.Fatalf("could not get product: %v", err)
			}
			_, _ = fileWriter.WriteString(strconv.Itoa(int(r.GetResult())) + "\n")
			fileWriter.Flush()
		} else if input[0] == "div" {
			if num2 != 0 {
				r, err := c.Division(ctx, &pb.CalculationRequest{Int1: num1, Int2: num2})
				if err != nil {
					log.Fatalf("could not get division: %v", err)
				}
				_, _ = fileWriter.WriteString(strconv.Itoa(int(r.GetResult())) + "\n")
				fileWriter.Flush()
			} else {
				_, _ = fileWriter.WriteString("Division by zero: ERROR" + "\n")
				fileWriter.Flush()
			}
		}
	}

	writeFile.Close()
	//if scanner error do this
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
