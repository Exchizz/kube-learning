package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func getLine() (string, error) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan() // use `for scanner.Scan()` to keep reading
	line := strings.TrimSpace(s.Text())
	if len(line) == 0 {
		return "", errors.New("Empty line")
	}
	return line, nil
}

func cmd_health(c cmd) {
	fmt.Printf("cmd_health called with value: %t\n", c.cmd_value)
	probes.health = c.cmd_value
}

func cmd_liveness(c cmd) {
	fmt.Printf("cmd_liveness called with value: %t\n", c.cmd_value)
	probes.liveness = c.cmd_value
}

func cmd_readyness(c cmd) {
	fmt.Printf("cmd_readyness called with value: %t\n", c.cmd_value)
	probes.readyness = c.cmd_value
}

type Callback func(cmd)

type Probes struct {
	health    bool
	liveness  bool
	readyness bool
}

var commands = map[string]Callback{
	"health":    cmd_health,
	"liveness":  cmd_liveness,
	"readyness": cmd_readyness,
}

type cmd struct {
	cmd_name  string
	cmd_value bool
}

func createCmd(line string) (cmd, error) {
	retval := cmd{}

	line_splitted := strings.Split(line, ":")
	len := len(line_splitted)
	if len != 2 {
		return cmd{}, fmt.Errorf("received %d parameters, expected 2", len)
	}
	retval.cmd_name = strings.ToLower(strings.Trim(line_splitted[0], " "))
	retval.cmd_value, _ = strconv.ParseBool(strings.Trim(line_splitted[1], " "))

	return retval, nil
}

func stdin() {
	fmt.Println("-------------------------------")
	fmt.Println("Command syntax: <cmd>:<value>")
	fmt.Println("Examples: ")
	fmt.Println("  Health: false")
	fmt.Println("  Readyness: false")
	fmt.Println("  liveness: true")
	fmt.Println("-------------------------------")

	for {
		// Get line from stdin
		line, err := getLine()
		if err != nil {
			continue
		}

		// Parse the line and convert to instance of "cmd"
		cmd, err := createCmd(line)
		if err != nil {
			fmt.Println(err)
		}

		// if cmd.cmd_name exists, run the associated funk
		if funk, ok := commands[cmd.cmd_name]; ok {
			funk(cmd)
		} else {
			fmt.Printf("Command '%s' does not exist\n", cmd.cmd_name)
		}

	}

}

var probes = Probes{}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	// initialize variables
	pod_name := getEnv("POD_NAME", "default_pod_name")
	node_name := getEnv("NODE_NAME", "local")

	// initialize probes
	probes.health = true
	probes.liveness = false
	probes.readyness = false

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("01-02-2006 15:04:05")
		fmt.Fprintf(w, "%s|%s|%s|says hello", now, node_name, pod_name)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if probes.health {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
	})
	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		if probes.liveness {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
	})
	http.HandleFunc("/readynessz", func(w http.ResponseWriter, r *http.Request) {
		if probes.readyness {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
	})

	go stdin()
	log.Printf("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", logRequest(http.DefaultServeMux)))

}
