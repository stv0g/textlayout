package truetype

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// TableName represents the OpenType 'name' table. This contains
// human-readable meta-data about the font, for example the Author
// and Copyright.
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6name.html
type TableName []NameEntry

// returns the name entry with `name`, for both plaftorm,
// or nil if not found
func (t TableName) getEntry(name NameID) (windows, mac *NameEntry) {
	for _, e := range t {
		if e.NameID == name {
			if e.isWindows() && (e.LanguageID == PLMicrosoftEnglish || windows == nil) {
				windows = &e
			}
			if e.isMac() && (e.LanguageID == PLMacEnglish || mac == nil) {
				mac = &e
			}
		}
	}
	return windows, mac
}

//convert a UTF-16 name entry to ASCII
func asciiFromUTF16(entry []byte) string {
	if len(entry)%2 != 0 {
		entry = append(entry, 0)
	}
	out := make([]byte, len(entry)/2)
	for n := range out {
		code := binary.BigEndian.Uint16(entry[2*n:])
		if code < 32 || code > 127 {
			code = '?'
		}
		out[n] = byte(code)
	}
	return string(out)
}

// convert an Apple Roman or symbol name entry to ASCII
func asciiFromOther(entry []byte) string {
	out := make([]byte, len(entry))
	for n, code := range entry {
		if code < 32 || code > 127 {
			code = '?'
		}
		out[n] = code
	}
	return string(out)
}

// return an empty string is not found
func (names TableName) getName(name NameID) string {
	var (
		foundApple        = -1
		foundAppleRoman   = -1
		foundAppleEnglish = -1
		foundWin          = -1
		foundUnicode      = -1
		isEnglish         = false
	)

	for n, rec := range names {
		// According to the OpenType 1.3 specification, only Microsoft or
		// Apple platform IDs might be used in the `name' table.  The
		// `Unicode' platform is reserved for the `cmap' table, and the
		// `ISO' one is deprecated.
		//
		// However, the Apple TrueType specification doesn't say the same
		// thing and goes to suggest that all Unicode `name' table entries
		// should be coded in UTF-16 (in big-endian format I suppose).
		if rec.NameID == name && len(rec.Value) > 0 {
			switch rec.PlatformID {
			case PlatformUnicode, PlatformIso:
				// there is `languageID' to check there.  We should use this
				// field only as a last solution when nothing else is
				// available.
				foundUnicode = n
			case PlatformMac:
				// This is a bit special because some fonts will use either
				// an English language id, or a Roman encoding id, to indicate
				// the English version of its font name.
				if rec.LanguageID == PLMacEnglish {
					foundAppleEnglish = n
				} else if rec.EncodingID == PEMacRoman {
					foundAppleRoman = n
				}
			case PlatformMicrosoft:
				// we only take a non-English name when there is nothing
				// else available in the font
				if foundWin == -1 || (rec.LanguageID&0x3FF) == 0x009 {
					switch rec.EncodingID {
					case PEMicrosoftSymbolCs, PEMicrosoftUnicodeCs, PEMicrosoftUcs4:
						isEnglish = (rec.LanguageID & 0x3FF) == 0x009
						foundWin = n
					}
				}
			}
		}
	}

	foundApple = foundAppleRoman
	if foundAppleEnglish >= 0 {
		foundApple = foundAppleEnglish
	}

	// some fonts contain invalid Unicode or Macintosh formatted entries;
	// we will thus favor names encoded in Windows formats if available
	// (provided it is an English name)
	var (
		convert func(entry []byte) string
		rec     NameEntry
	)
	if foundWin >= 0 && !(foundApple >= 0 && !isEnglish) {
		rec = names[foundWin]
		switch rec.EncodingID {
		// all Unicode strings are encoded using UTF-16BE
		case PEMicrosoftUnicodeCs, PEMicrosoftSymbolCs:
			convert = asciiFromUTF16
		case PEMicrosoftUcs4:
			// Apparently, if this value is found in a name table entry, it is
			// documented as `full Unicode repertoire'.  Experience with the
			// MsGothic font shipped with Windows Vista shows that this really
			// means UTF-16 encoded names (UCS-4 values are only used within
			// charmaps).
			convert = asciiFromUTF16
		}
	} else if foundApple >= 0 {
		rec = names[foundApple]
		convert = asciiFromOther
	} else if foundUnicode >= 0 {
		rec = names[foundUnicode]
		convert = asciiFromUTF16
	}

	if convert != nil {
		return convert(rec.Value)
	}

	return ""
}

type nameHeader struct {
	Format       uint16
	Count        uint16
	StringOffset uint16
}

// PlatformID represents the platform id for entries in the name table.
type PlatformID uint16

const (
	PlatformUnicode PlatformID = iota
	PlatformMac
	PlatformIso /* deprecated */
	PlatformMicrosoft
	PlatformCustom
	_
	_
	PlatformAdobe /* artificial */
)

// String returns an idenfying string for each platform or "Platform X" for unknown values.
func (p PlatformID) String() string {
	switch p {
	case PlatformUnicode:
		return "Unicode"
	case PlatformMac:
		return "Mac"
	case PlatformMicrosoft:
		return "Microsoft"
	default:
		return "Platform " + strconv.Itoa(int(p))
	}
}

// PlatformEncodingID represents the platform specific id for entries in the name table.
// The most common values are provided as constants.
type PlatformEncodingID uint16

