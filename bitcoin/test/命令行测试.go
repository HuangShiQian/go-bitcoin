package main

import (
	"os"
	"fmt"
)

func main()  {
	//var Args []string
	input:=os.Args
	for i,v:=range input{
		fmt.Printf("i= %d,v= %s\n",i,v)

	}
}
