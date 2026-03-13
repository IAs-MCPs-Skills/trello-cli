package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

// mockAPI implements trello.API for command testing.
type mockAPI struct {
	trello.API
	listBoardsFn           func(ctx context.Context) ([]trello.Board, error)
	getBoardFn             func(ctx context.Context, id string) (trello.Board, error)
	createBoardFn          func(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error)
	listListsFn            func(ctx context.Context, boardID string) ([]trello.List, error)
	createListFn           func(ctx context.Context, boardID, name string) (trello.List, error)
	updateListFn           func(ctx context.Context, listID string, params trello.UpdateListParams) (trello.List, error)
	archiveListFn          func(ctx context.Context, listID string) (trello.List, error)
	moveListFn             func(ctx context.Context, listID, boardID string, pos *float64) (trello.List, error)
	listCardsByBoardFn     func(ctx context.Context, boardID string) ([]trello.Card, error)
	listCardsByListFn      func(ctx context.Context, listID string) ([]trello.Card, error)
	getCardFn              func(ctx context.Context, cardID string) (trello.Card, error)
	createCardFn           func(ctx context.Context, params trello.CreateCardParams) (trello.Card, error)
	updateCardFn           func(ctx context.Context, cardID string, params trello.UpdateCardParams) (trello.Card, error)
	moveCardFn             func(ctx context.Context, cardID, listID string, pos *float64) (trello.Card, error)
	archiveCardFn          func(ctx context.Context, cardID string) (trello.Card, error)
	deleteCardFn           func(ctx context.Context, cardID string) error
	listCommentsFn         func(ctx context.Context, cardID string) ([]trello.Comment, error)
	addCommentFn           func(ctx context.Context, cardID, text string) (trello.Comment, error)
	updateCommentFn        func(ctx context.Context, actionID, text string) (trello.Comment, error)
	deleteCommentFn        func(ctx context.Context, actionID string) error
	listChecklistsFn       func(ctx context.Context, cardID string) ([]trello.Checklist, error)
	createChecklistFn      func(ctx context.Context, cardID, name string) (trello.Checklist, error)
	deleteChecklistFn      func(ctx context.Context, checklistID string) error
	addCheckItemFn         func(ctx context.Context, checklistID, name string) (trello.CheckItem, error)
	updateCheckItemFn      func(ctx context.Context, cardID, itemID, state string) (trello.CheckItem, error)
	deleteCheckItemFn      func(ctx context.Context, checklistID, itemID string) error
	listAttachmentsFn      func(ctx context.Context, cardID string) ([]trello.Attachment, error)
	addFileAttachmentFn    func(ctx context.Context, cardID, filePath string, name *string) (trello.Attachment, error)
	addURLAttachmentFn     func(ctx context.Context, cardID, urlStr string, name *string) (trello.Attachment, error)
	deleteAttachmentFn     func(ctx context.Context, cardID, attachmentID string) error
	listLabelsFn           func(ctx context.Context, boardID string) ([]trello.Label, error)
	createLabelFn          func(ctx context.Context, boardID, name, color string) (trello.Label, error)
	addLabelToCardFn       func(ctx context.Context, cardID, labelID string) error
	removeLabelFromCardFn  func(ctx context.Context, cardID, labelID string) error
	listCustomFieldsByBoardFn  func(ctx context.Context, boardID string) ([]trello.CustomField, error)
	getCustomFieldFn           func(ctx context.Context, fieldID string) (trello.CustomField, error)
	createCustomFieldFn        func(ctx context.Context, params trello.CreateCustomFieldParams) (trello.CustomField, error)
	updateCustomFieldFn        func(ctx context.Context, fieldID string, params trello.UpdateCustomFieldParams) (trello.CustomField, error)
	deleteCustomFieldFn        func(ctx context.Context, fieldID string) error
	listCustomFieldOptionsFn   func(ctx context.Context, fieldID string) ([]trello.CustomFieldOption, error)
	createCustomFieldOptionFn  func(ctx context.Context, fieldID string, params trello.CreateCustomFieldOptionParams) (trello.CustomFieldOption, error)
	updateCustomFieldOptionFn  func(ctx context.Context, fieldID, optionID string, params trello.UpdateCustomFieldOptionParams) (trello.CustomFieldOption, error)
	deleteCustomFieldOptionFn  func(ctx context.Context, fieldID, optionID string) error
	listCardCustomFieldItemsFn func(ctx context.Context, cardID string) ([]trello.CardCustomFieldItem, error)
	setCardCustomFieldItemFn   func(ctx context.Context, cardID, fieldID string, params trello.SetCardCustomFieldItemParams) (trello.CardCustomFieldItem, error)
	clearCardCustomFieldItemFn func(ctx context.Context, cardID, fieldID string) error
	listMembersFn          func(ctx context.Context, boardID string) ([]trello.Member, error)
	addMemberToCardFn      func(ctx context.Context, cardID, memberID string) error
	removeMemberFromCardFn func(ctx context.Context, cardID, memberID string) error
	searchCardsFn          func(ctx context.Context, query string) (trello.CardSearchResult, error)
	searchBoardsFn         func(ctx context.Context, query string) (trello.BoardSearchResult, error)
}

