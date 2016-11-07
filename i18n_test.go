package i18n

import "testing"

func TestTranslate(t *testing.T) {
	if err := Load("translations"); err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	var tests = []struct {
		value    string
		language string
		expected string
	}{
		{"Hello World", "", "Hello World"},
		{"Hello World", "it", "Ciao Mondo"},
		{"Hello World", "es", "Hola Mundo"},
		{"Hello Gopher", "", "Hello Gopher"},
	}
	for _, test := range tests {
		i18n := New(test.language)
		res := i18n.Print(test.value)
		if res != test.expected {
			t.Errorf("%s, expected %v, got %v", test.value, test.expected, res)
		}
	}
}

func TestPlural(t *testing.T) {
	if err := Load("translations"); err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	var tests = []struct {
		value    int
		language string
		zero     string
		one      string
		many     string
		expected string
	}{
		{-1, "", "no records found.", "one record found.", "%d records found.", "no records found."},
		{0, "it", "no records found.", "one record found.", "%d records found.", "no records found."},
		{1, "es", "no records found.", "one record found.", "%d records found.", "one record found."},
		{2, "en", "no records found.", "one record found.", "%d records found.", "2 records found."},
	}
	for _, test := range tests {
		i18n := New(test.language)
		res := i18n.Plural(test.value, test.zero, test.one, test.many)
		if res != test.expected {
			t.Errorf("%d, expected %v, got %v", test.value, test.expected, res)
		}
	}
}

func BenchmarkPrint(b *testing.B) {
	b.ReportAllocs()

	if err := Load("translations"); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		i18n := New("it")
		i18n.Print("Hello world")
	}
}

func BenchmarkPrintf(b *testing.B) {
	b.ReportAllocs()

	if err := Load("translations"); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		i18n := New("it")
		i18n.Print("Hello world %d", 20)
	}
}
