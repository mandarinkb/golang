package assets

type RootError struct {
	Success ErrorCode `mapstructure:"success"`
}
type ErrorCode struct {
	StatusCode           int    `mapstructure:"status" json:"statusCode"`
	MessageCode          string `mapstructure:"code" json:"messageCode,omitempty"`
	MessageDescription   string `mapstructure:"en" json:"messageDescription,omitempty"`
	MessageDescriptionTh string `mapstructure:"th" json:"messageDescriptionTH,omitempty"`
}
