package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prompt.Debugf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func getLine() (string, error) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan() // use `for scanner.Scan()` to keep reading
	line := strings.TrimSpace(s.Text())
	if len(line) == 0 {
		return "", errors.New("empty line")
	}
	return line, nil
}

func cmd_liveness(c cmd) {
	prompt.Printf("cmd_liveness called with value: %t", c.cmd_value)
	probes.liveness = c.cmd_value
}

func cmd_readiness(c cmd) {
	prompt.Printf("cmd_readiness called with value: %t", c.cmd_value)
	probes.readiness = c.cmd_value
}

func cmd_verbose(c cmd) {
	if c.cmd_value {
		log.Info("Enabling debug log")
		log.SetLevel(log.DebugLevel)
	} else {
		log.Info("Disabling debug log")
		log.SetLevel(log.WarnLevel)
	}
}

type Callback func(cmd)

type Probes struct {
	liveness  bool
	readiness bool
}

var commands = map[string]Callback{
	"liveness":  cmd_liveness,
	"readiness": cmd_readiness,
	"verbose":   cmd_verbose,
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

func show_status() {
	fmt.Printf("Status: \n")
	fmt.Printf("  Liveness: %t\n", probes.liveness)
	fmt.Printf("  readiness: %t\n", probes.readiness)
	fmt.Printf("  Node name: %s\n", node_name)
	fmt.Printf("  Pod name: %s\n", pod_name)

}
func stdin() {
	log.Printf("-------------------------------")
	log.Printf("Command syntax: <cmd>:<value>")
	log.Printf("Examples: ")
	log.Printf("  readiness: false")
	log.Printf("  liveness: true")
	log.Printf("Press \"?\" to see the current state of the probes")
	log.Printf("-------------------------------")

	for {
		// Get line from stdin
		line, err := getLine()
		if err != nil {
			continue
		}

		if line == "?" {
			show_status()
			continue
		}
		// Parse the line and convert to instance of "cmd"
		cmd, err := createCmd(line)
		if err != nil {
			log.Warn(err)
		}

		// if cmd.cmd_name exists, run the associated funk
		if funk, ok := commands[cmd.cmd_name]; ok {
			funk(cmd)
		} else {
			log.Warn("Command '%s' does not exist", cmd.cmd_name)
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

// Gets the env variable as a boolean.
// Handles "1", "t", "T", "TRUE", "true", "True" as truthy, everything else (including env variable being empty) as falsy
// If the env variable isn't set at all, uses the given fallback
func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		return parsed && err == nil
	}
	return fallback
}

var prompt *log.Entry
var pod_name string
var node_name string

func main() {

	// initialize variables
	pod_name = getEnv("POD_NAME", "default_pod_name")
	node_name = getEnv("NODE_NAME", "local")
	prompt = log.WithFields(log.Fields{
		"node_name": node_name,
		"pod_name":  pod_name,
	})

	// initialize probes
	probes.liveness = getEnvBool("KUBELEARN_ALIVE", true)
	probes.readiness = getEnvBool("KUBELEARN_READY", true)
	if getEnvBool("KUBELEARN_DEBUG", false) {
		cmd_verbose(cmd{
			cmd_name:  "verbose",
			cmd_value: true,
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("01-02-2006 15:04:05")
		fmt.Fprintf(w, "%s|%s|%s|says hello", now, node_name, pod_name)
	})

	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		if probes.liveness {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
	})
	http.HandleFunc("/readinessz", func(w http.ResponseWriter, r *http.Request) {
		if probes.readiness {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
	})

	go stdin()
	log.Info("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", logRequest(http.DefaultServeMux)))

}
