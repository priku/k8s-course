# DevOps/CI-CD Assistant Agent - Implementation Plan

A comprehensive guide to implementing an AI-powered DevOps assistant for our Kubernetes infrastructure using Azure OpenAI and Model Context Protocol (MCP).

---

## 🎯 Project Goals

Build an AI agent that can:
1. Monitor AKS cluster health across all namespaces
2. Diagnose deployment failures automatically
3. Analyze GitHub Actions workflow failures
4. Suggest fixes and optimizations
5. Generate operational reports
6. Auto-create issues for critical failures

---

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                  User (via Claude Code CLI)                   │
└────────────────────────────┬─────────────────────────────────┘
                             │
                             │ Natural Language
                             │
┌────────────────────────────▼─────────────────────────────────┐
│                    Claude AI Agent                            │
│              (with DevOps Assistant Context)                  │
└──────┬──────────────┬──────────────┬────────────────┬────────┘
       │              │              │                │
       │ MCP          │ MCP          │ MCP            │ Direct API
       │              │              │                │
┌──────▼──────┐ ┌────▼─────┐ ┌──────▼───────┐ ┌─────▼──────────┐
│ Kubernetes  │ │  GitHub  │ │    Azure     │ │ Azure OpenAI   │
│ MCP Server  │ │   MCP    │ │     MCP      │ │    Service     │
└──────┬──────┘ └────┬─────┘ └──────┬───────┘ └────────────────┘
       │              │              │
       │              │              │
┌──────▼──────────────▼──────────────▼───────┐
│        Azure Kubernetes Service (AKS)       │
│   Cluster: dwk-aks-cluster                  │
│   Resource Group: dwk-aks-rg                │
│                                             │
│   Namespaces:                               │
│   - project (main deployment)               │
│   - exercises                               │
│   - nginx-gateway                           │
│   - feature-* (branch environments)         │
└─────────────────────────────────────────────┘
```

---

## 📋 Phase 1: Setup & Basic Monitoring (Week 1)

### Objectives
- Set up Azure OpenAI resource
- Configure necessary MCP servers
- Implement basic cluster monitoring
- Test basic queries

### Tasks

#### 1.1 Azure OpenAI Setup

**Create Azure OpenAI Resource:**
```bash
# Resource details
Resource name: dwk-openai-devops
Region: East US (or Sweden Central)
Pricing tier: Standard S0
```

**Deploy Model:**
```bash
# Model deployment
Model: gpt-4 or gpt-3.5-turbo-16k
Deployment name: devops-assistant
Tokens per minute: 10K (start small)
```

**Get API Credentials:**
```bash
# Save these securely
AZURE_OPENAI_ENDPOINT=https://dwk-openai-devops.openai.azure.com/
AZURE_OPENAI_API_KEY=<your-key>
AZURE_OPENAI_DEPLOYMENT=devops-assistant
```

#### 1.2 MCP Server Research & Setup

**Option A: Use Existing MCP Servers**

1. **Kubernetes MCP Server**
   - Research: Check if Google's GKE MCP can work with AKS
   - Alternative: Look for community AKS MCP servers
   - Fallback: Create custom MCP server using kubectl

2. **GitHub MCP Server**
   - Search for existing GitHub MCP implementations
   - Test authentication with personal access token
   - Verify can read workflow runs and create issues

3. **Azure MCP Server**
   - Check official Google/Azure MCP support list
   - Research community implementations
   - Fallback: Direct Azure SDK integration

**Option B: Build Custom MCP Server (if needed)**

Create a simple MCP server wrapper around:
```bash
# Kubernetes operations
kubectl get pods --all-namespaces
kubectl describe pod <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace>
kubectl get events -n <namespace>

# GitHub CLI operations
gh run list --workflow=deploy-project.yml
gh run view <run-id> --log-failed
gh issue create --title "..." --body "..."
```

#### 1.3 Test Basic Queries

**Kubernetes Health Checks:**
```
Example queries to test:
- "Show all pods in the project namespace"
- "Are there any pods in error state?"
- "What's the status of the todo-project deployment?"
- "Show recent events in the project namespace"
- "Get logs from the todo-backend pod"
```

**GitHub Workflow Monitoring:**
```
Example queries:
- "Show the last 5 workflow runs for deploy-project"
- "What's the status of the latest deployment?"
- "Show me the logs from the failed workflow run"
```

**Expected Output Format:**
```json
{
  "namespace": "project",
  "pods": [
    {
      "name": "postgres-ss-0",
      "status": "Running",
      "restarts": 0,
      "age": "12h"
    },
    {
      "name": "todo-backend-dep-xxx",
      "status": "Running",
      "restarts": 0,
      "age": "12h"
    },
    {
      "name": "todo-project-dep-xxx",
      "status": "Running",
      "restarts": 0,
      "age": "12h"
    }
  ]
}
```

### Deliverables - Phase 1
- ✅ Azure OpenAI resource deployed
- ✅ At least one MCP server functional (Kubernetes preferred)
- ✅ Basic health check queries working
- ✅ Documentation of setup process

---

## 📋 Phase 2: Analysis & Diagnostics (Week 2)

### Objectives
- Integrate Azure OpenAI for intelligent analysis
- Build diagnostic capabilities
- Correlate failures across systems
- Generate actionable insights

### Tasks

#### 2.1 Build Diagnostic Prompts

**Root Cause Analysis Prompt Template:**
```
You are a DevOps expert analyzing a Kubernetes deployment failure.

