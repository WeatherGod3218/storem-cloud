package models

type GetUserInfoResponse struct {
	Role    string `json:"role"`
	Actions bool   `json:"actions"`
}
type UserToken struct {
	Validated bool    `json:"validated"`
	SubjectId *string `json:"sub"`
	Email     *string `json:"email"`
}

type UserConfig struct {
	UserID       string   `yaml:"userId"`
	Emails       []string `yaml:"emails"`
	DisplayName  string   `yaml:"displayName"`
	TotalStorage int      `yaml:"totalStorageGB"`
	Role         string   `yaml:"role"`
}
