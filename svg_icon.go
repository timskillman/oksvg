// Copyright 2017 The oksvg Authors. All rights reserved.
// created: 2/12/2017 by S.R.Wiley
//
// utils.go implements translation of an SVG2.0 path into a rasterx Path.

package oksvg

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/timskillman/rasterx"
)

// SvgIcon holds data from parsed SVGs.
type SvgIcon struct {
	ViewBox      struct{ X, Y, W, H float64 }
	Titles       []string // Title elements collect here
	Descriptions []string // Description elements collect here
	Grads        map[string]*rasterx.Gradient
	Defs         map[string][]definition
	SVGPaths     []SvgPath
	Transform    rasterx.Matrix2D
	classes      map[string]styleAttribute
}

// Draw the compiled SVG icon into the GraphicContext.
// All elements should be contained by the Bounds rectangle of the SvgIcon.
func (s *SvgIcon) Draw(r *rasterx.Dasher, opacity float64) {
	for _, svgp := range s.SVGPaths {
		svgp.DrawTransformed(r, opacity, s.Transform)
	}
}

// SetTarget sets the Transform matrix to draw within the bounds of the rectangle arguments
func (s *SvgIcon) SetTarget(x, y, w, h float64) {
	scaleW := w / s.ViewBox.W
	scaleH := h / s.ViewBox.H
	s.Transform = rasterx.Identity.Translate(x-s.ViewBox.X, y-s.ViewBox.Y).Scale(scaleW, scaleH)
}

// **NEW** Returns the SvgIcon as an image set to a given width and height.
// However, if width is set to -1 then the original width of the SvgIcon is used.
// If the height is set to -1 then the SvgIcon maintains its aspect ratio even when
// an arbitrary width is set
func (s *SvgIcon) AsImage(width float64, height float64) image.Image {
	if width < 1 {
		width = s.ViewBox.W
	}
	if height < 1 {
		sc := width / s.ViewBox.W
		height = s.ViewBox.H * sc
	}
	s.SetTarget(0, 0, width, height)
	w, h := int(width), int(height)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scannerGV := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scannerGV)
	s.Draw(raster, 1.0)
	return img
}

// **NEW** The SvgIcon is saved as a PNG file set to a given width and height.
// However, if width is set to -1 then the original width of the SvgIcon is used.
// If the height is set to -1 then the SvgIcon maintains its aspect ratio even when
// an arbitrary width is set
func (s *SvgIcon) SaveAsPng(filePath string, w float64, h float64) error {
	return s.saveImage(filePath, s.AsImage(w, h), true)
}

// **NEW** The SvgIcon is saved as a JPEG file set to a given width and height.
// However, if width is set to -1 then the original width of the SvgIcon is used.
// If the height is set to -1 then the SvgIcon maintains its aspect ratio even when
// an arbitrary width is set
func (s *SvgIcon) SaveAsJpeg(filePath string, w float64, h float64) error {
	return s.saveImage(filePath, s.AsImage(w, h), false)
}

func (s *SvgIcon) saveImage(filePath string, m image.Image, asPng bool) error {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create Writer from file
	w := bufio.NewWriter(f)

	// Write the image as either PNG or JPEG into the buffer
	if asPng {
		if err := png.Encode(w, m); err != nil {
			return err
		}
	} else {
		if err := jpeg.Encode(w, m, nil); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
