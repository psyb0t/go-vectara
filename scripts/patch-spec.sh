#!/bin/bash
# Removes the invalid property-level required: false flag using the YAML-aware
# tool in internal/cmd/patch-spec. The query property remains optional.
# Renames JsonSchemaObject.additionalProperties in Go to avoid the generator's
# unknown-field map name. The JSON key stays additionalProperties.
# Replaces the spec's dummy Slack bot token example, which GitHub push
# protection rejects, with an obvious placeholder.
set -euo pipefail

go tool patch-spec
