package internal

import "testing"

func TestStoreCreateWorkflowAndRun(t *testing.T) {
	store := NewStore()
	workflow, err := store.CreateWorkflow(WorkflowTemplate{
		Name:               "Space Adventure",
		Description:        "Generate safe science stories",
		ContentSuitability: "core",
		AgeBand:            "6-11",
		Steps:              []string{"script", "scene", "render"},
		ModelProfileID:     "nim-default",
		SafetyProfile:      "strict",
	}, "admin-1")
	if err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if workflow.ID == "" || workflow.Version != 1 {
		t.Fatalf("unexpected workflow fields: %+v", workflow)
	}
	run, err := store.CreateRun(WorkflowRun{WorkflowID: workflow.ID, Priority: "normal", AutoPublish: false}, "admin-1")
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	if run.ID == "" || run.Status != "requested" {
		t.Fatalf("unexpected run fields: %+v", run)
	}
}
