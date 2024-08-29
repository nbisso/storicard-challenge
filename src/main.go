package main

import (
	"fmt"

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
