package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/forstmeier/askpaulgraham/pkg/cnt"
	"github.com/forstmeier/askpaulgraham/pkg/db"
	"github.com/forstmeier/askpaulgraham/pkg/dct"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type mockCntClient struct {
	mockGetItemsOutput []cnt.ItemXML
	mockGetItemsError  error
	mockGetTextOutput  *string
	mockGetTextError   error
}

func (m *mockCntClient) GetItems(ctx context.Context, address string) ([]cnt.ItemXML, error) {
	return m.mockGetItemsOutput, m.mockGetItemsError
}

func (m *mockCntClient) GetText(ctx context.Context, address string) (*string, error) {
	return m.mockGetTextOutput, m.mockGetTextError
}

type mockDBClient struct {
	mockGetIDsOutput        []string
	mockGetIDsError         error
	mockStoreSummariesError error
	mockStoreTextError      error
	mockGetDocumentsOutput  []dct.Document
	mockGetDocumentsError   error
	mockStoreDocumentsError error
}

func (m *mockDBClient) GetIDs(ctx context.Context) ([]string, error) {
	return m.mockGetIDsOutput, m.mockGetIDsError
}

func (m *mockDBClient) GetSummaries(ctx context.Context) ([]db.Summary, error) {
	return nil, nil
}

func (m *mockDBClient) StoreSummaries(ctx context.Context, summaries []db.Summary) error {
	return m.mockStoreSummariesError
}

func (m *mockDBClient) StoreText(ctx context.Context, id, text string) error {
	return m.mockStoreTextError
}

func (m *mockDBClient) GetDocuments(ctx context.Context) ([]dct.Document, error) {
	return m.mockGetDocumentsOutput, m.mockGetDocumentsError
}

func (m *mockDBClient) StoreDocuments(ctx context.Context, documents []dct.Document) error {
	return m.mockStoreDocumentsError
}

func (m *mockDBClient) StoreQuestion(ctx context.Context, id, question string) error {
	return nil
}

func (m *mockDBClient) StoreAnswer(ctx context.Context, id, answer string) error {
	return nil
}

type mockNLPClient struct {
	mockGetSummaryOutput  *string
	mockGetSummaryError   error
	mockSetDocumentsError error
}

func (m *mockNLPClient) GetSummary(ctx context.Context, text string) (*string, error) {
	return m.mockGetSummaryOutput, m.mockGetSummaryError
}

func (m *mockNLPClient) SetDocuments(ctx context.Context, documents []dct.Document) error {
	return m.mockSetDocumentsError
}

func (m *mockNLPClient) GetAnswer(ctx context.Context, question, userID string) (*string, error) {
	return nil, nil
}

