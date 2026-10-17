// Command patch-spec removes an invalid required: false property flag from
// HallucinationCorrectionRequest.query. OpenAPI 3.0 requires an array of names
// on the parent schema. Removing this flag keeps the property optional.
package main

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	sourceSpec = "api/openapi.yaml"
	fixedSpec  = "api/openapi.fixed.yaml"
	fileMode   = 0o644
	pairSize   = 2
	// GitHub push protection rejects commits containing the spec's dummy Slack
	// bot token, which the generator copies into doc comments.
	slackExample     = `xoxb-\d{10}-\d{13}-[a-z]{24}`
	slackPlaceholder = "xoxb-EXAMPLE-BOT-TOKEN"
	proseWords       = "`[^`]*`|https?://\\S+|\\b([Ff]older|[Ss]ubfolder)s?\\b"
)

func main() {
	source, err := os.ReadFile(sourceSpec)
	if err != nil {
		panic(err)
	}

	var document yaml.Node
	if err := yaml.Unmarshal(source, &document); err != nil {
		panic(err)
	}

	patchSchema(property(&document,
		"components", "schemas", "HallucinationCorrectionRequest",
		"properties", "query",
	))
	renameAdditionalProperties(&document)

	words := regexp.MustCompile(strings.Join([]string{
		proseWords,
		" ?[—–] ?",
		"most powerful reranker available to",
		"comprehensive",
	}, "|"))
	normalizeDescriptions(&document, words)

	fixed, err := yaml.Marshal(&document)
	if err != nil {
		panic(err)
	}

	fixed = regexp.MustCompile(slackExample).ReplaceAll(
		fixed, []byte(slackPlaceholder),
	)

	if err := os.WriteFile(fixedSpec, fixed, fileMode); err != nil {
		panic(err)
	}
}

func property(node *yaml.Node, path ...string) *yaml.Node {
	if node.Kind == yaml.DocumentNode {
		node = node.Content[0]
	}

	for _, name := range path {
		var next *yaml.Node

		for i := 0; i+1 < len(node.Content); i += pairSize {
			if node.Content[i].Value == name {
				next = node.Content[i+1]

				break
			}
		}

		if next == nil {
			return nil
		}

		node = next
	}

	return node
}

func patchSchema(node *yaml.Node) {
	if node == nil {
		return
	}

	for i := 0; i+1 < len(node.Content); i += pairSize {
		key, value := node.Content[i], node.Content[i+1]
		if key.Value == "required" && value.Tag == "!!bool" &&
			value.Value == "false" {
			node.Content = append(
				node.Content[:i], node.Content[i+pairSize:]...,
			)

			return
		}
	}
}

// The named JSON Schema keyword collides with oapi-codegen's unknown-field
// map. x-go-name changes only the Go field, keeping both JSON behaviors.
func renameAdditionalProperties(document *yaml.Node) {
	schema := property(document, "components", "schemas", "JsonSchemaObject")
	if schema == nil {
		return
	}

	branches := []*yaml.Node{schema}
	if allOf := property(schema, "allOf"); allOf != nil {
		branches = append(branches, allOf.Content...)
	}

	for _, branch := range branches {
		node := property(branch, "properties", "additionalProperties")
		if node == nil || property(node, "x-go-name") != nil {
			continue
		}

		node.Content = append(node.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "x-go-name"},
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: "AllowsAdditionalProperties",
			},
		)
	}
}

// Only descriptions and summaries become generated documentation. Keep wire
// schemas, examples, URLs and inline identifiers intact during prose cleanup.
func normalizeDescriptions(node *yaml.Node, words *regexp.Regexp) {
	if node.Kind != yaml.MappingNode {
		for _, child := range node.Content {
			normalizeDescriptions(child, words)
		}

		return
	}

	for i := 0; i+1 < len(node.Content); i += pairSize {
		key, value := node.Content[i], node.Content[i+1]
		switch key.Value {
		case "example", "examples", "default", "enum", "const":
			continue
		}

		if (key.Value == "description" || key.Value == "summary") &&
			value.Kind == yaml.ScalarNode {
			value.Value = normalizeProse(value.Value, words)
		}

		normalizeDescriptions(value, words)
	}
}

func normalizeProse(text string, words *regexp.Regexp) string {
	replacer := strings.NewReplacer(
		"Subfolders", "Subdirectories", "subfolders", "subdirectories",
		"Subfolder", "Subdirectory", "subfolder", "subdirectory",
		"Folders", "Directories", "folders", "directories",
		"Folder", "Directory", "folder", "directory",
		"most powerful reranker available to", "reranker selected for",
		"comprehensive", "detailed",
	)

	return words.ReplaceAllStringFunc(text, func(word string) string {
		if strings.HasPrefix(word, "`") || strings.HasPrefix(word, "http") {
			return word
		}

		if strings.Contains(word, "—") {
			return ", "
		}

		if strings.Contains(word, "–") {
			return " to "
		}

		return replacer.Replace(word)
	})
}
