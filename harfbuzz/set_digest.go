package harfbuzz

import "github.com/benoitkugler/textlayout/fonts"

// ported from src/hb-set-digest.hh Copyright © 2012  Google, Inc. Behdad Esfahbod

const maskBits = 4 * 8 // 4 = size(setDigestLowestBits)

type setType = fonts.GlyphIndex

type setDigestLowestBits uint32

func maskFor(g setType, shift int) setDigestLowestBits {
	return 1 << ((g >> shift) & (maskBits - 1))
}

func (sd *setDigestLowestBits) add(g setType, shift int) { *sd |= maskFor(g, shift) }

func (sd *setDigestLowestBits) addRange(a, b setType, shift int) {
	if (b>>shift)-(a>>shift) >= maskBits-1 {
		*sd = ^setDigestLowestBits(0)
	} else {
		mb := maskFor(b, shift)
		ma := maskFor(a, shift)
		var op setDigestLowestBits
		if mb < ma {
			op = 1
		}
		*sd |= mb + (mb - ma) - op
	}
}

func (sd *setDigestLowestBits) addArray(arr []setType, shift int) {
	for _, v := range arr {
		sd.add(v, shift)
	}
}

func (sd setDigestLowestBits) mayHave(g setType, shift int) bool {
	return sd&maskFor(g, shift) != 0
}

/* This is a combination of digests that performs "best".
 * There is not much science to this: it's a result of intuition
 * and testing. */
const (
	shift0 = 4
	shift1 = 0
	shift2 = 9
)

// SetDigest implement various "filters" that support
// "approximate member query".  Conceptually these are like Bloom
// Filter and Quotient Filter, however, much smaller, faster, and
// designed to fit the requirements of our uses for glyph coverage
// queries.
//
// Our filters are highly accurate if the lookup covers fairly local
// set of glyphs, but fully flooded and ineffective if coverage is
// all over the place.
//
// The frozen-set can be used instead of a digest, to trade more
// memory for 100% accuracy, but in practice, that doesn't look like
// an attractive trade-off.
type SetDigest [3]setDigestLowestBits

// Add adds the given rune to the set.
func (sd *SetDigest) Add(g setType) {
	sd[0].add(g, shift0)
	sd[1].add(g, shift1)
	sd[2].add(g, shift2)
}

// AddRange adds the given, inclusive range to the set,
// in an efficient manner.
func (sd *SetDigest) AddRange(a, b setType) {
	sd[0].addRange(a, b, shift0)
	sd[1].addRange(a, b, shift1)
	sd[2].addRange(a, b, shift2)
}

// AddArray is a convenience method to add
// many runes.
func (sd *SetDigest) AddArray(arr []setType) {
	sd[0].addArray(arr, shift0)
	sd[1].addArray(arr, shift1)
	sd[2].addArray(arr, shift2)
}

// MayHave performs an "approximate member query": if the return value
// is `false`, then it is certain that `g` is not in the set.
// Otherwise, we don't kwow, it might be a false positive.
// Note that runes in the set are certain to return `true`.
func (sd SetDigest) MayHave(g setType) bool {
	return sd[0].mayHave(g, shift0) && sd[1].mayHave(g, shift1) && sd[2].mayHave(g, shift2)
}