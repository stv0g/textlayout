package type1

import (
	"encoding/hex"
	"errors"
	"fmt"

	tk "github.com/benoitkugler/pstokenizer"
	"github.com/benoitkugler/textlayout/fonts"
	"github.com/benoitkugler/textlayout/fonts/simpleencodings"
)

// constants for encryption
const (
	EEXEC_KEY      = 55665
	CHARSTRING_KEY = 4330
)

var none = tk.Token{} // null token

type parser struct {
	lexer lexer
}

type lexer struct {
	tk.Tokenizer
}

// constructs a new lexer given a header-less .pfb segment
func newLexer(data []byte) lexer {
	return lexer{*tk.NewTokenizer(data)}
}

func (l *lexer) nextToken() (tk.Token, error) {
	return l.Tokenizer.NextToken()
}

func (l lexer) peekToken() tk.Token {
	t, err := l.Tokenizer.PeekToken()
	if err != nil {
		return none
	}
	return t
}

// Encoding is either the standard encoding, or defined by the font
type Encoding struct {
	Standard bool
	Custom   simpleencodings.Encoding
}

/*
 Parses an Adobe Type 1 (.pfb) font, composed of `segment1` (ASCII) and `segment2` (Binary).
 It is used exclusively in Type1 font.

 The Type 1 font format is a free-text format which is somewhat difficult
 to parse. This is made worse by the fact that many Type 1 font files do
 not conform to the specification, especially those embedded in PDFs. This
 parser therefore tries to be as forgiving as possible.

 See "Adobe Type 1 Font Format, Adobe Systems (1999)"

 Ported from the code from John Hewson
*/
func Parse(segment1, segment2 []byte) (PFBFont, error) {
	p := parser{}
	out, err := p.parseASCII(segment1)
	if err != nil {
		return PFBFont{}, err
	}
	if len(segment2) > 0 {
		p.parseBinary(segment2, &out)
	}
	return out, nil
}

// Parses the ASCII portion of a Type 1 font.
func (p *parser) parseASCII(bytes []byte) (PFBFont, error) {
	if len(bytes) == 0 {
		return PFBFont{}, errors.New("bytes is empty")
	}

	// %!FontType1-1.0
	// %!PS-AdobeFont-1.0
	if len(bytes) < 2 || (bytes[0] != '%' && bytes[1] != '!') {
		return PFBFont{}, errors.New("Invalid start of ASCII segment")
	}

	var out PFBFont
	p.lexer = newLexer(bytes)

	// (corrupt?) synthetic font
	if p.lexer.peekToken().Value == "FontDirectory" {
		if err := p.readWithName(tk.Other, "FontDirectory"); err != nil {
			return out, err
		}
		if _, err := p.read(tk.Name); err != nil { // font name;
			return out, err
		}
		if err := p.readWithName(tk.Other, "known"); err != nil {
			return out, err
		}
		if _, err := p.read(tk.StartProc); err != nil {
			return out, err
		}
		if _, err := p.readProc(); err != nil {
			return out, err
		}
		if _, err := p.read(tk.StartProc); err != nil {
			return out, err
		}
		if _, err := p.readProc(); err != nil {
			return out, err
		}
		if err := p.readWithName(tk.Other, "ifelse"); err != nil {
			return out, err
		}
	}

	// font dict
	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return out, err
	}
	length, _ := lengthT.Int()
	if err := p.readWithName(tk.Other, "dict"); err != nil {
		return out, err
	}
	// found in some TeX fonts
	if _, err := p.readMaybe(tk.Other, "dup"); err != nil {
		return out, err
	}
	// if present, the "currentdict" is not required
	if err := p.readWithName(tk.Other, "begin"); err != nil {
		return out, err
	}

	for i := 0; i < length; i++ {
		// premature end
		token := p.lexer.peekToken()
		if token == none {
			break
		}
		if token.Kind == tk.Other && ("currentdict" == token.Value || "end" == token.Value) {
			break
		}

		// key/value
		keyT, err := p.read(tk.Name)
		if err != nil {
			return out, err
		}
		switch key := keyT.Value; key {
		case "FontInfo", "Fontinfo":
			dict, err := p.readSimpleDict()
			if err != nil {
				return out, err
			}
			out.PSInfo = p.readFontInfo(dict)
		case "Metrics":
			_, err = p.readSimpleDict()
		case "Encoding":
			out.Encoding, err = p.readEncoding()
		default:
			err = p.readSimpleValue(key, &out)
		}
		if err != nil {
			return out, err
		}
	}

	if _, err := p.readMaybe(tk.Other, "currentdict"); err != nil {
		return out, err
	}
	if err := p.readWithName(tk.Other, "end"); err != nil {
		return out, err
	}
	if err := p.readWithName(tk.Other, "currentfile"); err != nil {
		return out, err
	}
	if err := p.readWithName(tk.Other, "eexec"); err != nil {
		return out, err
	}
	return out, nil
}