func (m *mockAPI) ListBoards(ctx context.Context) ([]trello.Board, error) {
	if m.listBoardsFn != nil {
		return m.listBoardsFn(ctx)
	}
	return nil, nil
}

func (m *mockAPI) GetBoard(ctx context.Context, id string) (trello.Board, error) {
	if m.getBoardFn != nil {
		return m.getBoardFn(ctx, id)
	}
	return trello.Board{}, nil
}

func (m *mockAPI) CreateBoard(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error) {
	if m.createBoardFn != nil {
		return m.createBoardFn(ctx, params)
	}
	return trello.Board{}, nil
}

func (m *mockAPI) ListCustomFieldsByBoard(ctx context.Context, boardID string) ([]trello.CustomField, error) {
	if m.listCustomFieldsByBoardFn != nil {
		return m.listCustomFieldsByBoardFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) GetCustomField(ctx context.Context, fieldID string) (trello.CustomField, error) {
	if m.getCustomFieldFn != nil {
		return m.getCustomFieldFn(ctx, fieldID)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) CreateCustomField(ctx context.Context, params trello.CreateCustomFieldParams) (trello.CustomField, error) {
	if m.createCustomFieldFn != nil {
		return m.createCustomFieldFn(ctx, params)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) UpdateCustomField(ctx context.Context, fieldID string, params trello.UpdateCustomFieldParams) (trello.CustomField, error) {
	if m.updateCustomFieldFn != nil {
		return m.updateCustomFieldFn(ctx, fieldID, params)
	}
	return trello.CustomField{}, nil
}

func (m *mockAPI) DeleteCustomField(ctx context.Context, fieldID string) error {
	if m.deleteCustomFieldFn != nil {
		return m.deleteCustomFieldFn(ctx, fieldID)
	}
	return nil
}

func (m *mockAPI) ListCustomFieldOptions(ctx context.Context, fieldID string) ([]trello.CustomFieldOption, error) {
	if m.listCustomFieldOptionsFn != nil {
		return m.listCustomFieldOptionsFn(ctx, fieldID)
	}
	return nil, nil
}

func (m *mockAPI) CreateCustomFieldOption(ctx context.Context, fieldID string, params trello.CreateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
	if m.createCustomFieldOptionFn != nil {
		return m.createCustomFieldOptionFn(ctx, fieldID, params)
	}
	return trello.CustomFieldOption{}, nil
}

func (m *mockAPI) UpdateCustomFieldOption(ctx context.Context, fieldID, optionID string, params trello.UpdateCustomFieldOptionParams) (trello.CustomFieldOption, error) {
	if m.updateCustomFieldOptionFn != nil {
		return m.updateCustomFieldOptionFn(ctx, fieldID, optionID, params)
	}
	return trello.CustomFieldOption{}, nil
}

func (m *mockAPI) DeleteCustomFieldOption(ctx context.Context, fieldID, optionID string) error {
	if m.deleteCustomFieldOptionFn != nil {
		return m.deleteCustomFieldOptionFn(ctx, fieldID, optionID)
	}
	return nil
}

func (m *mockAPI) ListCardCustomFieldItems(ctx context.Context, cardID string) ([]trello.CardCustomFieldItem, error) {
	if m.listCardCustomFieldItemsFn != nil {
		return m.listCardCustomFieldItemsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) SetCardCustomFieldItem(ctx context.Context, cardID, fieldID string, params trello.SetCardCustomFieldItemParams) (trello.CardCustomFieldItem, error) {
	if m.setCardCustomFieldItemFn != nil {
		return m.setCardCustomFieldItemFn(ctx, cardID, fieldID, params)
	}
	return trello.CardCustomFieldItem{}, nil
}

func (m *mockAPI) ClearCardCustomFieldItem(ctx context.Context, cardID, fieldID string) error {
	if m.clearCardCustomFieldItemFn != nil {
		return m.clearCardCustomFieldItemFn(ctx, cardID, fieldID)
	}
	return nil
}

func TestBoardsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listBoardsFn: func(ctx context.Context) ([]trello.Board, error) {
			return []trello.Board{
				{ID: "b1", Name: "Board One"},
				{ID: "b2", Name: "Board Two"},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}
	data, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(data))
	}
}

func TestBoardsListRequiresAuth(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"boards", "list"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards list should fail without auth")
	}
}

func TestBoardsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getBoardFn: func(ctx context.Context, id string) (trello.Board, error) {
			if id != "b1" {
				t.Errorf("board ID = %q, want %q", id, "b1")
			}
			return trello.Board{ID: "b1", Name: "My Board"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "get", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards get failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	data := envelope["data"].(map[string]any)
	if data["id"] != "b1" {
		t.Errorf("data.id = %v, want b1", data["id"])
	}
}

func TestBoardsGetMissingFlag(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"boards", "get"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards get should fail without --board")
	}
}

func TestBoardsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	desc := "Board description"
	orgID := "org1"
	sourceBoardID := "template1"
	apiClient = &mockAPI{
		createBoardFn: func(ctx context.Context, params trello.CreateBoardParams) (trello.Board, error) {
			if params.Name != "Project Board" {
				t.Fatalf("params.Name = %q, want %q", params.Name, "Project Board")
			}
			if params.Desc == nil || *params.Desc != desc {
				t.Fatalf("params.Desc = %v, want %q", params.Desc, desc)
			}
			if params.DefaultLists == nil || !*params.DefaultLists {
				t.Fatalf("params.DefaultLists = %v, want true", params.DefaultLists)
			}
			if params.DefaultLabels == nil || !*params.DefaultLabels {
				t.Fatalf("params.DefaultLabels = %v, want true", params.DefaultLabels)
			}
			if params.IDOrganization == nil || *params.IDOrganization != orgID {
				t.Fatalf("params.IDOrganization = %v, want %q", params.IDOrganization, orgID)
			}
			if params.IDBoardSource == nil || *params.IDBoardSource != sourceBoardID {
				t.Fatalf("params.IDBoardSource = %v, want %q", params.IDBoardSource, sourceBoardID)
			}
			return trello.Board{ID: "b3", Name: params.Name, Desc: desc}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{
		"boards", "create",
		"--name", "Project Board",
		"--desc", desc,
		"--default-lists",
		"--default-labels",
		"--organization", orgID,
		"--source-board", sourceBoardID,
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards create failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	data := envelope["data"].(map[string]any)
	if data["id"] != "b3" {
		t.Errorf("data.id = %v, want b3", data["id"])
	}
}

func TestBoardsCreateMissingName(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	assertContractCode(t, executeRootArgs("boards", "create"), "VALIDATION_ERROR")
}

func TestBoardsHelpIncludesCreate(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"boards", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards help failed: %v", err)
	}

	help := buf.String()
	for _, want := range []string{"create", "Create a board"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q\nhelp:\n%s", want, help)
		}
	}
}

func TestBoardsCreateHelpIncludesFlags(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"boards", "create", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards create help failed: %v", err)
	}

	help := buf.String()
	for _, want := range []string{"--name", "--desc", "--default-lists", "--default-labels", "--organization", "--source-board"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q\nhelp:\n%s", want, help)
		}
	}
}
