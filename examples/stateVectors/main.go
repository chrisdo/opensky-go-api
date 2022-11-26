package main

import (
	"log"
	"time"

	"github.com/chrisdo/opensky-go-api"
)

func main() {

	client := opensky.NewClient()

	req := opensky.NewStateVectorRequest().
		IncludeCategory()
	//	WithBoundingBox(opensky.NewBoundingBox(49.96708, 50.252014, 19.5057256, 20.3726393))

	vectors, err := client.RequestStateVectors(req)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d vectors received at: %v\n", len(vectors.States), time.Unix(vectors.Time, 0))

	for i, v := range vectors.States {
		if i > 100 { //lets skip, just for demonstration purpose
			break
		}
		log.Printf("%v\n", v)
	}

}