func (p *parser) readSimpleValue(key string, font *PFBFont) error {
	value, err := p.readDictValue()
	if err != nil {
		return err
	}
	switch key {
	case "FontName", "PaintType", "FontType", "UniqueID", "StrokeWidth", "FID":
		if len(value) == 0 {
			return fmt.Errorf("missing value for key %s", key)
		}
	}
	switch key {
	case "FontName":
		font.FontName = value[0].Value
	case "PaintType":
		font.PaintType, _ = value[0].Int()
	case "FontType":
		font.FontType, _ = value[0].Int()
	case "UniqueID":
		font.UniqueID, _ = value[0].Int()
	case "StrokeWidth":
		font.StrokeWidth, _ = value[0].Float()
	case "FID":
		font.FontID = value[0].Value
	case "FontMatrix":
		font.FontMatrix, err = p.arrayToNumbers(value)
	case "FontBBox":
		font.FontBBox, err = p.arrayToNumbers(value)
	}
	return err
}

func (p *parser) readEncoding() (Encoding, error) {
	var out Encoding
	if p.lexer.peekToken().Kind == tk.Other {
		nameT, err := p.lexer.nextToken()
		if err != nil {
			return out, err
		}
		name_ := nameT.Value
		if name_ == "StandardEncoding" {
			out.Standard = true
		} else {
			return out, errors.New("Unknown encoding: " + name_)
		}
		if _, err := p.readMaybe(tk.Other, "readonly"); err != nil {
			return out, err
		}
		if err := p.readWithName(tk.Other, "def"); err != nil {
			return out, err
		}
	} else {
		if _, err := p.read(tk.Integer); err != nil {
			return out, err
		}
		if _, err := p.readMaybe(tk.Other, "array"); err != nil {
			return out, err
		}

		// 0 1 255 {1 index exch /.notdef put } for
		// we have to check "readonly" and "def" too
		// as some fonts don't provide any dup-values, see PDFBOX-2134
		for !(p.lexer.peekToken().Kind == tk.Other &&
			(p.lexer.peekToken().Value == "dup" ||
				p.lexer.peekToken().Value == "readonly" ||
				p.lexer.peekToken().Value == "def")) {
			_, err := p.lexer.nextToken()
			if err != nil {
				return out, err
			}
		}

		for p.lexer.peekToken().Kind == tk.Other &&
			p.lexer.peekToken().Value == "dup" {
			if err := p.readWithName(tk.Other, "dup"); err != nil {
				return out, err
			}
			codeT, err := p.read(tk.Integer)
			if err != nil {
				return out, err
			}
			code, _ := codeT.Int()
			nameT, err := p.read(tk.Name)
			if err != nil {
				return out, err
			}
			if err := p.readWithName(tk.Other, "put"); err != nil {
				return out, err
			}
			out.Custom[byte(code)] = nameT.Value
		}
		if _, err := p.readMaybe(tk.Other, "readonly"); err != nil {
			return out, err
		}
		if err := p.readWithName(tk.Other, "def"); err != nil {
			return out, err
		}
	}
	return out, nil
}

// Extracts values from an array as numbers.
func (p *parser) arrayToNumbers(value []tk.Token) ([]float64, error) {
	var numbers []float64
	for i, size := 1, len(value)-1; i < size; i++ {
		token := value[i]
		if token.Kind == tk.Float || token.Kind == tk.Integer {
			f, _ := token.Float()
			numbers = append(numbers, f)
		} else {
			return nil, fmt.Errorf("Expected INTEGER or REAL but got %s", token.Kind)
		}
	}
	return numbers, nil
}

