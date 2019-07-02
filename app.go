package main

import (
	"github.com/BonusPlay/VueHoster/routers"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	log.Println("Starting VueHoster")

	g.Go(func() error {
		log.Println("Starting main on port 3010")
		return http.ListenAndServe(":3010", routers.MainRouter())
	})

	g.Go(func() error {
		log.Println("Starting dev on port 3011")
		return http.ListenAndServe(":3011", routers.DevRouter())
	})

	g.Go(func() error {
		log.Println("Starting api on port 3020")
		return http.ListenAndServe(":3020", routers.ApiRouter())
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}