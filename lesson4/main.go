package main

import (
	"fmt"
)

type person struct {
	firstName string
	lastName  string
}

func main() {
	personA := person{
		firstName: "LL",
		lastName:  "00",
	}
	(&personA).changePerson();
	fmt.Println(personA)

	// fmt.Println(personA)

}

func (p *person )changePerson( )  {
	p.firstName="LLLL";
}