// Extracts values from the /FontInfo dictionary.
func (p *parser) readFontInfo(fontInfo map[string][]tk.Token) fonts.PSInfo {
	var out fonts.PSInfo
	for key, value := range fontInfo {
		switch key {
		case "version":
			out.Version = value[0].Value
		case "Notice":
			out.Notice = value[0].Value
		case "FullName":
			out.FullName = value[0].Value
		case "FamilyName":
			out.FamilyName = value[0].Value
		case "Weight":
			out.Weight = value[0].Value
		case "isFixedPitch":
			out.IsFixedPitch = value[0].Value == "true"
		case "ItalicAngle":
			out.ItalicAngle, _ = value[0].Int()
		case "UnderlinePosition":
			out.UnderlinePosition, _ = value[0].Int()
		case "UnderlineThickness":
			out.UnderlineThickness, _ = value[0].Int()
		}
	}
	return out
}

// Reads a dictionary whose values are simple, i.e., do not contain nested dictionaries.
func (p *parser) readSimpleDict() (map[string][]tk.Token, error) {
	dict := map[string][]tk.Token{}

	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return nil, err
	}
	length, _ := lengthT.Int()
	if err := p.readWithName(tk.Other, "dict"); err != nil {
		return nil, err
	}
	if _, err := p.readMaybe(tk.Other, "dup"); err != nil {
		return nil, err
	}
	if err := p.readWithName(tk.Other, "begin"); err != nil {
		return nil, err
	}

	for i := 0; i < length; i++ {
		if p.lexer.peekToken() == none {
			break
		}
		if p.lexer.peekToken().Kind == tk.Other &&
			!(p.lexer.peekToken().Value == "end") {
			if _, err := p.read(tk.Other); err != nil {
				return nil, err
			}
		}
		// premature end
		if p.lexer.peekToken() == none {
			break
		}
		if p.lexer.peekToken().Kind == tk.Other &&
			p.lexer.peekToken().Value == "end" {
			break
		}

		// simple value
		keyT, err := p.read(tk.Name)
		if err != nil {
			return nil, err
		}
		value, err := p.readDictValue()
		if err != nil {
			return nil, err
		}
		dict[keyT.Value] = value
	}

	if err := p.readWithName(tk.Other, "end"); err != nil {
		return nil, err
	}
	if _, err := p.readMaybe(tk.Other, "readonly"); err != nil {
		return nil, err
	}
	if err := p.readWithName(tk.Other, "def"); err != nil {
		return nil, err
	}

	return dict, nil
}

// Reads a simple value from a dictionary.
func (p *parser) readDictValue() ([]tk.Token, error) {
	value, err := p.readValue()
	if err != nil {
		return nil, err
	}
	err = p.readDef()
	return value, err
}

// Reads a simple value. This is either a number, a string,
// a name, a literal name, an array, a procedure, or a charstring.
// This method does not support reading nested dictionaries unless they're empty.
func (p *parser) readValue() ([]tk.Token, error) {
	var value []tk.Token
	token, err := p.lexer.nextToken()
	if err != nil {
		return nil, err
	}
	if p.lexer.peekToken() == none {
		return value, nil
	}
	value = append(value, token)

	switch token.Kind {
	case tk.StartArray:
		openArray := 1
		for {
			if p.lexer.peekToken() == none {
				return value, nil
			}
			if p.lexer.peekToken().Kind == tk.StartArray {
				openArray++
			}

			token, err = p.lexer.nextToken()
			if err != nil {
				return nil, err
			}
			value = append(value, token)

			if token.Kind == tk.EndArray {
				openArray--
				if openArray == 0 {
					break
				}
			}
		}
	case tk.StartProc:
		proc, err := p.readProc()
		if err != nil {
			return nil, err
		}
		value = append(value, proc...)
	case tk.StartDic:
		// skip "/GlyphNames2HostCode << >> def"
		if _, err := p.read(tk.EndDic); err != nil {
			return nil, err
		}
		return value, nil
	}
	err = p.readPostScriptWrapper(value)
	return value, err
}