Context:
- Namespace: {namespace}
- Deployment: {deployment_name}
- GitHub Workflow: {workflow_run_id}

Pod Status:
{pod_status_json}

Recent Events:
{k8s_events}

Pod Logs (last 100 lines):
{pod_logs}

GitHub Workflow Logs:
{workflow_logs}

Task: Analyze the failure and provide:
1. Root cause (most likely reason for failure)
2. Contributing factors
3. Step-by-step fix recommendations
4. Prevention strategies

Format your response as JSON:
{
  "root_cause": "...",
  "contributing_factors": ["...", "..."],
  "fix_steps": ["...", "..."],
  "prevention": ["...", "..."]
}
```

**Resource Optimization Prompt:**
```
You are a Kubernetes resource optimization expert.

Deployment: {deployment_name}
Current Resources:
{current_resource_requests}

Actual Usage (last 7 days):
{usage_metrics}

Task: Suggest optimal resource requests and limits.
Consider:
- Cost efficiency
- Performance headroom (20% buffer)
- Burst capacity needs

Provide recommendations as JSON:
{
  "current": { "requests": {...}, "limits": {...} },
  "recommended": { "requests": {...}, "limits": {...} },
  "savings_estimate": "$X/month",
  "rationale": "..."
}
```

#### 2.2 Implement Analysis Functions

**Core Analysis Capabilities:**

1. **Deployment Failure Analyzer**
   ```python
   def analyze_deployment_failure(namespace, deployment_name):
       # 1. Get pod status via K8s MCP
       pods = k8s_mcp.get_pods(namespace, deployment_name)

       # 2. Get events
       events = k8s_mcp.get_events(namespace)

       # 3. Get logs from failed pods
       logs = []
       for pod in pods:
           if pod.status != "Running":
               logs.append(k8s_mcp.get_logs(pod.name, namespace))

       # 4. Get workflow info
       workflow_run = github_mcp.get_latest_run()

       # 5. Send to Azure OpenAI for analysis
       context = build_context(pods, events, logs, workflow_run)
       analysis = azure_openai.analyze(context, diagnostic_prompt)

       return analysis
   ```

2. **Pattern Detection**
   ```python
   def detect_error_patterns(namespace, time_window="1h"):
       # Aggregate logs from all pods
       all_logs = k8s_mcp.get_all_logs(namespace, since=time_window)

       # Use Azure OpenAI to identify patterns
       prompt = f"Analyze these logs and identify error patterns: {all_logs}"
       patterns = azure_openai.analyze(prompt)

       return patterns
   ```

3. **Cross-Service Correlation**
   ```python
   def correlate_failures():
       # Get failures from different sources
       k8s_failures = k8s_mcp.get_failing_pods()
       github_failures = github_mcp.get_failed_workflows()

       # Use Azure OpenAI to correlate
       prompt = f"""
       Correlate these failures:
       K8s: {k8s_failures}
       GitHub: {github_failures}

       Are they related? What's the common root cause?
       """
       correlation = azure_openai.analyze(prompt)

       return correlation
   ```

#### 2.3 Build Reporting

**Daily Health Report:**
```
# Daily DevOps Health Report
Date: 2025-12-11

## Cluster Status
✅ All namespaces healthy
- project: 3/3 pods running
- exercises: 2/2 pods running
- nginx-gateway: 1/1 pods running

## Recent Deployments (Last 24h)
✅ deploy-project.yml: Success (3m 15s)
   - Branch: main
   - Commit: abc1234
   - Images: todo-project:main-abc1234, todo-backend:main-abc1234

## Issues Detected
⚠️  todo-backend pod restarted 2x due to OOMKill
   Recommendation: Increase memory limit from 256Mi to 512Mi

## Resource Utilization
- CPU: 45% average (good headroom)
- Memory: 78% average (consider adding nodes)
- Storage: 34% used

