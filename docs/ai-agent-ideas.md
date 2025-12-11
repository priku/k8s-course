# AI Agent Ideas for K8s Infrastructure

This document contains ideas for implementing AI agents across our Kubernetes infrastructure and applications, focusing on Azure AI services for AI-102 learning and practical DevOps automation.

---

## 1. DevOps/CI-CD Assistant Agent ⭐ (Priority)

**MCP Servers:** Kubernetes + GitHub + Azure

### What the Agent Could Do
- Monitor deployment health across all namespaces
- Auto-diagnose failed deployments and suggest fixes
- Analyze GitHub Actions workflow failures
- Optimize resource allocation (CPU/memory limits)
- Generate deployment reports
- Auto-create issues for failed builds
- Predict deployment failures based on patterns

### Example Use Cases
```
"Check why the last deployment to project namespace failed"
"Are there any pods in CrashLoopBackOff state?"
"Optimize resource requests for todo-project based on actual usage"
"Create a summary report of this week's deployments"
"What's the root cause of the failed workflow run #123?"
"Suggest fixes for the current failing deployment"
```

### Technical Implementation

#### MCP Servers Needed
- **Kubernetes MCP** - Monitor pods, deployments, services
  - Official Google GKE MCP (adapt for AKS)
  - Or custom MCP server for Azure Kubernetes Service
- **GitHub MCP** - Access workflows, issues, PRs
- **Azure MCP** - Access Azure resources and metrics

#### Azure AI Services (AI-102 Coverage)
- **Azure OpenAI** - Root cause analysis, suggestion generation
- **Azure Monitor / Log Analytics** - Metrics and logs aggregation
- **Anomaly Detector** - Detect unusual deployment patterns
- **Text Analytics** - Analyze error messages and logs

#### Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    AI Agent (Claude)                     │
│              via Model Context Protocol                  │
└─────────────┬───────────────────┬──────────────┬────────┘
              │                   │              │
    ┌─────────▼─────────┐ ┌──────▼──────┐ ┌────▼─────┐
    │  Kubernetes MCP   │ │ GitHub MCP  │ │ Azure MCP│
    │   - Get pods      │ │ - Workflows │ │ - Metrics│
    │   - Get logs      │ │ - Issues    │ │ - Logs   │
    │   - Get events    │ │ - PRs       │ │ - Costs  │
    └───────────────────┘ └─────────────┘ └──────────┘
              │                   │              │
    ┌─────────▼───────────────────▼──────────────▼────────┐
    │           Azure Kubernetes Service (AKS)             │
    │  Namespaces: project, exercises, nginx-gateway, etc. │
    └──────────────────────────────────────────────────────┘
```

#### Data Flow
1. User asks: "Why did the last deployment fail?"
2. Agent queries GitHub MCP for recent workflow runs
3. Agent queries Kubernetes MCP for pod status/events
4. Agent queries Azure MCP for resource metrics
5. Agent sends all context to Azure OpenAI for analysis
6. Azure OpenAI generates root cause and suggestions
7. Agent returns human-readable explanation with fixes

### Implementation Phases

#### Phase 1: Read-Only Monitoring (Week 1)
- Install/configure necessary MCP servers
- Query pod status across namespaces
- Read GitHub workflow results
- Basic health checks

**Commands:**
```
"Show me all pods in all namespaces"
"What's the status of the last 5 workflow runs?"
"Are there any failing pods right now?"
```

#### Phase 2: Analysis & Diagnostics (Week 2)
- Integrate Azure OpenAI for analysis
- Parse logs for error patterns
- Correlate GitHub failures with K8s events
- Generate diagnostic reports

**Commands:**
```
"Analyze why workflow run #123 failed"
"What's causing the CrashLoopBackOff in todo-project?"
"Generate a health report for all services"
```

#### Phase 3: Automated Actions (Week 3)
- Auto-create GitHub issues for failures
- Suggest resource optimizations
- Auto-restart failed pods (with approval)
- Schedule reports

**Commands:**
```
"Create an issue for the todo-backend deployment failure"
"Apply the suggested resource optimization"
"Restart all pods in error state"
```

### Azure Resources Required

```yaml
# Azure Resources
- Azure OpenAI Service
  - Model: GPT-4 or GPT-3.5-turbo
  - Deployment: devops-assistant

- Azure Monitor / Log Analytics Workspace
  - Connected to AKS cluster
  - Retention: 30 days

- Azure Application Insights (optional)
  - Track agent performance
  - Monitor API calls
