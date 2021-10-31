package main

import (
	"fmt"
	"io"

	"github.com/Warashi/esa-action/frontmatter"
	"github.com/goccy/go-yaml"
)

type Meta struct {
	Title     string `yaml:"title"`
	Category  string `yaml:"category"`
	Tags      string `yaml:"tags"`
	Published bool   `yaml:"published"`
	Number    uint    `yaml:"number"`
}

func ParseMeta(r io.Reader) (*Meta, error) {
	b, err := frontmatter.Extract(r, "---")
	if err != nil {
		return nil, fmt.Errorf("frontmatter.Extract: %w", err)
	}
	var meta Meta
	if err := yaml.Unmarshal(b, &meta); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}
	return &meta, nil
}

func ParseFile(r io.Reader) (body string, meta *Meta, err error) {
	meta, err = ParseMeta(r)
	if err != nil {
		return "", nil, fmt.Errorf("ParseMeta: %w", err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return "", nil,  fmt.Errorf("io.ReadAll: %w", err)
	}
	return string(b), meta, nil
}
