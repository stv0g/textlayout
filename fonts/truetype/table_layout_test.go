package truetype

import (
	"fmt"
	"os"
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

// func TestFindSub(t *testing.T) {
// 	dir := "/home/benoit/Téléchargements/harfbuzz/test/api/fonts"
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// mainLoop:
// 	for _, fi := range files {
// 		file, err := os.Open(filepath.Join(dir, fi.Name()))
// 		if err != nil {
// 			t.Fatalf("Failed to open %q: %s\n", fi.Name(), err)
// 		}

// 		font, err := Parse(file)
// 		if err != nil {
// 			t.Fatalf("Parse(%q) err = %q, want nil", fi.Name(), err)
// 		}

// 		sub, err := font.GsubTable()
// 		if err != nil {
// 			continue
// 		}
// 		for _, l := range sub.Lookups {
// 			for _, s := range l.Subtables {
// 				if s.Data != nil && s.Data.Type() == SubMultiple {
// 					fmt.Println("found :", fi.Name())
// 					continue mainLoop
// 				}
// 			}
// 		}
// 	}
// }

func TestGSUB(t *testing.T) {
	filenames := [...]string{
		"testdata/Raleway-v4020-Regular.otf",
		"testdata/Estedad-VF.ttf",
		"testdata/Mada-VF.ttf",
	}
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
			t.Fatal(err)
		}
		for _, l := range sub.Lookups {
			for _, s := range l.Subtables {
				if s.Data == nil {
					continue
				}
				if s.Data.Type() == SubMultiple {
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
