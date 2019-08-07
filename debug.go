package main

import (
	"GF/cipher"
	"fmt"
	"strings"
)

func main(){
	s := cipher.XXTEAEncryptFromString("Powered by GFHelper Project.\ngfhelper.github.io | gfhelper@outlook.com")
	xx := ""
	for _, part := range s{
		xx += fmt.Sprintf("\\x%2x", part)
	}
	fmt.Println(strings.Replace(xx, " ", "0", -1))
	fmt.Println(s)
	fmt.Printf("%x", s)
}