func (p *parser) readPostScriptWrapper(value []tk.Token) error {
	// postscript wrapper (not in the Type 1 spec)
	if p.lexer.peekToken().Value != "systemdict" {
		return nil
	}
	if err := p.readWithName(tk.Other, "systemdict"); err != nil {
		return err
	}
	if err := p.readWithName(tk.Name, "internaldict"); err != nil {
		return err
	}
	if err := p.readWithName(tk.Other, "known"); err != nil {
		return err
	}

	if _, err := p.read(tk.StartProc); err != nil {
		return err
	}
	if _, err := p.readProc(); err != nil {
		return err
	}

	if _, err := p.read(tk.StartProc); err != nil {
		return err
	}
	if _, err := p.readProc(); err != nil {
		return err
	}

	if err := p.readWithName(tk.Other, "ifelse"); err != nil {
		return err
	}

	// replace value
	if _, err := p.read(tk.StartProc); err != nil {
		return err
	}
	if err := p.readWithName(tk.Other, "pop"); err != nil {
		return err
	}
	value = nil
	other, err := p.readValue()
	if err != nil {
		return err
	}
	value = append(value, other...)
	if _, err := p.read(tk.EndProc); err != nil {
		return err
	}

	if err := p.readWithName(tk.Other, "if"); err != nil {
		return err
	}
	return nil
}

// Reads a procedure.
func (p *parser) readProc() ([]tk.Token, error) {
	var value []tk.Token
	openProc := 1
	for {
		if p.lexer.peekToken().Kind == tk.StartProc {
			openProc++
		}

		token, err := p.lexer.nextToken()
		if err != nil {
			return nil, err
		}
		value = append(value, token)

		if token.Kind == tk.EndProc {
			openProc--
			if openProc == 0 {
				break
			}
		}
	}
	executeonly, err := p.readMaybe(tk.Other, "executeonly")
	if err != nil {
		return nil, err
	}
	if executeonly != none {
		value = append(value, executeonly)
	}

	return value, nil
}

