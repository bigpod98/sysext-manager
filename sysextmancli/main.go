package main

import (
	"context"
	"fmt"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func main() {
	fmt.Println("well")

	fs, err := file.New("/tmp/")
	if err != nil {
		panic(err)
	}
	defer fs.Close()
	ctx := context.Background()

	mediaType := "application/vnd.test.file"
	fileNames := []string{"/tmp/myfile"}
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))
	for _, name := range fileNames {
		fileDescriptor, err := fs.Add(ctx, name, mediaType, "")
		if err != nil {
			panic(err)
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("file descriptor for %s: %v\n", name, fileDescriptor)
	}

	artifactType := "application/vnd.test.artifact"
	opts := oras.PackManifestOptions{
		Layers: fileDescriptors,
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("manifest descriptor:", manifestDescriptor)

	tag := "latest"
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		panic(err)
	}

	reg := "ghcr.io"
	repo, err := remote.NewRepository(reg + "/bigpod98/test")
	if err != nil {
		panic(err)
	}
	// Note: The below code can be omitted if authentication is not required
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: "",
			Password: "",
		}),
	}

	// 4. Copy from the file store to the remote repository
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	}
}
