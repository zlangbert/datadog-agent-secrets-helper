# Datadog Agent Secrets Helper

[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/zlangbert/datadog-agent-secrets-helper)](https://hub.docker.com/r/zlangbert/datadog-agent-secrets-helper)
[![Go Report Card](https://goreportcard.com/badge/github.com/zlangbert/datadog-agent-secrets-helper)](https://goreportcard.com/report/github.com/zlangbert/datadog-agent-secrets-helper)

The [Datadog Agent](https://github.com/DataDog/datadog-agent) often needs access to sensitive values in its 
configuration to preform checks against protected resources. As check configuration is often stored in version control
or accessible on the filesystem, the ability to securely store and use sensitive information in check configuration
is necessary. 

The agent offers a flexible method to define and retrieve secret values in agent or check configuration described 
[here](https://github.com/DataDog/datadog-agent/blob/master/docs/agent/secrets.md). This repository contains an
implementation of a secrets backend with support for multiple secrets providers. 

TODO: installation and usage docs