// Parses the binary portion of a Type 1 font.
func (p *parser) parseBinary(bytes []byte, font *PFBFont) error {
	var decrypted []byte
	// Sometimes, fonts use the hex format, so this needs to be converted before decryption
	if isBinary(bytes) {
		decrypted = decrypt(bytes, EEXEC_KEY, 4)
	} else {
		decrypted = decrypt(hexToBinary(bytes), EEXEC_KEY, 4)
	}

	p.lexer = newLexer(decrypted)

	// find /Private dict
	peekToken := p.lexer.peekToken()
	for peekToken.Value != "Private" {
		// for a more thorough validation, the presence of "begin" before Private
		// determines how code before and following charstrings should look
		// it is not currently checked anyway
		_, err := p.lexer.nextToken()
		if err != nil {
			return err
		}
		peekToken = p.lexer.peekToken()
	}
	if peekToken == none {
		return errors.New("/Private token not found")
	}

	// Private dict
	if err := p.readWithName(tk.Name, "Private"); err != nil {
		return err
	}
	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return err
	}
	length, _ := lengthT.Int()
	if err := p.readWithName(tk.Other, "dict"); err != nil {
		return err
	}
	// actually could also be "/Private 10 dict def Private begin"
	// instead of the "dup"
	if _, err := p.readMaybe(tk.Other, "dup"); err != nil {
		return err
	}
	if err := p.readWithName(tk.Other, "begin"); err != nil {
		return err
	}

	lenIV := 4 // number of random bytes at start of charstring

	for i := 0; i < length; i++ {
		// premature end
		if p.lexer.peekToken().Kind != tk.Name {
			break
		}

		// key/value
		key, err := p.read(tk.Name)
		if err != nil {
			return err
		}

		switch key.Value {
		case "Subrs":
			font.subrs, err = p.readSubrs(lenIV)
		case "OtherSubrs":
			err = p.readOtherSubrs()
		case "lenIV":
			vs, err := p.readDictValue()
			if err != nil {
				return err
			}
			lenIV, err = vs[0].Int()
		case "ND":
			if _, err := p.read(tk.StartProc); err != nil {
				return err
			}
			// the access restrictions are not mandatory
			if _, err := p.readMaybe(tk.Other, "noaccess"); err != nil {
				return err
			}
			if err := p.readWithName(tk.Other, "def"); err != nil {
				return err
			}
			if _, err := p.read(tk.EndProc); err != nil {
				return err
			}
			if _, err := p.readMaybe(tk.Other, "executeonly"); err != nil {
				return err
			}
			if err := p.readWithName(tk.Other, "def"); err != nil {
				return err
			}
		case "NP":
			if _, err := p.read(tk.StartProc); err != nil {
				return err
			}
			if _, err := p.readMaybe(tk.Other, "noaccess"); err != nil {
				return err
			}
			if _, err := p.read(tk.Other); err != nil {
				return err
			}
			if _, err := p.read(tk.EndProc); err != nil {
				return err
			}
			if _, err := p.readMaybe(tk.Other, "executeonly"); err != nil {
				return err
			}
			if err := p.readWithName(tk.Other, "def"); err != nil {
				return err
			}
		case "RD":
			// /RD {string currentfile exch readstring pop} bind executeonly def
			if _, err := p.read(tk.StartProc); err != nil {
				return err
			}
			if _, err := p.readProc(); err != nil {
				return err
			}
			if _, err := p.readMaybe(tk.Other, "bind"); err != nil {
				return err
			}
			if _, err := p.readMaybe(tk.Other, "executeonly"); err != nil {
				return err
			}
			if err := p.readWithName(tk.Other, "def"); err != nil {
				return err
			}
		default:
			var vs []tk.Token
			vs, err = p.readDictValue()
			if err != nil {
				return err
			}
			err = p.readPrivate(key.Value, vs)
		}

		if err != nil {
			return err
		}
	}

	// some fonts have "2 index" here, others have "end noaccess put"
	// sometimes followed by "put". Either way, we just skip until
	// the /CharStrings dict is found
	for p.lexer.peekToken() != (tk.Token{Kind: tk.Name, Value: "CharStrings"}) {
		_, err := p.lexer.nextToken()
		if err != nil {
			return err
		}
	}

	// CharStrings dict
	if err := p.readWithName(tk.Name, "CharStrings"); err != nil {
		return err
	}
	font.charstrings, err = p.readCharStrings(lenIV)
	return err
}

// Extracts values from the /Private dictionary.
func (p *parser) readPrivate(key string, value []tk.Token) error {
	// TODO: complete if needed
	// 		 switch (key)
	// 		 {
	// 			 case "BlueValues":
	// 				 font.blueValues = arrayToNumbers(value);
	// 				 break;
	// 			 case "OtherBlues":
	// 				 font.otherBlues = arrayToNumbers(value);
	// 				 break;
	// 			 case "FamilyBlues":
	// 				 font.familyBlues = arrayToNumbers(value);
	// 				 break;
	// 			 case "FamilyOtherBlues":
	// 				 font.familyOtherBlues = arrayToNumbers(value);
	// 				 break;
	// 			 case "BlueScale":
	// 				 font.blueScale = value[0].floatValue();
	// 				 break;
	// 			 case "BlueShift":
	// 				 font.blueShift = value[0].intValue();
	// 				 break;
	// 			 case "BlueFuzz":
	// 				 font.blueFuzz = value[0].intValue();
	// 				 break;
	// 			 case "StdHW":
	// 				 font.stdHW = arrayToNumbers(value);
	// 				 break;
	// 			 case "StdVW":
	// 				 font.stdVW = arrayToNumbers(value);
	// 				 break;
	// 			 case "StemSnapH":
	// 				 font.stemSnapH = arrayToNumbers(value);
	// 				 break;
	// 			 case "StemSnapV":
	// 				 font.stemSnapV = arrayToNumbers(value);
	// 				 break;
	// 			 case "ForceBold":
	// 				 font.forceBold = value[0].booleanValue();
	// 				 break;
	// 			 case "LanguageGroup":
	// 				 font.languageGroup = value[0].intValue();
	// 				 break;
	// 			 default:
	// 				 break;
	// 		 }
	return nil
}

