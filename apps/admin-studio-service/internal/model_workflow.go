package internal

type WorkflowTemplate struct {
	ID                 string
	Name               string
	Description        string
	ContentSuitability string
	AgeBand            string
	Steps              []string
	ModelProfileID     string
	SafetyProfile      string
	Version            int
}
