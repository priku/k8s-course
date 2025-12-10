# Contributing to DevOps with Kubernetes Course

## Professional Development Workflow

This project follows industry best practices for GitOps and Infrastructure as Code (IaC), even though it's a learning project. This demonstrates professional development standards.

## Branch Strategy

We use **GitHub Flow** - a simplified Git workflow suitable for continuous delivery:

```
main (protected)
  ↑
  └── feature/exercise-3.3 → PR → Review → Merge
  └── fix/ingress-path → PR → Review → Merge
```

### Branch Naming Convention

- `feature/exercise-X.Y` - New exercise implementations
- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

## Development Workflow

### 1. Create a Feature Branch

```bash
# Update main
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/exercise-3.3
```

### 2. Make Changes

Make your changes following these guidelines:

**For Infrastructure Changes (Terraform):**
- Update `.tf` files in `terraform/` directory
- Run `terraform fmt -recursive` before committing
- Test locally with `terraform plan` (read-only)
- **NEVER run `terraform apply` locally** - let CI/CD handle this

**For Application Changes:**
- Update application code
- Build and test locally
- Update manifests if needed
- Follow existing code style

**For Kubernetes Manifests:**
- Keep k3d and AKS manifests separate
- Document any new resources
- Include comments for complex configurations

### 3. Commit with Conventional Commits

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```bash
git add .
git commit -m "feat: implement exercise 3.3 with persistent volumes"
git commit -m "fix: correct ingress path routing"
git commit -m "docs: update README with AKS deployment steps"
git commit -m "refactor: simplify database connection logic"
```

**Commit Types:**
- `feat:` - New feature or exercise
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks
- `ci:` - CI/CD changes

### 4. Push to Remote

```bash
git push -u origin feature/exercise-3.3
```

### 5. Create Pull Request

```bash
# Using GitHub CLI (recommended)
gh pr create --title "feat: Implement Exercise 3.3" \
  --body "## Summary
Implements exercise 3.3 requirements for persistent volumes.

## Changes
- Added PersistentVolumeClaim for log storage
- Updated deployment to use volume mounts
- Tested with both k3d and AKS

## Testing
- [x] Tested locally with k3d
- [x] Verified volume persistence
- [ ] Deployed to AKS (will be done via CI/CD)

## Checklist
- [x] Code follows project conventions
- [x] Terraform formatted (if applicable)
- [x] Documentation updated
- [x] Ready for review"
```

**Or via GitHub Web UI:**
1. Go to https://github.com/priku/k8s-course
2. Click "Pull requests" → "New pull request"
3. Select your branch
4. Fill in title and description
5. Create PR

### 6. Automated CI Checks

Once PR is created, GitHub Actions will automatically:

✅ **Terraform Format Check** - Ensures code is properly formatted
✅ **Terraform Validate** - Validates Terraform configuration
✅ **Terraform Plan** - Shows what changes will be made
✅ **Plan Comment** - Posts Terraform plan to PR for review

**Review the Terraform plan carefully** in the PR comments before merging!

### 7. Review & Approve

For this learning project, you can self-review, but in a team:
- Another developer reviews your code
- Reviews the Terraform plan
- Suggests improvements
- Approves when satisfied

### 8. Merge Pull Request

```bash
# Via GitHub CLI
gh pr merge --squash --delete-branch

# Or via GitHub Web UI
# Click "Squash and merge"
# Delete branch after merge
```

### 9. Automated Deployment

After merge to `main`:
- CI/CD triggers automatically
- Terraform applies infrastructure changes (requires approval)
- Deployment is tracked in GitHub Actions

## What NOT to Do

❌ **Never push directly to `main`** - Always use PRs
❌ **Never run `terraform apply` locally** - Let CI/CD handle it
❌ **Never commit secrets** - Use GitHub Secrets or Azure Key Vault
❌ **Never skip CI checks** - They exist for a reason
❌ **Never force push to `main`** - This is blocked by branch protection

## What You SHOULD Do

✅ **Always create a feature branch**
✅ **Always write descriptive commit messages**
✅ **Always run `terraform fmt` before committing**
✅ **Always test locally before pushing** (plan, build, etc.)
✅ **Always review Terraform plans in PRs**
✅ **Always delete merged branches**

## Emergency Procedures

### If You Accidentally Pushed to Main

```bash
# This should be blocked by branch protection, but if it happens:
# 1. Create a revert commit
git revert <commit-hash>
git push origin main

# 2. Create a proper PR with the fix
git checkout -b fix/emergency-fix
# make fixes
git push -u origin fix/emergency-fix
gh pr create
```

### If Terraform State is Corrupted

```bash
# 1. Don't panic!
# 2. Check state in Azure Portal
# 3. Contact team lead (or in this case, review Terraform docs)
# 4. Use terraform state commands carefully:
terraform state list
terraform state show <resource>

# 5. Last resort: import resources
terraform import <resource_type>.<name> <azure_resource_id>
```

## Local Development

### Testing Terraform Changes

```bash
cd terraform

# Always format first
terraform fmt -recursive

# Initialize (if needed)
terraform init

# Validate syntax
terraform validate

# Preview changes (safe)
terraform plan

# STOP HERE - Don't run apply locally!
# Push to PR and let CI/CD handle it
```

### Testing Kubernetes Manifests

```bash
# Use k3d for local testing
k3d cluster create test-cluster

# Test manifests
kubectl apply -f <manifest>.yaml --dry-run=client
kubectl apply -f <manifest>.yaml --dry-run=server

# Deploy to local cluster
kubectl apply -f <manifest>.yaml

# Clean up
k3d cluster delete test-cluster
```

## Tools Required

- **Git** - Version control
- **GitHub CLI** (`gh`) - PR management
- **Terraform** - Infrastructure as Code
- **kubectl** - Kubernetes CLI
- **Docker** - Container runtime
- **k3d** - Local Kubernetes (optional for testing)
- **Azure CLI** - Azure management

## Getting Help

- **Documentation**: Check the [README.md](../README.md)
- **Terraform Docs**: https://registry.terraform.io/
- **Kubernetes Docs**: https://kubernetes.io/docs/
- **Course Material**: https://devopswithkubernetes.com/
- **Issues**: Create a GitHub issue for questions

## Code of Conduct

- Be professional
- Write clean, documented code
- Help others learn
- Admit when you don't know something
- Continuous improvement mindset

---

**Remember**: These practices may seem like overhead for a learning project, but they represent real-world professional standards. You're not just learning Kubernetes - you're learning how to work like a professional DevOps engineer! 🚀
