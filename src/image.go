package main

import

func init() {
    emptyImage, _ = NewImage(16, 16, FilterDefault)
}

func NewImage(width, height int, filter Filter) (*Image, error) {
    s := shareable.NewImage(width, height)
    i := &Image{
        shareableImages: []*shareable.Image{s},
        filter:          filter,
    }
    i.addr = i
    runtime.SetFinalizer(i, (*Image).Dispose)
    return i, nil
}

type Image struct {
    // addr holds self to check copying.
    // See strings.Builder for similar examples.
    addr *Image

    // shareableImages is a set of shareable.Image sorted by the order of mipmap level.
    // The level 0 image is a regular image and higher-level images are used for mipmap.
    shareableImages []*shareable.Image

    filter Filter
}

