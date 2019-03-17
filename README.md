# AWS Secrets Manager Provider for the Datadog Agent

The [Datadog Agent](https://github.com/DataDog/datadog-agent) often needs access to sensitive values in its 
configuration to preform checks against protected resources. As check configuration is often stored in version control
or accessible on the filesystem, the ability to securely store and use sensitive information in check configuration
is necessary. 

The agent offers a flexible method to define and retrieve secret values in agent or check configuration described 
[here](https://github.com/DataDog/datadog-agent/blob/master/docs/agent/secrets.md). This repository contains an
implementation of a secrets provider backed by [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/), allowing
agent configuration to reference and retrieve secrets stored in Secrets Manager.

TODO: installation and usage docs