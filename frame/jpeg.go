// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frame

import (
	"bytes"
	"image"
	"image/jpeg"
	"runtime"
)

type fJPEG struct {
	image.Image
	release func()
}

// Register this framer for this format.
func init() {
	RegisterFramer("JPEG", newJPEGFramer)
}

// Return a framer for JPEG.
func newJPEGFramer(w, h, stride, size int) func([]byte, func()) (Framer, error) {
	return jpegFramer
}

// Wrap a jpeg block in a Framer so that it can be used as an image.
func jpegFramer(f []byte, rel func()) (Framer, error) {
	img, err := jpeg.Decode(bytes.NewBuffer(f))
	if err != nil {
		if rel != nil {
			rel()
		}
		return nil, err
	}
	fr := &fJPEG{Image: img, release: rel}
	runtime.SetFinalizer(fr, func(obj Framer) {
		obj.Release()
	})
	return fr, nil
}

// Done with frame, release back to camera (if required).
func (f *fJPEG) Release() {
	if f.release != nil {
		f.release()
		// Make sure it only gets called once.
		f.release = nil
	}
}
