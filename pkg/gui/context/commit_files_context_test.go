package context

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestBranchDiffPullRequestDiff(t *testing.T) {
	tr := i18n.EnglishTranslationSet()
	feature := &models.Branch{Name: "feature"}
	branchDiff := &models.BranchDiff{
		Branch:      feature,
		BaseBranch:  "origin/main",
		BaseRefName: "main",
		MergeBase:   "base",
		Tip:         "tip",
	}

	scenarios := []struct {
		name         string
		branchDiff   *models.BranchDiff
		pullRequests map[string]*models.GithubPullRequest
		expected     types.PullRequestDiff
	}{
		{
			name:         "the pull request starts and ends where the diff does",
			branchDiff:   branchDiff,
			pullRequests: map[string]*models.GithubPullRequest{"feature": {BaseRefName: "main", HeadRefOid: "tip"}},
			expected:     types.PullRequestDiff{Branch: "feature", WholePullRequest: true},
		},
		{
			name:         "the branch has no pull request",
			branchDiff:   branchDiff,
			pullRequests: map[string]*models.GithubPullRequest{},
			expected:     types.PullRequestDiff{Branch: "feature", WholePullRequest: true},
		},
		{
			name:         "the pull request is opened against another branch",
			branchDiff:   branchDiff,
			pullRequests: map[string]*models.GithubPullRequest{"feature": {BaseRefName: "feature-a", HeadRefOid: "tip"}},
			expected: types.PullRequestDiff{
				Branch:           "feature",
				WholePullRequest: true,
				MismatchReason:   "The pull request is opened against feature-a, which this diff doesn't start from",
			},
		},
		{
			name:         "the pull request is at another commit",
			branchDiff:   branchDiff,
			pullRequests: map[string]*models.GithubPullRequest{"feature": {BaseRefName: "main", HeadRefOid: "pushed"}},
			expected: types.PullRequestDiff{
				Branch:           "feature",
				WholePullRequest: true,
				MismatchReason:   tr.BranchDiffNotAtPullRequestHead,
			},
		},
		{
			name: "a remote branch has no pull request of its own",
			branchDiff: &models.BranchDiff{
				Branch:      &models.RemoteBranch{Name: "feature", RemoteName: "origin"},
				BaseRefName: "main",
				Tip:         "tip",
			},
			pullRequests: map[string]*models.GithubPullRequest{"feature": {BaseRefName: "main", HeadRefOid: "tip"}},
			expected:     types.PullRequestDiff{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, branchDiffPullRequestDiff(s.branchDiff, s.pullRequests, tr))
		})
	}
}
