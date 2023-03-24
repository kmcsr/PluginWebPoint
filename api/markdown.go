
package api

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark-emoji"
	emojiDef "github.com/yuin/goldmark-emoji/definition"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/mermaid"
	"github.com/microcosm-cc/bluemonday"
)

type Option struct{
	URLPrefix string
	DataURLPrefix string
	HeadingIDPrefix string
	HeadingIDSuffix string
}

func (o *Option)fixRelLink(prefix, src string)(string){
	if !strings.Contains(src, "://") && !strings.HasPrefix(src, "/") && !strings.HasPrefix(src, "#") {
		out, err := url.JoinPath(prefix, src)
		if err != nil {
			out = "about:blank#err:" + err.Error()
		}
		return out
	}
	return src
}

func (o *Option)Transform(doc *ast.Document, reader text.Reader, ctx parser.Context){
	ast.Walk(doc, func(node ast.Node, enter bool)(ast.WalkStatus, error){
		if !enter {
			return ast.WalkContinue, nil
		}
		switch m := node.(type) {
		case *ast.Link:
			if len(o.URLPrefix) > 0 {
				m.Destination = ([]byte)(o.fixRelLink(o.URLPrefix, (string)(m.Destination)))
			}
		case *ast.Image:
			if len(o.DataURLPrefix) > 0 {
				m.Destination = ([]byte)(o.fixRelLink(o.DataURLPrefix, (string)(m.Destination)))
			}else if len(o.URLPrefix) > 0 {
				m.Destination = ([]byte)(o.fixRelLink(o.URLPrefix, (string)(m.Destination)))
			}
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
}

var langClassRe = regexp.MustCompile("language-[^\"' ]")

func (o *Option)NewPolicy()(p *bluemonday.Policy){
	p = bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(langClassRe).OnElements("code")
	p.AllowAttrs("class").Matching(regexp.MustCompile("anchor")).OnElements("a")
	p.AllowAttrs("class").Matching(regexp.MustCompile("mermaid")).OnElements("div", "pre")
	return
}

var DefaultOption = &Option{}

func RenderMarkdown(src []byte, opt *Option)(out []byte, err error){
	buf := bytes.NewBuffer(make([]byte, 0, len(src) * 3 / 2))
	renderer := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
		),
		goldmark.WithExtensions(
			emoji.Emoji,
		),
		goldmark.WithExtensions(
			highlighting.Highlighting,
		),
		goldmark.WithExtensions(
			&anchor.Extender{Texter: anchor.Text("#")},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(opt, 100)),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
	// src = bytes.ReplaceAll(src, ([]byte)("\r\n"), ([]byte)("\n"))
	if err = renderer.Convert(src, buf); err != nil {
		return
	}
	out = opt.NewPolicy().SanitizeReader(buf).Bytes()
	return
}

var emojiSet = emojiDef.Github()

func ReplaceEmoji(src []byte)(out []byte){
	out = make([]byte, len(src))
	s := 0
	l := 0
	for s < len(src) {
		if i := bytes.IndexByte(src[s:], ':'); i != -1 {
			i += s
			l += copy(out[l:], src[s:i])
			s = i
			if i := bytes.IndexByte(src[s + 1:], ':'); i != -1 {
				s++
				i += s
				emo, ok := emojiSet.Get((string)(src[s:i]))
				if ok {
					l += copy(out[l:], ([]byte)((string)(emo.Unicode)))
				}else{
					l += copy(out[l:], src[s:i + 1])
				}
				s = i + 1
				continue
			}
		}
		l += copy(out[l:], src[s:])
		break
	}
	return out[:l]
}