const (
	PEUnicodeDefault     = PlatformEncodingID(0)
	PEMacRoman           = PEUnicodeDefault
	PEMicrosoftSymbolCs  = PlatformEncodingID(0)
	PEMicrosoftUnicodeCs = PlatformEncodingID(1)
	PEMicrosoftUcs4      = PlatformEncodingID(10)
)

// PlatformLanguageID represents the language used by an entry in the name table,
// the three most common values are provided as constants.
type PlatformLanguageID uint16

const (
	PLMacEnglish       = PlatformLanguageID(0)
	PLUnicodeDefault   = PlatformLanguageID(0)
	PLMicrosoftEnglish = PlatformLanguageID(0x0409)
)

// NameID is the ID for entries in the font table.
type NameID uint16

const (
	NameCopyrightNotice NameID = iota
	NameFontFamily
	NameFontSubfamily
	NameUniqueIdentifier
	NameFull
	NameVersion
	NamePostscript
	NameTrademark
	NameManufacturer
	NameDesigner
	NameDescription
	NameVendorURL
	NameDesignerURL
	NameLicenseDescription
	_NameReserved
	NameLicenseURL
	NamePreferredFamily    // or Typographic Family
	NamePreferredSubfamily // or Typographic Subfamily
	NameCompatibleFull
	NameSampleText
	NamePostscriptCID
	NameWWSFamily
	NameWWSSubfamily
	NameLightBackgroundPalette
	NameDarkBackgroundPalette
)

// String returns an identifying
func (nameId NameID) String() string {
	switch nameId {
	case NameCopyrightNotice:
		return "Copyright Notice"
	case NameFontFamily:
		return "Font Family"
	case NameFontSubfamily:
		return "Font Subfamily"
	case NameUniqueIdentifier:
		return "Unique Identifier"
	case NameFull:
		return "Full Name"
	case NameVersion:
		return "Version"
	case NamePostscript:
		return "PostScript Name"
	case NameTrademark:
		return "Trademark Notice"
	case NameManufacturer:
		return "Manufacturer"
	case NameDesigner:
		return "Designer"
	case NameDescription:
		return "Description"
	case NameVendorURL:
		return "Vendor URL"
	case NameDesignerURL:
		return "Designer URL"
	case NameLicenseDescription:
		return "License Description"
	case NameLicenseURL:
		return "License URL"
	case NamePreferredFamily:
		return "Preferred Family"
	case NamePreferredSubfamily:
		return "Preferred Subfamily"
	case NameCompatibleFull:
		return "Compatible Full"
	case NameSampleText:
		return "Sample Text"
	case NamePostscriptCID:
		return "PostScript CID"
	case NameWWSFamily:
		return "WWS Family"
	case NameWWSSubfamily:
		return "WWS Subfamily"
	case NameLightBackgroundPalette:
		return "Light Background Palette"
	case NameDarkBackgroundPalette:
		return "Dark Background Palette"
	default:
		return "Name " + strconv.Itoa(int(nameId))
	}

}

type nameRecord struct {
	PlatformID PlatformID
	EncodingID PlatformEncodingID
	LanguageID PlatformLanguageID
	NameID     NameID
	Length     uint16
	Offset     uint16
}

type NameEntry struct {
	PlatformID PlatformID
	EncodingID PlatformEncodingID
	LanguageID PlatformLanguageID
	NameID     NameID
	Value      []byte
}

func (n NameEntry) isWindows() bool {
	return n.PlatformID == PlatformMicrosoft && (n.EncodingID == PEMicrosoftUnicodeCs || n.EncodingID == PEUnicodeDefault)
}

func (n NameEntry) isMac() bool {
	return n.PlatformID == PlatformMac && n.EncodingID == PEMacRoman
}

// String is a best-effort attempt to get a UTF-8 encoded version of
// Value. Only MicrosoftUnicode (3,1 ,X), MacRomain (1,0,X) and Unicode platform
// strings are supported.
func (nameEntry *NameEntry) String() string {

	if nameEntry.PlatformID == PlatformUnicode || (nameEntry.PlatformID == PlatformMicrosoft &&
		nameEntry.EncodingID == PEMicrosoftUnicodeCs) {

		decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()

		outstr, _, err := transform.String(decoder, string(nameEntry.Value))

		if err == nil {
			return outstr
		}
	}

	if nameEntry.isMac() {

		decoder := charmap.Macintosh.NewDecoder()

		outstr, _, err := transform.String(decoder, string(nameEntry.Value))

		if err == nil {
			return outstr
		}
	}

	return string(nameEntry.Value)
}

func (nameEntry *NameEntry) Label() string {
	return nameEntry.NameID.String()
}

func (nameEntry *NameEntry) Platform() string {
	return nameEntry.PlatformID.String()
}

func parseTableName(buf []byte) (TableName, error) {
	r := bytes.NewReader(buf)

	var header nameHeader
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	table := make(TableName, 0, header.Count)

	for i := 0; i < int(header.Count); i++ {
		var record nameRecord
		if err := binary.Read(r, binary.BigEndian, &record); err != nil {
			return nil, err
		}

		start := header.StringOffset + record.Offset
		end := start + record.Length

		if int(start) > len(buf) || int(end) > len(buf) {
			return nil, io.ErrUnexpectedEOF
		}

		table = append(table, NameEntry{
			record.PlatformID,
			record.EncodingID,
			record.LanguageID,
			record.NameID,
			buf[start:end],
		})
	}

	return table, nil
}