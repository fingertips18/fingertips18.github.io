// Package metadata provides utilities for computing blurhashes.
package metadata

import (
	"image"

	bh "github.com/buckket/go-blurhash"
)

type BlurHashAPI interface {
	Encode(xComponents, yComponents int, rgba image.Image) (string, error)
	Decode(hash string, width, height, punch int) (image.Image, error)
	IsValid(hash string) bool
}

type blurHashImpl struct{}

// Encode encodes the provided image into a BlurHash string using the specified
// number of horizontal (xComponents) and vertical (yComponents) frequency
// components. Higher component counts yield more detail and a longer hash,
// while lower counts produce shorter, blurrier representations.
//
// The rgba parameter is the source image to encode. The method returns the
// encoded BlurHash string and an error if encoding fails (for example when the
// component counts are outside the valid range, the image is nil, or the image
// cannot be processed).
//
// Choose component counts to balance compactness and fidelity: small values
// produce compact, low-detail hashes; larger values preserve more visual
// information at the cost of a longer string.
func (b *blurHashImpl) Encode(xComponents, yComponents int, rgba image.Image) (string, error) {
	return bh.Encode(xComponents, yComponents, rgba)
}

// Decode decodes a BlurHash-encoded string into an image.Image with the
// requested pixel dimensions and punch (contrast) factor. The hash parameter
// must contain a valid BlurHash; width and height specify the output image
// dimensions in pixels and should be positive. The punch parameter adjusts the
// contrast/strength of the reconstruction (values around 1.0 produce a faithful
// reconstruction; higher values increase contrast).
//
// It returns the decoded image on success, or an error if the hash is invalid
// or decoding fails.
func (b *blurHashImpl) Decode(hash string, width, height, punch int) (image.Image, error) {
	return bh.Decode(hash, width, height, punch)
}

// IsValid checks whether the provided hash string is a valid blur hash.
// It returns true if the hash can be successfully decoded into its components,
// and false otherwise.
func (b *blurHashImpl) IsValid(hash string) bool {
	_, _, err := bh.Components(hash)
	return err == nil
}

// NewBlurHashAPI returns a new instance of the package's default BlurHashAPI
// implementation (blurHashImpl). The returned value implements the
// BlurHashAPI interface and is non-nil, ready for use with its default
// configuration.
func NewBlurHashAPI() BlurHashAPI {
	return &blurHashImpl{}
}
