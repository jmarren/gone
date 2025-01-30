package main

import (
	"fmt"

	"github.com/jmarren/gone/gone/generate/consts"
	"strings"
)

// "os"

func main() {
	// Get all .gds files in target directory

	converter := GenerateConverter("*app.AppData")

	fmt.Printf("generated converter: %s\n", converter)
	/*
		find instances of

		gone struct {

		}

	*/

	// Generate a corresponding type in .gds.go

	/*
		Generate a Datastore struct that can be
		and passed to Route
	*/

	/*
		Generate a keytype nonce and keyvalue nonce that
		can be used to get the data and ensure that it
		is the correct type before returning it
	*/

	/*

		--> Look for instances of Routes in source code and
		    generate a store that takes an interface and converts
		    it to the type of the object that is passed in as the
		    interface to Route. It can provide an error if the
		    data is not of the correct type.
	*/
}

func RemoveUntilPeriod(dataType string) string {
	strs := strings.SplitAfter(dataType, ".")
	return strs[len(strs)-1]
}

func RemoveAsterisk(dataType string) string {
	return strings.TrimLeft(dataType, " *")
}

func GenerateConverter(dataType string) string {
	noAsterisk := RemoveAsterisk(dataType)
	noPeriod := RemoveUntilPeriod(noAsterisk)
	return fmt.Sprintf(consts.Converter, noPeriod, dataType, dataType, dataType)
}
