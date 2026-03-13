package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListMembers(ctx context.Context, boardID string) ([]Member, error) {
	var members []Member
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/members", boardID), nil, &members)
	return members, err
}

func (c *Client) AddMemberToCard(ctx context.Context, cardID, memberID string) error {
	return c.Post(ctx, fmt.Sprintf("/1/cards/%s/idMembers", cardID), map[string]string{
		"value": memberID,
	}, nil)
}

func (c *Client) RemoveMemberFromCard(ctx context.Context, cardID, memberID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/idMembers/%s", cardID, memberID), nil)
}

func (c *Client) GetMe(ctx context.Context) (Member, error) {
	var member Member
	err := c.Get(ctx, "/1/members/me", nil, &member)
	return member, err
}