```

### AI-102 Learning Outcomes

✅ **Azure OpenAI Integration** - Core exam topic
✅ **Prompt Engineering** - Context building, chain-of-thought
✅ **API Authentication** - Managed Identity, API keys
✅ **Error Handling** - Retry logic, graceful degradation
✅ **Monitoring & Logging** - Application Insights integration
✅ **Real-world AI application** - Production use case

---

## 2. Infrastructure Cost Optimizer Agent

**MCP Servers:** Azure + Kubernetes metrics

### What the Agent Could Do
- Analyze Azure resource costs (AKS nodes, storage, networking)
- Suggest right-sizing for AKS nodes
- Identify unused namespaces/resources
- Recommend reserved instances or savings plans
- Alert on cost spikes or anomalies
- Predict monthly costs based on trends
- Compare actual vs. budgeted costs

### Example Use Cases
```
"What's our monthly AKS spend breakdown?"
"Which namespaces are using the most resources?"
"Suggest cost optimizations for our cluster"
"Are we paying for any unused resources?"
"What would we save by reserving our current instances?"
"Alert me if daily costs exceed $50"
```

### Azure AI Services
- **Azure OpenAI** - Cost analysis and recommendations
- **Anomaly Detector** - Detect unusual spending patterns
- **Azure Cost Management API** - Real cost data

---

## 3. Security & Compliance Agent

**MCP Servers:** Kubernetes + Azure Security Center

### What the Agent Could Do
- Scan for vulnerabilities in container images
- Check for exposed secrets in configs
- Audit RBAC permissions
- Monitor for security policy violations (Gatekeeper)
- Generate compliance reports (PCI, SOC2, etc.)
- Detect misconfigurations
- Track CVEs in deployed images

### Example Use Cases
```
"Scan all deployments for security vulnerabilities"
"Are there any secrets exposed in our manifests?"
"Check if all pods follow our security policies"
"Generate a SOC2 compliance report"
"What CVEs exist in our current images?"
"Audit who has cluster-admin access"
```

### Azure AI Services
- **Azure OpenAI** - Policy interpretation, report generation
- **Azure Security Center** - Vulnerability data
- **Azure Defender for Containers** - Runtime protection

---

## 4. Database Management Agent

**MCP Servers:** PostgreSQL + Monitoring tools

### What the Agent Could Do
- Monitor PostgreSQL StatefulSet health
- Analyze slow queries and execution plans
- Suggest index optimizations
- Backup/restore automation
- Predict storage needs based on growth
- Detect connection leaks
- Query optimization recommendations

### Example Use Cases
```
"Show me slow queries from the todos database"
"Is the postgres-ss pod healthy?"
"When will we run out of PVC storage?"
"Suggest indexes to improve performance"
"Backup the todos database to Azure Blob Storage"
"What's the current connection pool usage?"
```

### Azure AI Services
- **Azure OpenAI** - Query analysis, optimization suggestions
- **Azure Database for PostgreSQL** - Managed service integration
- **Log Analytics** - Query performance insights

---

## 5. Todo App Intelligence Agent (AI-102 Perfect!)

**MCP Servers:** Gmail + Azure OpenAI + Todo Backend API

### What the Agent Could Do
- Parse emails → automatically create todos
- Analyze todo completion patterns and productivity
- Smart task scheduling suggestions
- Auto-categorize and tag todos based on content
- Generate daily/weekly productivity summaries
- Voice-to-todo via Speech API
- Multi-language support via Translator
- Sentiment analysis for task urgency

### Example Use Cases
```
"Read my unread emails and create todos"
"What tasks should I focus on today?"
"Summarize my week's productivity"
"Add todo via voice: [speak into mic]"
"Translate this task to Spanish"
"Which tasks seem most urgent based on content?"
"Suggest optimal time slots for my tasks"
```

### Azure AI Services (Full AI-102 Coverage!)
- **Azure OpenAI** - Task extraction, NLP, suggestions
- **Speech-to-Text** - Voice input for todos
- **Text-to-Speech** - Read todos aloud
- **Text Analytics** - Sentiment analysis, key phrase extraction
- **Translator** - Multi-language support
- **Language Understanding (LUIS)** - Intent recognition for commands

### Implementation Details

#### Email → Todo Extraction
```
User email:
"Hi Team, please review the Q4 report by Friday and
schedule a meeting with the client next week."

Azure OpenAI extracts:
[
  {
    "description": "Review Q4 report",
    "dueDate": "Friday",
    "priority": "high"
  },
  {
    "description": "Schedule meeting with client",
    "dueDate": "next week",
    "priority": "medium"
  }
]
```

#### API Integration
```go
// New endpoint in todo-backend
POST /api/todos/parse-email
{
  "email_content": "email text here"
}