func Test_handler(t *testing.T) {
	mockGetItemsErr := errors.New("mock get items error")
	mockGetIDsErr := errors.New("mock get ids error")
	mockGetTextErr := errors.New("mock get text error")
	mockGetSummaryErr := errors.New("mock get summary error")
	mockStoreSummariesErr := errors.New("mock store summaries error")
	mockStoreTextErr := errors.New("mock store text error")
	mockGetDocumentsErr := errors.New("mock get documents error")
	mockSetDocumentsErr := errors.New("mock set documents error")
	mockStoreDocumentsErr := errors.New("mock store documents error")

	mockText := "full text"
	mockSummary := "summary"

	tests := []struct {
		description             string
		mockGetItemsOutput      []cnt.ItemXML
		mockGetItemsError       error
		mockGetIDsOutput        []string
		mockGetIDsError         error
		mockGetTextOutput       *string
		mockGetTextError        error
		mockGetSummaryOutput    *string
		mockGetSummaryError     error
		mockStoreSummariesError error
		mockStoreTextError      error
		mockGetDocumentsOutput  []dct.Document
		mockGetDocumentsError   error
		mockSetDocumentsError   error
		mockStoreDocumentsError error
		error                   error
	}{
		{
			description:        "error getting items",
			mockGetItemsOutput: nil,
			mockGetItemsError:  mockGetItemsErr,
			error:              mockGetItemsErr,
		},
		{
			description:        "error getting ids",
			mockGetItemsOutput: []cnt.ItemXML{},
			mockGetItemsError:  nil,
			mockGetIDsOutput:   nil,
			mockGetIDsError:    mockGetIDsErr,
			error:              mockGetIDsErr,
		},
		{
			description: "error getting text",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError: nil,
			mockGetIDsOutput:  []string{"old_id"},
			mockGetIDsError:   nil,
			mockGetTextOutput: nil,
			mockGetTextError:  mockGetTextErr,
			error:             mockGetTextErr,
		},
		{
			description: "error getting summary",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:    nil,
			mockGetIDsOutput:     []string{"old_id"},
			mockGetIDsError:      nil,
			mockGetTextOutput:    &mockText,
			mockGetTextError:     nil,
			mockGetSummaryOutput: nil,
			mockGetSummaryError:  mockGetSummaryErr,
			error:                mockGetSummaryErr,
		},
		{
			description: "error storing summary",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: mockStoreSummariesErr,
			error:                   mockStoreSummariesErr,
		},
		{
			description: "error storing text",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: nil,
			mockStoreTextError:      mockStoreTextErr,
			error:                   mockStoreTextErr,
		},
		{
			description: "error getting documents",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: nil,
			mockStoreTextError:      nil,
			mockGetDocumentsOutput:  nil,
			mockGetDocumentsError:   mockGetDocumentsErr,
			error:                   mockGetDocumentsErr,
		},
		{
			description: "error setting documents",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: nil,
			mockStoreTextError:      nil,
			mockGetDocumentsOutput:  []dct.Document{},
			mockGetDocumentsError:   nil,
			mockSetDocumentsError:   mockSetDocumentsErr,
			error:                   mockSetDocumentsErr,
		},
		{
			description: "error storing documents",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: nil,
			mockStoreTextError:      nil,
			mockGetDocumentsOutput:  []dct.Document{},
			mockGetDocumentsError:   nil,
			mockSetDocumentsError:   nil,
			mockStoreDocumentsError: mockStoreDocumentsErr,
			error:                   mockStoreDocumentsErr,
		},
		{
			description: "successful invocation",
			mockGetItemsOutput: []cnt.ItemXML{
				{
					Link: "http://www.paulgraham.com/new_id.html",
				},
			},
			mockGetItemsError:       nil,
			mockGetIDsOutput:        []string{"old_id"},
			mockGetIDsError:         nil,
			mockGetTextOutput:       &mockText,
			mockGetTextError:        nil,
			mockGetSummaryOutput:    &mockSummary,
			mockGetSummaryError:     nil,
			mockStoreSummariesError: nil,
			mockStoreTextError:      nil,
			mockGetDocumentsOutput: []dct.Document{
				{
					Metadata: "mock_id",
					Text:     "mock text",
				},
			},
			mockGetDocumentsError: nil,
			mockSetDocumentsError: nil,
			error:                 nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			c := &mockCntClient{
				mockGetItemsOutput: test.mockGetItemsOutput,
				mockGetItemsError:  test.mockGetItemsError,
				mockGetTextOutput:  test.mockGetTextOutput,
				mockGetTextError:   test.mockGetTextError,
			}

			d := &mockDBClient{
				mockGetIDsOutput:        test.mockGetIDsOutput,
				mockGetIDsError:         test.mockGetIDsError,
				mockStoreSummariesError: test.mockStoreSummariesError,
				mockStoreTextError:      test.mockStoreTextError,
				mockGetDocumentsOutput:  test.mockGetDocumentsOutput,
				mockGetDocumentsError:   test.mockGetDocumentsError,
				mockStoreDocumentsError: test.mockStoreDocumentsError,
			}

			n := &mockNLPClient{
				mockGetSummaryOutput:  test.mockGetSummaryOutput,
				mockGetSummaryError:   test.mockGetSummaryError,
				mockSetDocumentsError: test.mockSetDocumentsError,
			}

			handlerFunc := handler(c, d, n, "rss_url")

			err := handlerFunc(context.Background())

			if err != test.error {
				t.Errorf("incorrect error, received: %v, expected: %v", err, test.error)
			}
		})
	}
}
