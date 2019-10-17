package markdown

import (
	"github.com/niklasfasching/go-org/org"
)

func (r *Renderer) Before(doc *org.Document) {
	r.DocumentHeader(r.out)
}

func (r *Renderer) After(doc *org.Document) {
	r.DocumentFooter(r.out)
}

func (r *Renderer) String() string {
	return r.out.String()
}

func (r *Renderer) WriteKeyword(item org.Keyword) {
}

func (r *Renderer) WriteInclude(item org.Include) {
}

func (r *Renderer) WriteComment(item org.Comment) {
}

func (r *Renderer) WriteNodeWithMeta(item org.NodeWithMeta) {
}

func (r *Renderer) WriteHeadline(item org.Headline) {
	r.Header(r.out, func() bool {
		org.WriteNodes(r, item.Title...)
		return true
	}, item.Lvl, item.ID())
}

func (r *Renderer) WriteBlock(item org.Block) {
}

func (r *Renderer) WriteExample(item org.Example) {
	r.BlockCode(r.out, []byte(item.String()), "")
}

func (r *Renderer) WriteDrawer(item org.Drawer) {
}

func (r *Renderer) WritePropertyDrawer(item org.PropertyDrawer) {
}

func (r *Renderer) WriteList(item org.List) {
}

func (r *Renderer) WriteListItem(item org.ListItem) {
	r.ListItem(r.out, []byte(item.String()), 0)
}

func (r *Renderer) WriteDescriptiveListItem(item org.DescriptiveListItem) {
}

func (r *Renderer) WriteTable(item org.Table) {
}

func (r *Renderer) WriteHorizontalRule(item org.HorizontalRule) {
}

func (r *Renderer) WriteParagraph(item org.Paragraph) {
	r.Paragraph(r.out, func() bool {
		org.WriteNodes(r, item.Children...)
		return true
	})
}

func (r *Renderer) WriteText(item org.Text) {
	r.NormalText(r.out, []byte(item.Content))
}

func (r *Renderer) WriteEmphasis(item org.Emphasis) {
}

func (r *Renderer) WriteLatexFragment(item org.LatexFragment) {
}

func (r *Renderer) WriteStatisticToken(item org.StatisticToken) {
}

func (r *Renderer) WriteExplicitLineBreak(item org.ExplicitLineBreak) {
}

func (r *Renderer) WriteLineBreak(item org.LineBreak) {
	for i := 0; i < item.Count; i++ {
		r.LineBreak(r.out)
	}
}

func (r *Renderer) WriteRegularLink(item org.RegularLink) {
	r.Link(r.out, []byte(item.URL), nil, []byte(item.String()))
}

func (r *Renderer) WriteTimestamp(item org.Timestamp) {
}

func (r *Renderer) WriteFootnoteLink(item org.FootnoteLink) {
}

func (r *Renderer) WriteFootnoteDefinition(item org.FootnoteDefinition) {
}
