package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hori-ryota/esa-go/esa"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var teamName, filename string
	flag.StringVar(&teamName, "team", "", "team name. ex) warashi for warashi.esa.io")
	flag.StringVar(&filename, "file", "", "filename")
	flag.Parse()

	if teamName == "" || filename == "" {
		flag.Usage()
		return 1
	}

	token := os.Getenv("ESA_API_TOKEN")
	if token == "" {
		fmt.Println("set api token as environment variable named ESA_API_TOKEN")
		return 1
	}

	client := esa.NewClient(token, teamName)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	body, meta, err := parseFile(filename)
	if err != nil {
		log.Println(err)
		return 1
	}

	if meta.Number == 0 {
		meta, err := create(ctx, client, body, meta)
		if err != nil {
			log.Println(err)
			return 1
		}
		if err :=  writeFile(filename, body, meta); err!= nil {
			log.Println(err)
			return 1
		}
	}
	meta, err = update(ctx, client, body, meta)
	if err != nil {
		log.Println(err)
		return 1
	}
	if err := writeFile(filename, body, meta); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func parseFile(filename string) (body string, meta *Meta, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", nil, fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()
	return ParseFile(f)
}

func writeFile(filename string, body string, meta *Meta) error {
	f, err := os.CreateTemp("", "")
	tmpName := f.Name()
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	if err := write(f, body, meta); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}
	if err := os.Rename(tmpName, filename); err != nil {
		return fmt.Errorf("os.Rename: %w", err)
	}
	return nil
}

func write(w interface {
	io.Writer
	io.StringWriter
}, body string, meta *Meta) error {
	if _, err := w.WriteString("---\n"); err != nil {
		return fmt.Errorf("w.WriteString: %w", err)
	}
	b, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("yaml.Marshal: %w", err)
	}
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("w.Write: %w", err)
	}
	if _, err := w.WriteString("---\n"); err != nil {
		return fmt.Errorf("w.WriteString: %w", err)
	}
	if _, err := w.WriteString(body); err != nil {
		return fmt.Errorf("w.WriteString: %w", err)
	}
	return nil
}

func create(ctx context.Context, client esa.Client, body string, meta *Meta) (*Meta, error) {
	tags := strings.Split(meta.Tags, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	p := esa.CreatePostParam{
		Name:     meta.Title,
		BodyMD:   &body,
		Tags:     &tags,
		Category: &meta.Category,
		WIP:      esa.BoolP(!meta.Published),
	}
	post, err := client.CreatePost(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("client.CreatePost: %w", err)
	}
	return &Meta{
		Title:     post.Name,
		Category:  post.Category,
		Tags:      strings.Join(post.Tags, ","),
		Published: !post.WIP,
		Number:    post.Number,
	}, nil
}

func update(ctx context.Context, client esa.Client, body string, meta *Meta) (*Meta, error) {
	tags := strings.Split(meta.Tags, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	p := esa.UpdatePostParam{
		Name:     &meta.Title,
		BodyMD:   &body,
		Tags:     &tags,
		Category: &meta.Category,
		WIP:      esa.BoolP(!meta.Published),
	}
	post, err := client.UpdatePost(ctx, meta.Number, p)
	if err != nil {
		return nil, fmt.Errorf("client.CreatePost: %w", err)
	}
	return &Meta{
		Title:     post.Name,
		Category:  post.Category,
		Tags:      strings.Join(post.Tags, ","),
		Published: !post.WIP,
		Number:    post.Number,
	}, nil

}
