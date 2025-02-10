package main

import (
	"context"
	"fmt"
	"log"
	"templates_new/internal/app"
)

type SomeInterface interface {
	SomeMethod(in string) (out string)
}

type SomeStruct struct {
	prifix string
}

func (s SomeStruct) SomeMethod(in string) (out string) {
	return s.prifix + ": " + in
}

func main() {
	var iv1 SomeInterface = SomeStruct{"I say"}

	fmt.Println(iv1 == nil)
	var f any
	fmt.Println(f == nil)

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
