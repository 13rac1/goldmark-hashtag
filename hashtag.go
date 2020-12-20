package hashtag

import (
	"fmt"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Option is a functional option type for this extension.
type Option func(*htExtension)

type htExtension struct{}

// New returns a new Hashtag extension.
func New(opts ...Option) goldmark.Extender {
	e := &htExtension{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Extend adds a hashtag parser to a Goldmark parser
func (e *htExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(util.Prioritized(
			NewParser(), 120),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
}

// Hashtag struct represents a Hashtag of the Markdown text.
type Hashtag struct {
	ast.Link
}

// KindHashtag is a NodeKind of the Hashtag node.
var KindHashtag = ast.NewNodeKind("Hashtag")

// Kind implements Node.Kind.
func (n *Hashtag) Kind() ast.NodeKind {
	return KindHashtag
}

// NewHashtag returns a new Hashtag node.
func NewHashtag(title, destination []byte) *Hashtag {
	link := ast.NewLink()
	link.Title = title
	link.Destination = destination

	c := &Hashtag{
		Link: *link,
	}

	return c
}

// Linker describes the method to convert hashtags into URLs.
type Linker interface {
	Link(text string) string
}

type defaultLinker struct{}

func (t *defaultLinker) Link(text string) string {
	return "/tags/" + text[1:]
}

// Parser implements the goldmark parser.Parser interface
type Parser struct {
	linker Linker
}

// NewParser gives you back a parser that you can use to process hashtags.
func NewParser() *Parser {
	return &Parser{
		linker: &defaultLinker{},
	}
}

// WithNormalizer is the fluent interface for replacing the default normalizer plugin.
func (p *Parser) WithNormalizer(fn Linker) *Parser {
	p.linker = fn
	return p
}

// Trigger looks for the # beginning of a hashtag
func (p *Parser) Trigger() []byte {
	return []byte{'#'}
}

// Parse implements the parser.Parser interface for Hashtags in markdown
func (p *Parser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// TODO: Confirm all unicode work, may need to switch to runes instead of bytes.

	line, segment := block.PeekLine()

	if line[1] == ' ' {
		// Must not be a header
		return nil
	}
	pos := 2

	for ; pos < len(line); pos++ {
		if isHashtagEnd(rune(line[pos])) {
			break
		}
	}

	title := block.Value(text.NewSegment(segment.Start, segment.Start+pos))

	block.Advance(pos)
	destination := p.linker.Link(string(title))
	return NewHashtag(title, []byte(destination))
}

func isMarkdownSpecial(char rune) bool {
	return char == '#' || char == '_' || char == '*' || char == '`'
}

func isHashtagEnd(char rune) bool {
	return unicode.IsSpace(char) || isMarkdownSpecial(char) || char == '.'
}

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct{}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer() renderer.NodeRenderer {
	return &HTMLRenderer{}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHashtag, r.render)
}

func (r *HTMLRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkContinue, nil
	}

	ht := node.(*Hashtag)
	out := fmt.Sprintf(`<a href="%s">%s</a>`, ht.Destination, ht.Title)
	w.Write([]byte(out))
	return ast.WalkContinue, nil
}
