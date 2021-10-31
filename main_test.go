package main

import (
	"os"
)

func ExampleWrite() {
	body := `# hogehoge

aaa
bbb

# foobar
foo
bar
`
	meta := &Meta{
		Title:     "title",
		Category:  "foo/bar/hoge",
		Tags:      "tag1,tag2",
		Published: false,
		Number:    10,
	}

	if err := write(os.Stdout, body, meta); err != nil {
		panic(err)
	}

	// output:
	// ---
	// title: title
	// category: foo/bar/hoge
	// tags: tag1,tag2
	// published: false
	// number: 10
	// ---
	// # hogehoge
	//
	// aaa
	// bbb
	//
	// # foobar
	// foo
	// bar
}
