package util

import (
	"fmt"
	"strings"

	//  Import the crypto sha256 algorithm for the docker image parser to work
	_ "crypto/sha256"
	//  Import the crypto/sha512 algorithm for the docker image parser to work with 384 and 512 sha hashes
	_ "crypto/sha512"

	dockerref "github.com/docker/distribution/reference"
)

const (
	DefaultImageTag = "latest"
)

// ParseImageName parses a docker image string into three parts: repo, tag and digest.
// If both tag and digest are empty, a default image tag will be returned.
func ParseImageName(image string) (string, string, string, error) {
	named, err := dockerref.ParseNormalizedNamed(image)
	if err != nil {
		return "", "", "", fmt.Errorf("couldn't parse image name: %v", err)
	}

	repoToPull := named.Name()
	var tag, digest string

	tagged, ok := named.(dockerref.Tagged)
	if ok {
		tag = tagged.Tag()
	}

	digested, ok := named.(dockerref.Digested)
	if ok {
		digest = digested.Digest().String()
	}
	// If no tag was specified, use the default "latest".
	if len(tag) == 0 && len(digest) == 0 {
		tag = DefaultImageTag
	}
	return repoToPull, tag, digest, nil
}

func GetImageId(image string, imageId string) string {
	if strings.HasPrefix(imageId, "docker-pullable://") {
		img := strings.TrimPrefix(imageId, "docker-pullable://")
		named, _ := dockerref.ParseNormalizedNamed(img)

		return named.String()

	} else {
		return fmt.Sprintf("%v@%v", image, imageId)
	}
}
