---
applyTo: "**/*.yml,**/*.yaml"
---

# GitHub Actions Workflow File Standards

## Required Header Block

Every workflow file must start with a structured comment header immediately after `name:`:

```yaml
# ============================================================
# Purpose: [Brief description of what this workflow does]
# Triggers: [When this workflow runs]
# Permissions: [What permissions this workflow needs]
# Inputs: [What Inputs the Workflow is using, can be empty]
# Outputs: [What outputs the workflow is creating, can be empty
# Secrets: [Which Secretes the Workflow needs]
# Variables [What Variables are used]
# Precondition [For instance any plumbing required like IdP Trust Relations]
# ============================================================
```

## Example

```yaml
name: CI Pipeline

# ============================================================
# Purpose: Run linting, tests, and security checks on push/PR
# Triggers: push to main, pull requests
# Permissions: contents:read, pull-requests:write
# Inputs: Text Files in the /corpus directory
# Outputs: A Release with the intermediate vector store
# Secrets: None
# Variables: None
# Preconditions: None
# ============================================================

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
```

## Best Practices

1. **Always include the header block** - It provides quick context for workflow purpose and configuration
2. **Keep descriptions concise** - Use brief, clear language in the Purpose field
3. **List all triggers** - Document all events that trigger the workflow
4. **Document permissions** - Clearly state what permissions the workflow requires
5. **Document inputs** - Clearly state what inputs the workflow requires
6. **Document outputs** - Clearly state what outputs the workflow requires
7. **Document secrets** - Clearly state what github secrets the workflow requires
5. **Document variables** - Clearly state what github variables the workflow requires
6. **Document precondition** - Clearly state what type of setup has to be done before the workflow can be done
8. **Use consistent formatting** - Follow the exact format shown in the template
