package cmd_test

import (
	"testing"

	"github.com/MaikelVeen/branch/pkg/cmd"
	"github.com/MaikelVeen/branch/pkg/jira"
	"github.com/stretchr/testify/require"
)

func TestBranchNameFromTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
		issue    *jira.Issue
		want     string
		wantErr  bool
	}{
		{
			name:     "valid template and issue",
			template: "{{.type}}/{{.key}}-{{.summary}}",
			issue: &jira.Issue{
				Key: "TEST-123",
				Fields: jira.IssueFields{
					Issuetype: jira.IssueType{Name: "Bug"},
					Summary:   "This is a test issue",
				},
			},
			want:    "bug/TEST-123-this-is-a-test-issue",
			wantErr: false,
		},
		{
			name:     "invalid template",
			template: "{{.invalid}/{{.key}}-{{.summary}}",
			issue: &jira.Issue{
				Key: "TEST-123",
				Fields: jira.IssueFields{
					Issuetype: jira.IssueType{Name: "Bug"},
					Summary:   "This is a test issue",
				},
			},
			wantErr: true,
		},
		{
			name:     "issue with special characters in summary",
			template: "{{.type}}/{{.key}}-{{.summary}}",
			issue: &jira.Issue{
				Key: "TEST-123",
				Fields: jira.IssueFields{
					Issuetype: jira.IssueType{Name: "Bug"},
					Summary:   "This is a test issue with special characters: !@#$%^&*()",
				},
			},
			want:    "bug/TEST-123-this-is-a-test-issue-with-special-characters",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := cmd.BranchNameFromTemplate(tt.template, tt.issue)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
