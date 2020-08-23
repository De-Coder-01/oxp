package oxp

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Parser struct {
	content *html.Node
	entry   *Entry
}

func queryInnerText(node *html.Node, expr string) (string, error) {
	node, err := htmlquery.Query(node, expr)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", nil
	}
	return htmlquery.InnerText(node), nil
}

func (p *Parser) Parse(content *html.Node) (*Entry, error) {
	var err error

	p.content = content
	p.entry = &Entry{}

	p.entry.Headword, err = queryInnerText(content, "//h1[@class='headword']")
	if err != nil {
		return nil, err
	}
	p.entry.POS, err = queryInnerText(content, "//*[@class='pos']")
	if err != nil {
		return nil, err
	}
	if err = p.parseSenses(); err != nil {
		return nil, err
	}
	return p.entry, nil
}

func (p *Parser) parseSense(node *html.Node) (sense Sense, err error) {
	sense.Shcut, err = queryInnerText(node, "../*[@class='shcut']")
	if err != nil {
		return
	}
	sense.Def, err = queryInnerText(node, "//span[@class='def']")
	if err != nil {
		return
	}
	err = p.parseExamples(&sense, node)
	return
}

func (p *Parser) parseSensesIn(node *html.Node) error {
	nodes, err := htmlquery.QueryAll(node, "//li[@class='sense']")
	if err != nil {
		return err
	}
	for _, node := range nodes {
		sense, err := p.parseSense(node)
		if err != nil {
			return err
		}
		p.entry.Senses = append(p.entry.Senses, sense)
	}
	return nil
}

func (p *Parser) parseSenses() error {
	single, err := htmlquery.Query(p.content, "//div[@class='entry']/ol[@class='sense_single']")
	if err != nil {
		return err
	}
	if single != nil {
		return p.parseSensesIn(single)
	}

	multiple, err := htmlquery.Query(p.content, "//div[@class='entry']/ol[@class='senses_multiple']")
	if err != nil {
		return err
	}
	if multiple != nil {
		return p.parseSensesIn(multiple)
	}
	return nil
}

func (p *Parser) parseExamples(sense *Sense, node *html.Node) error {
	nodes, err := htmlquery.QueryAll(node, "./ul[@class='examples']/li/span[@class='x']")
	if err != nil {
		return err
	}
	for _, node := range nodes {
		example := htmlquery.InnerText(node)
		if example != "" {
			sense.Examples = append(sense.Examples, example)
		}
	}
	return nil
}
