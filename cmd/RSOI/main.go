package main

import (
	"flag"
	"fmt"
)

func main() {
	service := flag.String("service", "", "Service name. Either post, poststats, comment or api")
	connString := flag.String("conn", "", "PostgreSQL connection string")
	portNum := flag.Int("port", -1, "Port on which service well listen")
	postServerAddr := flag.String("post-server", "", "Address of post server")
	commentServerAddr := flag.String("comment-server", "", "Address of comment server")
	postStatsServerAddr := flag.String("post-stats-server", "", "Address of post stats server")
	jaegerAddr := flag.String("jaeger-addr", "", "Jaeger address")

	flag.Parse()

	port := *portNum
	conn := *connString
	ps := *postServerAddr
	cs := *commentServerAddr
	pss := *postStatsServerAddr
	ja := *jaegerAddr

	var err error
	switch *service {
	case "post":
		fmt.Printf("running post service on port %d\n", port)
		err = runPost(port, conn, ja)
	case "comment":
		fmt.Printf("running comment service on port %d\n", port)
		err = runComment(port, conn, ja)
	case "poststats":
		fmt.Printf("running post stats service on port %d\n", port)
		err = runPostStats(port, conn, ja)
	case "api":
		fmt.Printf("running API service on port %d\n", port)
		err = runAPI(port, ps, cs, pss, ja)
	default:
		fmt.Printf("unknown service %v\n", service)
	}

	if err != nil {
		fmt.Printf("finished with error %v", err)
	}
}
