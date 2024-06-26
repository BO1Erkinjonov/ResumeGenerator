package entity

// CREATE TABLE IF NOT EXISTS links (
// id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// resume_id UUID REFERENCES resumes(id),
// link_name VARCHAR(200),
// link_url TEXT
// );
type Link struct {
	ID       string `json:"id"`
	ResumeID string `json:"resume_id"`
	LinkName string `json:"link_name"`
	LinkURL  string `json:"link_url"`
}

type LinksUpdateReq struct {
	LinkName string `json:"link_name"`
	LinkURL  string `json:"link_url"`
}
