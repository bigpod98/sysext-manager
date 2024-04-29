package libsysextman

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Pull(location string) {
	fmt.Println("pulling from repository")
	pulldistributionregistry(location)

}

func pulldistributionregistry(location string) {
	fmt.Println("pulling from OCI distribution registry")
	registry := ""
	namespace := ""
	tag := ""
	//"	/v2/<name>/manifests/<reference>" endpoint-3
	_, err := http.Get(fmt.Sprintf("https://%s//v2/%s/manifests/%s", registry, namespace, tag))
	if err != nil {
		panic(err)
	}
	//read the manifest, check if single layer has right format, then pull the the layer

	//" /v2/<name>/blobs/<digest>" endpoint-2
	digest := ""

	resp, err := http.Get(fmt.Sprintf("https://%s//v2/%s/blobs/%s", registry, namespace, digest))
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(resp.Body)
	responseData, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	err1 := os.WriteFile("/tmp/dat1", responseData, 0644)
	if err1 != nil {
		panic(err1)
	}
}