// Reads the /Subrs array.
// `lenIV` is he number of random bytes used in charstring encryption.
func (p *parser) readSubrs(lenIV int) ([][]byte, error) {
	// allocate size (array indexes may not be in-order)
	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return nil, err
	}
	length, _ := lengthT.Int()
	subrs := make([][]byte, length)
	if err := p.readWithName(tk.Other, "array"); err != nil {
		return nil, err
	}

	for i := 0; i < length; i++ {
		// premature end
		if p.lexer.peekToken() != (tk.Token{Kind: tk.Other, Value: "dup"}) {
			break
		}

		if err := p.readWithName(tk.Other, "dup"); err != nil {
			return nil, err
		}
		indexT, err := p.read(tk.Integer)
		if err != nil {
			return nil, err
		}
		index, _ := indexT.Int()
		if _, err := p.read(tk.Integer); err != nil {
			return nil, err
		}
		if index >= length {
			return nil, fmt.Errorf("out of range charstring index %d (for %d)", index, length)
		}

		// RD
		charstring, err := p.read(tk.CharString)
		if err != nil {
			return nil, err
		}
		subrs[index] = decrypt([]byte(charstring.Value), CHARSTRING_KEY, lenIV)
		err = p.readPut()
		if err != nil {
			return nil, err
		}
	}
	err = p.readDef()
	return subrs, err
}

// OtherSubrs are embedded PostScript procedures which we can safely ignore
func (p *parser) readOtherSubrs() error {
	if p.lexer.peekToken().Kind == tk.StartArray {
		if _, err := p.readValue(); err != nil {
			return err
		}
		err := p.readDef()
		return err
	}
	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return err
	}
	length, _ := lengthT.Int()
	if err := p.readWithName(tk.Other, "array"); err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		if err := p.readWithName(tk.Other, "dup"); err != nil {
			return err
		}
		if _, err := p.read(tk.Integer); err != nil { // index
			return err
		}
		if _, err := p.readValue(); err != nil { // PostScript
			return err
		}
		if err := p.readPut(); err != nil {
			return err
		}
	}
	err = p.readDef()
	return err
}

// Reads the /CharStrings dictionary.
// `lenIV` is the number of random bytes used in charstring encryption.
func (p *parser) readCharStrings(lenIV int) ([]charstring, error) {
	lengthT, err := p.read(tk.Integer)
	if err != nil {
		return nil, err
	}
	length, _ := lengthT.Int()
	if err := p.readWithName(tk.Other, "dict"); err != nil {
		return nil, err
	}
	// could actually be a sequence ending in "CharStrings begin", too
	// instead of the "dup begin"
	if err := p.readWithName(tk.Other, "dup"); err != nil {
		return nil, err
	}
	if err := p.readWithName(tk.Other, "begin"); err != nil {
		return nil, err
	}

	charstrings := make([]charstring, length)
	for i := range charstrings {
		// premature end
		if tok := p.lexer.peekToken(); tok == none || tok == (tk.Token{Kind: tk.Other, Value: "end"}) {
			break
		}
		// key/value
		nameT, err := p.read(tk.Name)
		if err != nil {
			return nil, err
		}

		// RD
		_, err = p.read(tk.Integer)
		if err != nil {
			return nil, err
		}
		charstring, err := p.read(tk.CharString)
		if err != nil {
			return nil, err
		}

		charstrings[i].name = nameT.Value
		charstrings[i].data = decrypt([]byte(charstring.Value), CHARSTRING_KEY, lenIV)

		err = p.readDef()
		if err != nil {
			return nil, err
		}
	}

	// some fonts have one "end", others two
	err = p.readWithName(tk.Other, "end")
	// since checking ends here, this does not matter ....
	// more thorough checking would see whether there is "begin" before /Private
	// and expect a "def" somewhere, otherwise a "put"
	return charstrings, err
}

// Reads the sequence "noaccess def" or equivalent.
func (p *parser) readDef() error {
	if _, err := p.readMaybe(tk.Other, "readonly"); err != nil {
		return err
	}
	// allows "noaccess ND" (not in the Type 1 spec)
	if _, err := p.readMaybe(tk.Other, "noaccess"); err != nil {
		return err
	}

	token, err := p.read(tk.Other)
	if err != nil {
		return err
	}
	switch token.Value {
	case "ND", "|-":
		return nil
	case "noaccess":
		token, err = p.read(tk.Other)
		if err != nil {
			return err
		}
	}
	if token.Value == "def" {
		return nil
	}
	return fmt.Errorf("Found %s but expected ND", token.Value)
}

