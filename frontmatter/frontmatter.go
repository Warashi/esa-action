package frontmatter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func Extract(r io.Reader, sep string) ([]byte, error) {
	var buf bytes.Buffer
	s := bufio.NewScanner(r)
	if ok := s.Scan(); !ok {
		return nil, fmt.Errorf("s.Scan: %w", s.Err())
	}
	if s.Text() != sep {
		return nil, fmt.Errorf("frontmatter does not start with %q", sep)
	}
	for s.Scan() {
		if s.Text() == sep {
			return buf.Bytes(), nil
		}
		if _, err := buf.WriteString(s.Text()); err != nil {
			return nil, fmt.Errorf("buf.WriteString: %w", err)
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("s.Err: %w", err)
	}

	return nil, fmt.Errorf("frontmatter ends without %q", sep)
}
