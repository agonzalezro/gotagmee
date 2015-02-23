package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/agonzalezro/gotagmee/meetup"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s MEETUP_API_TOKEN GROUP_URL_NAME:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var neo4jDB = flag.String("neo4j", "", "If set it will store the scraped data there. Use the form: protocol://host:port/db/data")
	flag.Parse()

	api, err := meetup.NewAPI(flag.Arg(0), flag.Arg(1), *neo4jDB)
	if err != nil {
		log.Fatal(err)
	}

	membersChan := make(chan meetup.Member, 1)
	go api.Members(membersChan)
	for m := range membersChan {
		log.Println(m)
	}
}
