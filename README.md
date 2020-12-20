# goldmark-hashtag

goldmark-hashtag is an extension for the [goldmark][goldmark] library to handle
`#example` hashtags providing a Hashtag AST type.

[goldmark]: http://github.com/yuin/goldmark

## Demo

This markdown:

```md
# Hello goldmark-hashtag

#example
```

With the default configuration, becomes this HTML:

```html
<h1>Hello goldmark-hashtag</h1>
<p><a href="/tags/example">#example</a></p>
```

### Installation

```bash
go get github.com/13rac1/goldmark-hashtag
```

## Usage

```go
  markdown := goldmark.New(
    goldmark.WithExtensions(
      hashtag.Extension,
    ),
  )
  var buf bytes.Buffer
  if err := markdown.Convert([]byte(source), &buf); err != nil {
    panic(err)
  }
  fmt.Print(buf)
}
```

## License

MIT

## Author

Brad Erickson
