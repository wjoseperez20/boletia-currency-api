package main

import (
	"fmt"
	"github.com/wjoseperez20/boletia-currency-api/pkg/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
