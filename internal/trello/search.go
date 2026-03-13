package trello

import "context"

func (c *Client) SearchCards(ctx context.Context, query string) (CardSearchResult, error) {
	var raw struct {
		Cards []Card `json:"cards"`
	}
	err := c.Get(ctx, "/1/search", map[string]string{
		"query":      query,
		"modelTypes": "cards",
	}, &raw)
	return CardSearchResult{Query: query, Cards: raw.Cards}, err
}

func (c *Client) SearchBoards(ctx context.Context, query string) (BoardSearchResult, error) {
	var raw struct {
		Boards []Board `json:"boards"`
	}
	err := c.Get(ctx, "/1/search", map[string]string{
		"query":      query,
		"modelTypes": "boards",
	}, &raw)
	return BoardSearchResult{Query: query, Boards: raw.Boards}, err
}
