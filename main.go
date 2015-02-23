package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/agonzalezro/gotagmee/db"
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

	if *neo4jDB == "" {
		fmt.Fprintf(os.Stderr, "I know that `-neo4j` is an option, but it is not. The code without Neo4J is not implemented, sorry! :(")
	}

	api := meetup.NewAPI(flag.Arg(0), flag.Arg(1))

	membersChan := make(chan db.Member, 1)
	go api.Members(membersChan)

	db, err := db.NewDB(*neo4jDB)
	if err != nil {
		log.Fatal(err)
	}

	for m := range membersChan {
		db.Store(m)
	}
}
