package git

// filter will first remove any values that are between () and [].
// Then it will filter all the special characters.
func filter(s string) string {
	return ""
}

func GetBranchName(key, base, title string) string {
	return ""
}

/*TODO: this function needs some love
func GetBranchNameFromIssue(issue jira.IssueBean) (string, error) {
	base := getBranchBase(issue)

	// TODO trim whitespace on end.
	filtered, err := removeSpecialChars(issue.Fields.Summary)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.ToLower(filtered), " ") // TODO: limit to ~12 entries
	hyphenated := strings.Join(parts, "-")

	// TODO: check if string would be a valid branch name
	return base + issue.Key + "-" + hyphenated, nil
}

func removeSpecialChars(s string) (string, error) {
	re, err := regexp.Compile(`[^\w!(-\/_ )]`)
	if err != nil {
		return "", nil
	}

	return re.ReplaceAllString(s, ""), nil
}

func getBranchBase(issue jira.IssueBean) string {
	if StringInSliceCaseInsensitive(issue.Fields.Issuetype.Name, []string{"bug"}) {
		return "hotfix/"
	}

	return "feature/"
}
*/
