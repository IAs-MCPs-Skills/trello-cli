package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListComments(ctx context.Context, cardID string) ([]Comment, error) {
	var comments []Comment
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/actions", cardID), map[string]string{
		"filter": "commentCard",
	}, &comments)
	return comments, err
}

func (c *Client) AddComment(ctx context.Context, cardID, text string) (Comment, error) {
	var comment Comment
	err := c.Post(ctx, fmt.Sprintf("/1/cards/%s/actions/comments", cardID), map[string]string{
		"text": text,
	}, &comment)
	return comment, err
}

func (c *Client) UpdateComment(ctx context.Context, actionID, text string) (Comment, error) {
	var comment Comment
	err := c.Put(ctx, fmt.Sprintf("/1/actions/%s/text", actionID), map[string]string{
		"value": text,
	}, &comment)
	return comment, err
}

func (c *Client) DeleteComment(ctx context.Context, actionID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/actions/%s", actionID), nil)
}
