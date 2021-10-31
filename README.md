# esacli
esacli is command line tool for esa.

# Usage
```sh
$ esacli -team warashi -file path/to/post.md
```

team: ex) warashi for warashi.esa.io

# Example file
this example posted as "foo/bar/category/Title", and its body is `# Hello, world!`
```markdwon:post.md
---
title: "Title"
category: foo/bar/category
tags: tag1,tag2
published: true
number: 123
---
# Hello, world!
```
