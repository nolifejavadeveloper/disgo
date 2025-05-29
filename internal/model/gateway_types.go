package model

type Snowflake int64

type Activity struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Url   string `json:"url,omitempty"`
	State string `json:"state,omitempty"`
}

type User struct {
	ID                   Snowflake         `json:"id"`
	Username             string            `json:"username"`
	Discriminator        string            `json:"discriminator"`
	GlobalName           *string           `json:"global_name,omitempty"`
	Avatar               *string           `json:"avatar,omitempty"`
	Bot                  bool              `json:"bot"`
	System               *bool             `json:"system,omitempty"`
	MFAEnabled           bool              `json:"mfa_enabled"`
	Banner               *string           `json:"banner,omitempty"`
	AccentColor          *int              `json:"accent_color,omitempty"`
	Locale               string            `json:"locale"`
	Verified             *bool             `json:"verified,omitempty"`
	Email                *string           `json:"email,omitempty"`
	Flags                int               `json:"flags"`
	PremiumType          int               `json:"premium_type"`
	PublicFlags          int               `json:"public_flags"`
	AvatarDecorationData *AvatarDecoration `json:"avatar_decoration_data,omitempty"`
}

type AvatarDecoration struct {
	Asset string    `json:"asset"`
	SkuId Snowflake `json:"sku_id"`
}