## Recommendations
1. Scale todo-backend horizontally (add 1 replica)
2. Increase memory limits for todo-backend
3. Clean up feature-exercise-3-7-branch-environments (merged)
```

### Deliverables - Phase 2
- ✅ Diagnostic prompts tested and refined
- ✅ Root cause analysis working
- ✅ Pattern detection functional
- ✅ Daily health reports generated
- ✅ Example analysis results documented

---

## 📋 Phase 3: Automated Actions (Week 3)

### Objectives
- Implement safe automated remediation
- Auto-create GitHub issues
- Resource optimization application
- Scheduled reporting

### Tasks

#### 3.1 GitHub Issue Automation

**Auto-Create Issues for Failures:**
```python
def auto_create_issue(failure_analysis):
    """
    Create GitHub issue when deployment fails
    """
    title = f"Deployment Failure: {failure_analysis['deployment']}"

    body = f"""
## Deployment Failure Report

**Namespace:** {failure_analysis['namespace']}
**Deployment:** {failure_analysis['deployment']}
**Time:** {failure_analysis['timestamp']}

### Root Cause
{failure_analysis['root_cause']}

### Contributing Factors
{format_list(failure_analysis['contributing_factors'])}

### Recommended Fixes
{format_list(failure_analysis['fix_steps'])}

### Prevention Strategies
{format_list(failure_analysis['prevention'])}

---
🤖 Auto-generated by DevOps Assistant Agent
Workflow Run: {failure_analysis['workflow_run_url']}
"""

    # Create issue via GitHub MCP
    issue = github_mcp.create_issue(
        title=title,
        body=body,
        labels=['deployment-failure', 'auto-created']
    )

    return issue
```

#### 3.2 Resource Optimization Application

**Apply Optimization (with approval):**
```python
def apply_resource_optimization(deployment, recommendations, auto_apply=False):
    """
    Apply resource optimization recommendations
    """
    if not auto_apply:
        # Present recommendations to user for approval
        print(f"Recommended resource changes for {deployment}:")
        print(f"Current: {recommendations['current']}")
        print(f"Recommended: {recommendations['recommended']}")
        print(f"Estimated savings: {recommendations['savings_estimate']}")

        approval = input("Apply these changes? (yes/no): ")
        if approval.lower() != 'yes':
            return "Optimization cancelled by user"

    # Update deployment via K8s MCP
    result = k8s_mcp.patch_deployment(
        deployment_name=deployment,
        resources=recommendations['recommended']
    )

    # Create PR with changes (better approach)
    create_optimization_pr(deployment, recommendations)

    return result
```

#### 3.3 Safe Automated Remediation

**Restart Failed Pods (with safety checks):**
```python
def auto_restart_failed_pods(namespace, max_restarts=3):
    """
    Automatically restart pods in CrashLoopBackOff
    Only if restart count < max_restarts
    """
    failed_pods = k8s_mcp.get_pods(
        namespace=namespace,
        status="CrashLoopBackOff"
    )

    for pod in failed_pods:
        if pod.restart_count >= max_restarts:
            # Create issue instead of restarting
            auto_create_issue({
                'title': f"Pod {pod.name} stuck in CrashLoopBackOff",
                'description': f"Restart count: {pod.restart_count}",
                'logs': k8s_mcp.get_logs(pod.name, namespace)
            })
        else:
            # Safe to restart
            k8s_mcp.delete_pod(pod.name, namespace)
            print(f"Restarted pod {pod.name} (restart {pod.restart_count})")
```

#### 3.4 Scheduled Reporting

**Cron-like Scheduled Tasks:**
```python
# Run daily at 9 AM
def daily_health_check():
    report = generate_health_report()

    # Send to Slack/Email/GitHub Discussion
    notify(report)

    # If critical issues, create GitHub issue
    if report.has_critical_issues():
        auto_create_issue(report.critical_issues)

# Run after each deployment
def post_deployment_check(deployment_id):
    # Wait 5 minutes for stabilization
    time.sleep(300)

    # Check deployment health
    health = check_deployment_health(deployment_id)

    if not health.is_healthy:
        analysis = analyze_deployment_failure(deployment_id)
        auto_create_issue(analysis)
```

### Deliverables - Phase 3
- ✅ GitHub issue automation working
- ✅ Resource optimization with approval flow
- ✅ Safe automated remediation rules
- ✅ Scheduled reporting implemented
- ✅ Safety guardrails tested

---

## 🔧 Technical Implementation Details

### Environment Setup

**Required Environment Variables:**
```bash
# Azure OpenAI
export AZURE_OPENAI_ENDPOINT="https://dwk-openai-devops.openai.azure.com/"
export AZURE_OPENAI_API_KEY="your-key-here"
export AZURE_OPENAI_DEPLOYMENT="devops-assistant"

# Azure AKS
export AKS_RESOURCE_GROUP="dwk-aks-rg"
export AKS_CLUSTER="dwk-aks-cluster"

