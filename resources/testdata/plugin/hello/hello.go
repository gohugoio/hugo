package main

import "fmt"

const Language = "english"

var HelloFmt = map[string]string{
	"english": "Hello %s",
	"french":  "Salutation %s",
	"spanish": "Hola %s",
}

func Hello(name string) string { return fmt.Sprintf(HelloFmt[Language], name) }
