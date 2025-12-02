// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package himage provides some high level image types and interfaces.
package himage

import "image"

// AnimatedImage represents an animated image.
// This is currently supported for GIF and WebP images.
type AnimatedImage interface {
	image.Image        // The first frame.
	GetRaw() any       // *gif.GIF or *WEBP.
	GetLoopCount() int // Number of times to loop the animation. 0 means infinite.
	ImageFrames
}

// ImageFrames provides access to the frames of an animated image.
type ImageFrames interface {
	GetFrames() []image.Image

	// Frame durations in milliseconds.
	// Note that Gif frame durations are in 100ths of a second,
	// so they need to be multiplied by 10 to get milliseconds and vice versa.
	GetFrameDurations() []int

	SetFrames(frames []image.Image)
	SetWidthHeight(width, height int)
}

// ImageConfigProvider provides access to the image.Config of an image.
type ImageConfigProvider interface {
	GetImageConfig() image.Config
}

// FrameDurationsToGifDelays converts frame durations in milliseconds to
// GIF delays in 100ths of a second.
func FrameDurationsToGifDelays(frameDurations []int) []int {
	delays := make([]int, len(frameDurations))
	for i, fd := range frameDurations {
		delays[i] = fd / 10
		if delays[i] == 0 && fd > 0 {
			delays[i] = 1
		}
	}
	return delays
}

// GifDelaysToFrameDurations converts GIF delays in 100ths of a second to
// frame durations in milliseconds.
func GifDelaysToFrameDurations(delays []int) []int {
	frameDurations := make([]int, len(delays))
	for i, d := range delays {
		frameDurations[i] = d * 10
	}
	return frameDurations
}
