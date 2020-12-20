package hashtag_test

import (
	"testing"

	hashtag "github.com/13rac1/goldmark-hashtag"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestWikilink(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			hashtag.New(),
		),
	)

	testutil.DoTestCaseFile(markdown, "hashtags.txt", t)

}
