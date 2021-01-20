package fcfonts

import (
	"github.com/benoitkugler/textlayout/fontconfig"
	fc "github.com/benoitkugler/textlayout/fontconfig"
	"github.com/benoitkugler/textlayout/pango"
)

type coverage fc.Charset

// Convert the given `charset` into a new Coverage object.
func fromCharset(charset fc.Charset) pango.Coverage {
	return (*coverage)(&charset)
}

// Get returns true if the rune is covered
func (c *coverage) Get(index rune) bool { return (*fc.Charset)(c).HasChar(index) }

func (c *coverage) Set(index rune, covered bool) {
	if covered {
		(*fc.Charset)(c).AddChar(index)
	} else {
		(*fc.Charset)(c).DelChar(index)
	}
}

// Copy returns a deep copy of the coverage
func (c *coverage) Copy() pango.Coverage {
	if c == nil {
		return c
	}
	cs := (*fc.Charset)(c).Copy()
	return (*coverage)(&cs)
}

// Decoder represents a decoder that an application provides
// for handling a font that is encoded in a custom way.
type Decoder interface {
	// GetCharset returns a charset given a font that
	// includes a list of supported characters in the font.
	// The implementation must be fast because the method is called
	// separately for each character to determine Unicode coverage.
	GetCharset(font *fcFont) fc.Charset

	// GetGlyph returns a single glyph for a given Unicode code point.
	GetGlyph(font *fcFont, r rune) pango.Glyph
}

type fcFontKeyHash struct {
	pattern     string
	matrix      pango.Matrix
	context_key int
	variations  string
}

type fontHash map[fcFontKeyHash]*Font // (GHashFunc)pango_fc_font_key_hash,  (GEqualFunc)pango_fc_font_key_equal

func (m fontHash) lookup(p PangoFcFontKey) *Font {
	key := fcFontKeyHash{pattern: p.pattern.Hash(), matrix: p.matrix,
		context_key: p.context_key, variations: p.variations}
	return m[key]
}

func (m fontHash) insert(p PangoFcFontKey, v *Font) {
	key := fcFontKeyHash{pattern: p.pattern.Hash(), matrix: p.matrix,
		context_key: p.context_key, variations: p.variations}
	m[key] = v
}

func (m fontHash) remove(p PangoFcFontKey) {
	key := fcFontKeyHash{pattern: p.pattern.Hash(), matrix: p.matrix,
		context_key: p.context_key, variations: p.variations}
	delete(m, key)
}

type fontsetHash map[PangoFcFontsetKey]*fcFontset

func (m fontsetHash) lookup(p PangoFcFontsetKey) *fcFontset {
	p.desc = p.desc.AsHash()
	p.fontmap = nil
	return m[p]
}

func (m fontsetHash) insert(p PangoFcFontsetKey, v *fcFontset) {
	p.desc = p.desc.AsHash()
	p.fontmap = nil
	m[p] = v
}

func (m fontsetHash) remove(p PangoFcFontsetKey) {
	p.desc = p.desc.AsHash()
	p.fontmap = nil
	delete(m, p)
}

type fcPatternHash map[string]*fcPatterns

func (m fcPatternHash) lookup(p fontconfig.Pattern) *fcPatterns { return m[p.Hash()] }

func (m fcPatternHash) insert(p fontconfig.Pattern, pts *fcPatterns) { m[p.Hash()] = pts }

// ------------------------------------------------------------------------------------