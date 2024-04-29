package libsysextman

import (
	"fmt"
)

func Push(location string) {
	fmt.Println("pulling from repository")
	pulldistributionregistry(location)
}

func pushdistributionregistry(location string) {
	fmt.Println("pulling from OCI distribution registry")
}

