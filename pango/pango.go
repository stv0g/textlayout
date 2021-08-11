package pango

import (
	"fmt"
	"math"

	"github.com/benoitkugler/textlayout/fribidi"
	"golang.org/x/text/width"
)

// enables additional checks, to use only during developpement or testing
const debugMode = false

// reference commit : 23e033aead84005d5d253de7f5e891b9eef60baa

// assert is only used in debug mode
func assert(b bool, msg string) {
	if !b {
		panic("assertion error: " + msg)
	}
}

func showDebug(where string, line *layoutLineData, state *ParaBreakState) {
	lineWidth := line.getWidth()

	fmt.Printf("\trem %d + line %d = %d		%s\n",
		state.remaining_width,
		lineWidth,
		state.remaining_width+lineWidth,
		where)
}

type Fl = float32

// Alignment describes how to align the lines of a `Layout` within the
// available space. If the `Layout` is set to justify
// using SetJustify(), this only has effect for partial lines.
type Alignment uint8

const (
	ALIGN_LEFT   Alignment = iota // Put all available space on the right
	ALIGN_CENTER                  // Center the line within the available space
	ALIGN_RIGHT                   // Put all available space on the left
)

// Rectangle represents a rectangle. It is frequently
// used to represent the logical or ink extents of a single glyph or section
// of text. (See, for instance, pango_font_get_glyph_extents())
type Rectangle struct {
	X      int32 // X coordinate of the left side of the rectangle.
	Y      int32 // Y coordinate of the the top side of the rectangle.
	Width  int32 // width of the rectangle.
	Height int32 // height of the rectangle.
}

const maxInt = int(^uint32(0) >> 1)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func min32(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}

func minL(a, b fribidi.Level) fribidi.Level {
	if a > b {
		return b
	}
	return a
}

func maxG(a, b GlyphUnit) GlyphUnit {
	if a < b {
		return b
	}
	return a
}

// pangoIsZeroWidth checks `ch` to see if it is a character that should not be
// normally rendered on the screen.  This includes all Unicode characters
// with "ZERO WIDTH" in their name, as well as bidi formatting characters, and
// a few other ones.
func pangoIsZeroWidth(ch rune) bool {
	//  00AD  SOFT HYPHEN
	//  034F  COMBINING GRAPHEME JOINER
	//
	//  200B  ZERO WIDTH SPACE
	//  200C  ZERO WIDTH NON-JOINER
	//  200D  ZERO WIDTH JOINER
	//  200E  LEFT-TO-RIGHT MARK
	//  200F  RIGHT-TO-LEFT MARK
	//
	//  2028  LINE SEPARATOR
	//
	//  202A  LEFT-TO-RIGHT EMBEDDING
	//  202B  RIGHT-TO-LEFT EMBEDDING
	//  202C  POP DIRECTIONAL FORMATTING
	//  202D  LEFT-TO-RIGHT OVERRIDE
	//  202E  RIGHT-TO-LEFT OVERRIDE
	//
	//  2060  WORD JOINER
	//  2061  FUNCTION APPLICATION
	//  2062  INVISIBLE TIMES
	//  2063  INVISIBLE SEPARATOR
	//
	//  FEFF  ZERO WIDTH NO-BREAK SPACE
	return (ch & ^0x007F == 0x2000 &&
		((ch >= 0x200B && ch <= 0x200F) || (ch >= 0x202A && ch <= 0x202E) ||
			(ch >= 0x2060 && ch <= 0x2063) || (ch == 0x2028))) ||
		(ch == 0x00AD || ch == 0x034F || ch == 0xFEFF)
}

// return true for east asian wide characters
func isWide(r rune) bool {
	switch width.LookupRune(r).Kind() {
	case width.EastAsianFullwidth, width.EastAsianWide:
		return true
	default:
		return false
	}
}

// Matrix is a transformation between user-space
// coordinates and device coordinates. The transformation
// is given by
//
// x_device = x_user * matrix.xx + y_user * matrix.xy + matrix.x0;
// y_device = x_user * matrix.yx + y_user * matrix.yy + matrix.y0;
type Matrix struct {
	Xx, Xy, Yx, Yy, X0, Y0 Fl
}

var Identity = Matrix{1, 0, 0, 1, 0, 0}

// GetFontScaleFactors calculates the scale factor of a matrix on the width and height of the font.
// That is, `xscale` is the scale factor in the direction of the X coordinate,
// and `yscale` is the scale factor in the direction perpendicular to the
// vector that the X coordinate is mapped to.
//
// Note that output numbers will always be non-negative.
func (matrix *Matrix) GetFontScaleFactors() (xscale, yscale Fl) {
	// Based on cairo-matrix.c:_cairo_matrix_compute_scale_factors()
	// Copyright 2005, Keith Packard
	if matrix == nil {
		return 1, 1
	}

	x := matrix.Xx
	y := matrix.Yx
	xscale = float32(math.Sqrt(float64(x*x + y*y)))

	if xscale != 0 {
		det := matrix.Xx*matrix.Yy - matrix.Yx*matrix.Xy

		/*
		* ignore mirroring
		 */
		if det < 0 {
			det = -det
		}

		yscale = det / xscale
	} else {
		yscale = 0.
	}
	return
}
