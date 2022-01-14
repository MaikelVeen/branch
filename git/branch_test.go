package git

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	testCases := map[string]struct {
		input  string
		expect string
	}{
		"square brackets and its contents are removed": {
			input:  "[Bug] Fix the issues with branch",
			expect: "fix-the-issues-with-branch",
		},
		"round brackets and its contents are removed": {
			input:  "(Bug) Fix the issues with branch",
			expect: "fix-the-issues-with-branch",
		},
		"lowercase": {
			input:  "THIS ISSUE IS ANGRY",
			expect: "this-issue-is-angry",
		},
		"special characters are removed": {
			input:  `!@#$%^&*()_+?<>/:{}~\|[],.`,
			expect: "",
		},
		"special characters between words end up as hyphen": {
			input:  "test the issue/bug",
			expect: "test-the-issue-bug",
		},
		"long titles are cutoff to 12 max": {
			input:  "there is absolutely no way to explain this issue in less than twelve words",
			expect: "there-is-absolutely-no-way-to-explain-this-issue-in-less-than",
		},
		"trailing spaces are not hyphenated": {
			input:  "trailing spaces are not hyphenated ",
			expect: "trailing-spaces-are-not-hyphenated",
		},
		"leading spaces are not hyphenated": {
			input:  " leading spaces are not hyphenated",
			expect: "leading-spaces-are-not-hyphenated",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expect, FormatAsValidRef(tc.input))
		})
	}

	t.Run("spaces are always removed to one spaces", func(t *testing.T) {
		t.Parallel()

		start := 'a'
		end := 'b'

		for i := 1; i < 100; i++ {
			var sb strings.Builder
			sb.WriteRune(start)

			for j := 0; j < i; j++ {
				sb.WriteRune(' ')
			}

			sb.WriteRune(end)

			assert.Equal(t, "a-b", FormatAsValidRef(sb.String()))
			sb.Reset()
		}
	})
}

func TestGetBranchName(t *testing.T) {
	testCases := map[string]struct {
		base   string
		key    string
		title  string
		expect string
	}{
		"feature branch": {
			base:   "feature",
			key:    "DONUT-25",
			title:  "We need to fix this issue",
			expect: "feature/DONUT-25/we-need-to-fix-this-issue",
		},
		"hotfix branch": {
			base:   "hotfix",
			key:    "DONUT-26",
			title:  "Also this issue needs some love",
			expect: "hotfix/DONUT-26/also-this-issue-needs-some-love",
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expect, GetBranchName(tc.base, tc.key, tc.title))
		})
	}
}
