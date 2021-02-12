package truetype

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	filename := "testdata/Raleway-v4020-Regular.otf"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open %q: %s\n", filename, err)
	}

	font, err := Parse(file)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want nil", filename, err)
	}

	pos, err := font.GposTable()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := font.GsubTable()
	if err != nil {
		t.Fatal(err)
	}

	for _, table := range []TableLayout{pos.TableLayout, sub.TableLayout} {
		var tags []int
		for _, s := range table.Scripts {
			tags = append(tags, int(s.Tag))
		}
		if !sort.IntsAreSorted(tags) {
			t.Fatalf("tag not sorted: %v", tags)
		}
		for i, s := range table.Scripts {
			ptr := table.FindScript(s.Tag)
			if ptr != i {
				t.Errorf("binary search failed for script tag %s", s.Tag)
			}
		}

		s := table.FindScript(Tag(0)) // invalid
		if s != -1 {
			t.Errorf("binary search should have failed")
		}

		// now check the languages

		for _, script := range table.Scripts {
			var tags []int
			for _, s := range script.Languages {
				tags = append(tags, int(s.Tag))
			}
			if !sort.IntsAreSorted(tags) {
				t.Fatalf("tag not sorted: %v", tags)
			}
			for i, l := range script.Languages {
				ptr := script.FindLanguage(l.Tag)
				if ptr != i {
					t.Errorf("binary search failed for language tag %s", l.Tag)
				}
			}

			s := script.FindLanguage(Tag(0)) // invalid
			if s != -1 {
				t.Errorf("binary search should have failed")
			}
		}
	}
}

func TestFindSub(t *testing.T) {
	// dir := "/home/benoit/Téléchargements/harfbuzz/test/shaping/data/aots/fonts"
	dir := "/home/benoit/Téléchargements/harfbuzz/test/shaping/data/in-house/fonts"
	// dir := "/home/benoit/Téléchargements/harfbuzz/test/shaping/data/text-rendering-tests/fonts"
	// dir := "/home/benoit/Téléchargements/harfbuzz/test/api/fonts"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	// mainLoop:
	for _, fi := range files {
		file, err := os.Open(filepath.Join(dir, fi.Name()))
		if err != nil {
			t.Fatalf("Failed to open %q: %s\n", fi.Name(), err)
		}

		font, err := Parse(file)
		if err != nil {
			t.Logf("Parse(%q) err = %q, want nil", fi.Name(), err)
			continue
		}
		if font.tables[tagMort] != nil {
			fmt.Println("found mort:", fi.Name())
		}

		if font.tables[tagMorx] != nil {
			fmt.Println("found morx:", fi.Name())
		}
		// if font.tables[TagGsub] == nil {
		// 	continue
		// }
		// sub, err := font.GsubTable()
		// if err != nil {
		// 	t.Log(err)
		// 	continue
		// }
		// for _, l := range sub.Lookups {
		// 	for _, s := range l.Subtables {
		// 		if s.Data != nil && s.Data.Type() == SubReverse {
		// 			fmt.Println("found :", fi.Name())
		// 			continue mainLoop
		// 		}
		// 	}
		// }
	}
}

func dirFiles(t *testing.T, dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var filenames []string
	for _, fi := range files {
		filenames = append(filenames, filepath.Join(dir, fi.Name()))
	}
	return filenames
}

func TestGSUB(t *testing.T) {
	filenames := []string{
		"testdata/Raleway-v4020-Regular.otf",
		"testdata/Estedad-VF.ttf",
		"testdata/Mada-VF.ttf",
	}

	filenames = append(filenames, dirFiles(t, "testdata/layout_fonts/gsub")...)

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			t.Fatalf("Failed to open %q: %s\n", filename, err)
		}

		font, err := Parse(file)
		if err != nil {
			t.Fatalf("Parse(%q) err = %q, want nil", filename, err)
		}

		sub, err := font.GsubTable()
		if err != nil {
			t.Fatal(filename, err)
		}
		for _, l := range sub.Lookups {
			for _, s := range l.Subtables {
				if s.Data == nil {
					continue
				}
			}
		}
		fmt.Println(len(sub.Lookups), "lookups")
	}
}

func TestFeatureVariations(t *testing.T) {
	filename := "testdata/Commissioner-VF.ttf"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open %q: %s\n", filename, err)
	}

	font, err := Parse(file)
	if err != nil {
		t.Fatalf("Parse(%q) err = %q, want nil", filename, err)
	}

	gsub, err := font.GsubTable()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gsub.FeatureVariations)

	gdef, err := font.GDefTable()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gdef.Class)
}