# GitHub
export GITHUB_TOKEN="ghp_xxxxx"
export GITHUB_REPO="priku/k8s-course"

# Kubernetes
export KUBECONFIG="~/.kube/config"
```

### MCP Server Configuration

**Example MCP Server Config (claude_desktop_config.json):**
```json
{
  "mcpServers": {
    "kubernetes": {
      "command": "kubectl-mcp-server",
      "args": ["--kubeconfig", "~/.kube/config"],
      "env": {
        "CLUSTER": "dwk-aks-cluster"
      }
    },
    "github": {
      "command": "github-mcp-server",
      "env": {
        "GITHUB_TOKEN": "${GITHUB_TOKEN}",
        "GITHUB_REPO": "priku/k8s-course"
      }
    }
  }
}
```

### Azure OpenAI Integration

**Example Azure OpenAI Call:**
```python
from openai import AzureOpenAI

client = AzureOpenAI(
    api_key=os.getenv("AZURE_OPENAI_API_KEY"),
    api_version="2024-02-01",
    azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT")
)

def analyze_with_openai(context, prompt):
    response = client.chat.completions.create(
        model=os.getenv("AZURE_OPENAI_DEPLOYMENT"),
        messages=[
            {"role": "system", "content": "You are a DevOps expert."},
            {"role": "user", "content": f"{prompt}\n\nContext: {context}"}
        ],
        temperature=0.7,
        max_tokens=2000
    )

    return response.choices[0].message.content
```

---

## 🧪 Testing Strategy

### Unit Tests
```python
def test_deployment_failure_analysis():
    # Mock K8s MCP responses
    mock_pods = [{"name": "test-pod", "status": "CrashLoopBackOff"}]
    mock_events = [{"reason": "BackOff", "message": "..."}]

    # Test analysis
    result = analyze_deployment_failure("test-ns", "test-deploy")

    assert result['root_cause'] is not None
    assert len(result['fix_steps']) > 0
```

### Integration Tests
```python
def test_end_to_end_failure_detection():
    # Simulate a deployment failure
    deploy_broken_app()

    # Wait for agent to detect
    time.sleep(60)

    # Verify issue was created
    issues = github_mcp.list_issues(label='deployment-failure')
    assert len(issues) > 0
    assert 'auto-created' in issues[0].labels
```

### Manual Test Cases

1. **Deployment Failure Detection**
   - Deploy broken image
   - Verify agent detects failure
   - Check analysis quality

2. **Resource Optimization**
   - Query resource recommendations
   - Verify calculations are reasonable
   - Test approval flow

3. **Health Reporting**
   - Generate daily report
   - Verify all sections present
   - Check accuracy of metrics

---

## 📊 Success Metrics

### Technical Metrics
- Time to detect deployment failure: < 5 minutes
- Root cause accuracy: > 80% correct
- False positive rate: < 10%
- Analysis response time: < 30 seconds

### Business Metrics
- Reduction in MTTR (Mean Time To Recovery)
- Decrease in manual issue creation
- Cost savings from resource optimization
- Developer time saved

---

## 🚀 Future Enhancements

### Phase 4: Advanced Features
- Predictive failure analysis (ML-based)
- Auto-rollback on critical failures
- Capacity planning recommendations
- Performance trend analysis
- Multi-cluster support

### Phase 5: Integration Expansion
- Slack/Teams notifications
- PagerDuty integration
- Grafana dashboard automation
- Terraform drift detection

---

## 📚 Learning Resources

### Azure AI-102 Topics Covered
- ✅ Azure OpenAI Service deployment
- ✅ Prompt engineering and optimization
- ✅ API authentication and security
- ✅ Error handling and retry logic
- ✅ Token management and cost optimization
- ✅ Integration with external systems

### Documentation Links
- [Azure OpenAI Documentation](https://learn.microsoft.com/en-us/azure/ai-services/openai/)
- [Model Context Protocol Spec](https://modelcontextprotocol.io/)
- [Kubernetes API Reference](https://kubernetes.io/docs/reference/)
- [GitHub REST API](https://docs.github.com/en/rest)

---

## 📝 Next Steps

### This Week
1. Create Azure OpenAI resource
2. Research and test MCP servers
3. Implement basic health checks
4. Document setup process

### Week 2
1. Build diagnostic prompts
2. Test root cause analysis
3. Implement pattern detection
4. Create sample reports

### Week 3
1. Add GitHub issue automation
2. Build resource optimization
3. Test automated remediation
4. Deploy to production (with monitoring)

---

*Document created: 2025-12-10*
*Project: DevOps/CI-CD Assistant Agent*
*For: AI-102 Learning & Practical DevOps Automation*
