# Branch Protection Setup Guide

To enable branch protection on the `main` branch, follow these steps:

## Via GitHub Web UI (Recommended)

1. Go to https://github.com/priku/k8s-course/settings/branches

2. Click "Add rule" or "Edit" if rule exists

3. Configure the following settings:

### Branch name pattern
```
main
```

### Protect matching branches

#### ✅ Require a pull request before merging
- **Required approvals**: 0 (for solo project) or 1 (for team)
- ☑️ Dismiss stale pull request approvals when new commits are pushed
- ☑️ Require approval of the most recent reviewable push

#### ✅ Require status checks to pass before merging
- ☑️ Require branches to be up to date before merging
- **Status checks that are required**:
  - None configured (checks will run based on file changes)
  - Terraform Plan runs automatically when terraform/ files change
  - Add required checks as needed for your workflow

#### ✅ Require conversation resolution before merging
- Ensures all PR comments are addressed

#### ⚠️ Do not require signed commits
- Optional: Enable if you set up GPG signing

#### ✅ Require linear history
- Prevents merge commits, enforces clean history

#### ✅ Include administrators
- Even admins must follow the rules (recommended for learning)

### Rules applied to everyone including administrators

#### ✅ Allow force pushes
- **Specify who can force push**: Nobody
- Prevents rewriting history

#### ✅ Allow deletions
- ☐ Allow deletions (keep unchecked to prevent accidental deletion)

4. Click **"Create"** or **"Save changes"**

## Verify Branch Protection

After setting up, test with:

```bash
# This should fail
git checkout main
echo "test" >> README.md
git commit -am "test: direct push"
git push origin main
# Expected: ! [remote rejected] main -> main (protected branch hook declined)

# This should work
git checkout -b test/branch-protection
git push -u origin test/branch-protection
gh pr create --title "Test: Verify branch protection" --body "Testing branch protection rules"
```

## Bypass Protection (Emergency Only)

If you absolutely must bypass protection (e.g., initial setup):

1. Temporarily disable branch protection
2. Make the critical change
3. Re-enable protection immediately
4. Document why in commit message

**Note:** This should be extremely rare and documented!

## Current Protection Status

Check status: https://github.com/priku/k8s-course/settings/branches

Expected configuration:
- ✅ Require pull request before merging
- ✅ Require status checks to pass
- ✅ Require conversation resolution
- ✅ Require linear history
- ✅ No force pushes
- ✅ No deletions

## Troubleshooting

### "You can't push directly to main"
**Solution:** This is correct! Create a feature branch and PR.

### "Status check failed"
**Solution:** Fix the CI errors before merging. Review the GitHub Actions logs.

### "PR can't be merged - conflicts"
**Solution:**
```bash
git checkout your-branch
git pull origin main
# Resolve conflicts
git push origin your-branch
```

### "I accidentally committed to main locally"
**Solution:**
```bash
# Don't push! Instead:
git reset HEAD~1  # Undo last commit, keep changes
git checkout -b feature/proper-branch
git commit -m "feat: proper commit message"
git push -u origin feature/proper-branch
```

## Best Practices

1. **Never disable protection** unless absolutely necessary
2. **Always create PRs** even for small changes
3. **Review Terraform plans** in PR comments before merging
4. **Keep PRs small** - easier to review
5. **Delete branches** after merging
6. **Use descriptive branch names** following convention

---

**Why We Do This:**
- ✅ Prevents accidental breaking changes
- ✅ Enforces code review process
- ✅ Creates audit trail
- ✅ Validates changes before deployment
- ✅ Professional development practice
