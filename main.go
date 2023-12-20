package main

import (
	// "context"
	// "encoding/json"
	// "fmt"

	// "github.com/docker/docker/client"

	"context"
	"fmt"
	"log"
	"os"
	"time"

	manifesttypes "github.com/docker/cli/cli/manifest/types"
	"github.com/docker/distribution"

	"github.com/distribution/reference"
	"github.com/docker/cli/cli/registry/client"
	registrytypes "github.com/docker/docker/api/types/registry"
)

func main() {
	startTime := time.Now()

	if len(os.Args) < 2 {
		log.Fatal("please specify image url e.g., ubuntu:latst")
	}

	resolver := func(ctx context.Context, index *registrytypes.IndexInfo) registrytypes.AuthConfig {
		return registrytypes.AuthConfig{}
	}
	c := client.NewRegistryClient(resolver, "suraj's macbook", false)
	ref, err := normalizeReference(os.Args[1])
	if err != nil {
		panic(err)
	}

	var mlist []manifesttypes.ImageManifest
	imageManifest, err := c.GetManifest(context.Background(), ref)
	if err == nil {
		mlist = []manifesttypes.ImageManifest{
			imageManifest,
		}

	} else {
		mlist, err = c.GetManifestList(context.Background(), ref)
		if err != nil {
			panic(err)
		}
	}

	// os.Setenv("DOCKER_DEFAULT_PLATFORM", "linux/amd64")

	// b, _ := json.MarshalIndent(mlist, " ", " ")
	// fmt.Println("t", string(b))
	// fmt.Println("===============================")

	for _, m := range mlist {
		if m.Descriptor.Platform.OS == "linux" &&
			m.Descriptor.Platform.Architecture == "amd64" {
			var sizeInMB int64

			var layers []distribution.Descriptor
			if m.OCIManifest != nil {
				layers = m.OCIManifest.Layers
			} else {
				layers = m.SchemaV2Manifest.Layers
			}

			for _, l := range layers {
				sizeInMB += l.Size
			}

			fmt.Printf("%s is %vMB in size\n", ref.String(), sizeInMB/(1024*1024))

			// b, _ := json.MarshalIndent(m.OCIManifest, " ", " ")
			// fmt.Println("t", string(b))
			// fmt.Println("===============================")
		}
	}

	endTime := time.Now()

	fmt.Println("======================")
	fmt.Printf("that took %v seconds\n", endTime.Sub(startTime).Seconds())

}

// copy-paste from https://github.com/docker/cli/blob/fb2ba5d63ba4166ceeefa21c2fd866b06966874e/cli/command/manifest/util.go#L59-L68
func normalizeReference(ref string) (reference.Named, error) {
	namedRef, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return nil, err
	}
	if _, isDigested := namedRef.(reference.Canonical); !isDigested {
		return reference.TagNameOnly(namedRef), nil
	}
	return namedRef, nil
}
