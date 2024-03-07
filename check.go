package main

import (
	"context"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/grafana/regexp"
	"golang.org/x/exp/slices"
)

type checkResult struct {
	// ReviewSatisfied indicates that *any* review has been made on the PR. It is also set to
	// true if the test plan indicates that this PR does not need to be review.
	ReviewSatisfied bool
	// CanSkipTestPlan indicates that the test plan is not required for audit.
	CanSkipTestPlan bool
	// TestPlan is the content provided after the acceptance checklist checkbox.
	TestPlan string
	// ProtectedBranch indicates that the base branch for this PR is protected and merges
	// are considered to be exceptional and should always be justified.
	ProtectedBranch bool
	// Error indicating any issue that might have occured during the check.
	Error error
}

func (r checkResult) IsSatisfied() bool {
	return r.IsTestPlanSatisfied() && r.ReviewSatisfied && !r.ProtectedBranch
}

func (r checkResult) IsTestPlanSatisfied() bool {
	return r.CanSkipTestPlan || r.TestPlan != ""
}

var (
	testPlanDividerRegexp       = regexp.MustCompile("(?m)(#+ Test [pP]lan)|(Test [pP]lan:)")
	noReviewNeededDividerRegexp = regexp.MustCompile("(?m)([nN]o [rR]eview [rR]equired:)")

	markdownCommentRegexp = regexp.MustCompile("<!--((.|\n)*?)-->(\n)*")

	noReviewNeededLabels = []string{"no-review-required", "automerge"}
)

type checkOpts struct {
	SkipReviews        bool
	SkipReviewForUsers string
	SkipTestPlan       bool
	ProtectedBranch    string
}

func isProtectedBranch(payload *EventPayload, protectedBranch string) bool {
	return protectedBranch != "" && payload.PullRequest.Base.Ref == protectedBranch
}

func isReviewSatisfied(payload *EventPayload, opts checkOpts, checker ApprovalChecker) bool {
	pr := payload.PullRequest

	if opts.SkipReviews {
		return true
	}

	if opts.SkipReviewForUsers != "" {
		isPRAuthorExcluded := slices.Contains(strings.Split(opts.SkipReviewForUsers, ","), pr.User.Login)

		if isPRAuthorExcluded {
			return true
		}
	}

	// If the PR has explicit review comments, great we're done
	if pr.ReviewComments > 0 {
		return true
	}

	// Look for no review required explanation in the body
	if sections := noReviewNeededDividerRegexp.Split(pr.Body, 2); len(sections) > 1 {
		noReviewRequiredExplanation := cleanMarkdown(sections[1])
		if len(noReviewRequiredExplanation) > 0 {
			return true
		}
	}

	// Look for no review required labels
	for _, label := range pr.Labels {
		if slices.Contains(noReviewNeededLabels, label.Name) {
			return true
		}
	}

	// Else we have to check for an approval through the GitHub API
	return checker.IsApproved(*payload)
}

func checkPR(payload *EventPayload, opts checkOpts, checker ApprovalChecker) checkResult {
	pr := payload.PullRequest

	// Whether or not this PR was reviewed can be inferred from payload, but an approval
	// might not have any comments so we need to double-check through the GitHub API
	reviewed := isReviewSatisfied(payload, opts, checker)

	// Parse test plan data from body
	sections := testPlanDividerRegexp.Split(pr.Body, 2)
	if len(sections) < 2 {
		return checkResult{
			ReviewSatisfied: reviewed,
			CanSkipTestPlan: opts.SkipTestPlan,
		}
	}

	mergeAgainstProtected := isProtectedBranch(payload, opts.ProtectedBranch)

	return checkResult{
		ReviewSatisfied: reviewed,
		CanSkipTestPlan: opts.SkipTestPlan,
		TestPlan:        cleanMarkdown(sections[1]),
		ProtectedBranch: mergeAgainstProtected,
	}
}

func cleanMarkdown(s string) string {
	content := s
	// Remove comments
	content = markdownCommentRegexp.ReplaceAllString(content, "")
	// Remove whitespace
	content = strings.TrimSpace(content)

	return content
}

type ApprovalChecker interface {
	IsApproved(EventPayload) bool
}

type GithubApprovalChecker struct {
	client *github.Client
	ctx    context.Context
}

func (checker GithubApprovalChecker) IsApproved(payload EventPayload) bool {
	owner, repo := payload.Repository.GetOwnerAndName()
	var reviews []*github.PullRequestReview

	reviews, _, _ = checker.client.PullRequests.ListReviews(checker.ctx, owner, repo, payload.PullRequest.Number, &github.ListOptions{})
	return len(reviews) > 0
}
