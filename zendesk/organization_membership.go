package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type OrganizationMembership struct {
	ID             int64     `json:"id,omitempty"`
	URL            string    `json:"url,omitempty"`
	UserID         int64     `json:"user_id,omitempty"`
	OrganizationID int64     `json:"organization_id,omitempty"`
	Default        bool      `json:"default"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type OrganizationMembershipListOptions struct {
	PageOptions
}

type OrganizationMembershipAPI interface {
	GetOrganizationMemberships(ctx context.Context, opts *OrganizationMembershipListOptions) ([]OrganizationMembership, Page, error)
	DeleteOrganizationMembership(ctx context.Context, orgMemID int64) error
}

func (z *Client) GetOrganizationMemberships(ctx context.Context, userID int64) ([]OrganizationMembership, Page, error) {
	var data struct {
		OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
		Page
	}

	body, err := z.get(ctx, fmt.Sprintf("/users/%d/organization_memberships.json", userID))
	if err != nil {
		return []OrganizationMembership{}, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []OrganizationMembership{}, Page{}, err
	}

	return data.OrganizationMemberships, data.Page, nil
}
func (z *Client) DeleteOrganizationMembership(ctx context.Context, orgMemID int64) error {
	err := z.delete(ctx, fmt.Sprintf("/organization_memberships/%d.json", orgMemID))

	if err != nil {
		return err
	}

	return nil
}

// CreateOrganizationMembership はメンバーシップ(組織-ユーザ紐付け情報)を1件作成します
// ref: https://developer.zendesk.com/api-reference/ticketing/organizations/organization_memberships/#create-membership
func (z *Client) CreateOrganizationMembership(ctx context.Context, userID, organizationID int64) (OrganizationMembership, error) {
	var data, result struct {
		OrganizationMembership OrganizationMembership `json:"organization_membership"`
	}
	data.OrganizationMembership.UserID = userID
	data.OrganizationMembership.OrganizationID = organizationID

	body, err := z.post(ctx, "/organization_memberships.json", data)
	if err != nil {
		return OrganizationMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OrganizationMembership{}, err
	}

	return result.OrganizationMembership, nil
}
