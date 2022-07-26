package oxp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Client struct {
	url       string
	userAgent string

	http http.Client
}

func NewClient() *Client {
	c := &Client{
		url:       "https://www.oxfordlearnersdictionaries.com/search/english/?q=@",
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.136 YaBrowser/20.2.3.320 (beta) Yowser/2.5 Safari/537.36",
		http:      http.Client{},
	}
	return c
}

func (c *Client) getPage(ctx context.Context, url string,
) (*html.Node, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header["User-Agent"] = []string{c.userAgent}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, resp, errors.Errorf("status %s", resp.Status)
	}
	page, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return page, resp, nil
}

func (c *Client) parsePage(page *html.Node) (*Entry, error) {
	content, err := c.extractContent(page)
	if err != nil {
		return nil, err
	}
	parser := Parser{}
	return parser.Parse(content)
}

func (c *Client) extractURLs(page *html.Node, url url.URL) ([]string, error) {
	url.RawQuery = ""
	path := strings.TrimSuffix(url.Path, "_1")
	if path == url.Path {
		return nil, nil
	}

	urls := []string{}
	for i := 2; true; i++ {
		url.Path = fmt.Sprintf("%s_%d", path, i)
		filter := fmt.Sprintf("//*[@id='relatedentries']//a[@href='%s']", &url)
		el, err := htmlquery.Query(page, filter)
		if err != nil {
			return nil, err
		}
		if el == nil {
			break
		}
		urls = append(urls, url.String())
	}
	return urls, nil
}

func (c *Client) extractContent(page *html.Node) (*html.Node, error) {
	content, err := htmlquery.Query(page, "//div[@id='entryContent']")
	if err != nil {
		return nil, err
	}
	if content == nil {
		return nil, errors.New("entry content element not found")
	}
	return content, nil
}

func (c *Client) getPagesWithSameHeadwordInURL(
	ctx context.Context, page *html.Node, url *url.URL,
) ([]*html.Node, error) {
	urls, err := c.extractURLs(page, *url)
	if err != nil {
		return nil, err
	}

	pages := []*html.Node{page}
	for _, url := range urls {
		if page, _, err = c.getPage(ctx, url); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (c *Client) Search(ctx context.Context, text string) (interface{}, error) {
	url := strings.Replace(c.url, "@", text, -1)

	page, resp, err := c.getPage(ctx, url)
	if err != nil {
		return nil, err
	}
	if _, err = c.extractContent(page); err != nil {
		return nil, err
	}

	pages, err := c.getPagesWithSameHeadwordInURL(ctx, page, resp.Request.URL)
	if err != nil {
		return nil, err
	}

	entries := []*Entry{}
	for _, page := range pages {
		entry, err := c.parsePage(page)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