// Reads the sequence "noaccess put" or equivalent.
func (p *parser) readPut() error {
	_, err := p.readMaybe(tk.Other, "readonly")
	if err != nil {
		return err
	}

	token, err := p.read(tk.Other)
	if err != nil {
		return err
	}
	switch token.Value {
	case "NP", "|":
		return nil
	case "noaccess":
		token, err = p.read(tk.Other)
		if err != nil {
			return err
		}
	}

	if token.Value == "put" {
		return nil
	}
	return fmt.Errorf("found %s but expected NP", token.Value)
}

/// Reads the next token and throws an error if it is not of the given kind.
func (p *parser) read(kind tk.Kind) (tk.Token, error) {
	token, err := p.lexer.nextToken()
	if err != nil {
		return none, err
	}
	if token.Kind != kind {
		return none, fmt.Errorf("found token %s (%s) but expected token %s", token.Kind, token.Value, kind)
	}
	return token, nil
}

// Reads the next token and throws an error if it is not of the given kind
// and does not have the given value.
func (p *parser) readWithName(kind tk.Kind, name string) error {
	token, err := p.read(kind)
	if err != nil {
		return err
	}
	if token.Value != name {
		return fmt.Errorf("found %s but expected %s", token.Value, name)
	}
	return nil
}

// Reads the next token if and only if it is of the given kind and
// has the given value.
func (p *parser) readMaybe(kind tk.Kind, name string) (tk.Token, error) {
	token := p.lexer.peekToken()
	if token.Kind == kind && token.Value == name {
		return p.lexer.nextToken()
	}
	return none, nil
}

func decryptSegment(crypted []byte) ([]byte, error) {
	// Sometimes, fonts use the hex format, so this needs to be converted before decryption
	if isBinary(crypted) {
		return decrypt(crypted, EEXEC_KEY, 4), nil
	} else {
		dl := hex.DecodedLen(len(crypted))
		tmp := make([]byte, dl)
		_, err := hex.Decode(tmp, crypted)
		if err != nil {
			return nil, err
		}
		return decrypt(tmp, EEXEC_KEY, 4), nil
	}
}

// Type 1 Decryption (eexec, charstring).
// `r` is the key and `n` the number of random bytes (lenIV)
func decrypt(cipherBytes []byte, r, n int) []byte {
	// lenIV of -1 means no encryption (not documented)
	if n == -1 {
		return cipherBytes
	}
	// empty charstrings and charstrings of insufficient length
	if len(cipherBytes) == 0 || n > len(cipherBytes) {
		return nil
	}
	// decrypt
	c1 := 52845
	c2 := 22719
	plainBytes := make([]byte, len(cipherBytes)-n)
	for i := 0; i < len(cipherBytes); i++ {
		cipher := int(cipherBytes[i] & 0xFF)
		plain := int(cipher ^ r>>8)
		if i >= n {
			plainBytes[i-n] = byte(plain)
		}
		r = (cipher+r)*c1 + c2&0xffff
	}
	return plainBytes
}

// Check whether binary or hex encoded. See Adobe Type 1 Font Format specification
// 7.2 eexec encryption
func isBinary(bytes []byte) bool {
	if len(bytes) < 4 {
		return true
	}
	// "At least one of the first 4 ciphertext bytes must not be one of
	// the ASCII hexadecimal character codes (a code for 0-9, A-F, or a-f)."
	for i := 0; i < 4; i++ {
		by := bytes[i]

		if _, isHex := tk.IsHexChar(by); by != 0x0a && by != 0x0d && by != 0x20 && by != '\t' && !isHex {
			return true
		}
	}
	return false
}

func hexToBinary(data []byte) []byte {
	// white space characters may be interspersed
	tmp := make([]byte, 0, len(data))
	for _, c := range data {
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
			tmp = append(tmp, c)
		}
	}
	out := make([]byte, hex.DecodedLen(len(tmp)))
	hex.Decode(out, tmp)
	return out
}