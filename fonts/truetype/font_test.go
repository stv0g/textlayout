package truetype

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	for _, filename := range []string{
		"testdata/Roboto-BoldItalic.ttf",
		"testdata/Raleway-v4020-Regular.otf",
		"testdata/open-sans-v15-latin-regular.woff",
	} {
		file, err := os.Open(filename)
		if err != nil {
			t.Fatalf("Failed to open %q: %s\n", filename, err)
		}

		font, err := Parse(file)
		if err != nil {
			t.Fatalf("Parse(%q) err = %q, want nil", filename, err)
		}
		for tag := range font.tables {
			_ = tag.String()
		}

		_, err = font.OS2Table()
		if err != nil {
			t.Fatal(err)
		}
		_, err = font.GposTable()
		if err != nil {
			t.Fatal(err)
		}
		_, err = font.GsubTable()
		if err != nil {
			t.Fatal(err)
		}
		_, err = font.HeadTable()
		if err != nil {
			t.Fatal(err)
		}
		_, err = font.HheaTable()
		if err != nil {
			t.Fatal(err)
		}
		_, err = font.NameTable()
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}
}

func TestParseCrashers(t *testing.T) {
	font, err := Parse(bytes.NewReader([]byte{}))
	if font != nil || err == nil {
		t.Fail()
	}

	for range [50]int{} {
		input := make([]byte, 20000)
		rand.Read(input)
		font, err = Parse(bytes.NewReader(input))
		if font != nil || err == nil {
			t.Error("expected error on random input")
		}
	}
}

func TestTables(t *testing.T) {
	f, err := os.Open("testdata/LateefGR-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	font, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse err = %q, want nil", err)
	}
	fmt.Println(font.tables)
}