// Response
{
  "tasks": [
    {
      "description": "...",
      "due_date": "...",
      "priority": "..."
    }
  ]
}
```

---

## 6. Log Analysis & Troubleshooting Agent

**MCP Servers:** Kubernetes logs + Azure Monitor

### What the Agent Could Do
- Aggregate logs from all pods across namespaces
- Detect error patterns and anomalies
- Correlate errors across microservices
- Generate incident reports
- Predict failures based on log patterns
- Natural language log queries
- Auto-tag incidents by severity

### Example Use Cases
```
"Show me all errors from the last hour"
"Why is todo-backend restarting?"
"Analyze logs for the failed deployment"
"Are there any error patterns in the last 24 hours?"
"What happened before the pod crash?"
"Search logs for 'database timeout'"
```

### Azure AI Services
- **Azure OpenAI** - Log pattern analysis, incident summarization
- **Log Analytics** - Log aggregation and KQL queries
- **Anomaly Detector** - Detect unusual error rates
- **Text Analytics** - Extract entities from error messages

---

## 7. Traffic & Performance Agent

**MCP Servers:** Prometheus/Grafana metrics + Application Insights

### What the Agent Could Do
- Analyze request patterns and traffic trends
- Detect anomalies in traffic (DDoS, bot traffic)
- Performance bottleneck identification
- Load testing recommendations
- Auto-scaling suggestions based on usage
- Response time predictions
- Capacity planning

### Example Use Cases
```
"What's the response time trend for todo-project?"
"Detect any traffic anomalies this week"
"Should we scale up the backend?"
"What's our 95th percentile response time?"
"Predict traffic for next month"
"Recommend optimal HPA settings"
```

### Azure AI Services
- **Azure OpenAI** - Analysis and recommendations
- **Anomaly Detector** - Traffic pattern detection
- **Azure Monitor** - Metrics collection
- **Application Insights** - APM data

---

## 8. Branch Environment Manager Agent

**MCP Servers:** GitHub + Kubernetes + Custom APIs

### What the Agent Could Do
- Auto-create feature environments on branch creation
- Manage environment lifecycle
- Copy production data to test environments
- Auto-cleanup stale environments (we already have this!)
- Environment status dashboard
- Resource quota management per environment
- Auto-generate environment URLs

### Example Use Cases
```
"Create a test environment from main for feature X"
"List all active feature environments"
"Clone production data to my feature branch environment"
"Clean up environments older than 7 days"
"What's the status of the feature-xyz environment?"
"Generate a shareable URL for my feature branch"
```

### Azure AI Services
- **Azure OpenAI** - Environment planning, naming suggestions
- **Azure DevOps API** - Pipeline integration

### Current Implementation
We already have:
- ✅ Automatic namespace creation per branch
- ✅ Automatic cleanup on branch deletion
- ✅ Branch name sanitization
- ✅ Isolated environments

Could add:
- 🔄 Data seeding automation
- 🔄 Environment comparison reports
- 🔄 Cost tracking per environment
- 🔄 Auto-expiration policies

---

## Implementation Priority Recommendations

### For AI-102 Learning (Best Coverage)
1. **Todo App Intelligence Agent** - Uses most Azure AI services
2. **Log Analysis Agent** - Good for Text Analytics and Anomaly Detection
3. **DevOps Assistant** - Azure OpenAI integration patterns

### For Immediate Practical Value
1. **DevOps/CI-CD Assistant** - Solves real operational problems
2. **Cost Optimizer** - Direct business value
3. **Security & Compliance** - Risk mitigation

### For Portfolio/Resume
1. **Todo App Intelligence** - User-facing, impressive demo
2. **DevOps Assistant** - Shows DevOps + AI expertise
3. **Security Agent** - High-value security automation

---

## Next Steps

### This Week: DevOps/CI-CD Assistant - Phase 1
1. Research available MCP servers for Kubernetes/GitHub/Azure
2. Set up Azure OpenAI resource
3. Test basic queries: pod status, workflow results
4. Document findings

### Next Week: Build Core Features
1. Implement diagnostic analysis
2. Create GitHub issue automation
3. Add resource optimization suggestions

### Future: Expand to Other Agents
Based on success and learning from DevOps Assistant

---

## Resources

### MCP Servers
- Google Kubernetes Engine (GKE) - Official
- GitHub MCP - Community or custom
- Azure MCP - Check official support

### Azure AI-102 Documentation
- Azure OpenAI Service
- Azure Monitor / Log Analytics
- Anomaly Detector API
- Text Analytics API

### Our Infrastructure
- AKS Cluster: dwk-aks-cluster
- Resource Group: dwk-aks-rg
- Namespaces: project, exercises, nginx-gateway
- GitHub Workflows: deploy-project.yml, cleanup-environment.yml

---

*Document created: 2025-12-10*
*Last updated: 2025-12-10*
