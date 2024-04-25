package model

type CreateCAParam struct {
	Dn            string    `json:"dn"`
	Validity      int       `json:"validity"`
}

type CreateCAArgs struct {
	Path          string
	Dn            string
	MaxPathLength int
	Validity      int
	Overwrite     bool
}
