package main

import (
	"fmt"

	_ "github.com/nbisso/storicard-challenge/docs"

	"github.com/nbisso/storicard-challenge/infrastracture/conf"
	"github.com/nbisso/storicard-challenge/infrastracture/http"
)

func main() {

	s := http.NewServer()

	c, reg := s.Run(conf.Instance.Port)

	<-c

	reg.Register.CleanUp()

	fmt.Println("Shutting down gracefully")

}
