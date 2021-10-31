package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWrite(t *testing.T) {
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
	var buf bytes.Buffer
	if err := write(&buf, body, meta); err != nil {
		panic(err)
	}

	expected := `---
title: title
category: foo/bar/hoge
tags: tag1,tag2
published: false
number: 10
---
# hogehoge

aaa
bbb

# foobar
foo
bar
`
	if actual := buf.String(); actual != expected {
		t.Log(actual)
		t.Error(cmp.Diff(expected, actual))
	}
}
