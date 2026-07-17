// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2026- Leonardo Di Donato <leodidonato@gmail.com>

// Command actionpins verifies that this repository's GitHub Actions workflows
// use reviewed, immutable action references.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	checkoutSHA     = "9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0"
	checkoutVersion = "v7.0.0"
	setupGoSHA      = "b7ad1dad31e06c5925ef5d2fc7ad053ef454303e"
	setupGoVersion  = "v7.0.0"
)

var fullCommitSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

func main() {
	workflows := flag.String("workflows", ".github/workflows", "directory containing workflow YAML files")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: actionpins [-workflows directory]")
		os.Exit(2)
	}

	if err := checkWorkflows(*workflows); err != nil {
		fmt.Fprintf(os.Stderr, "actionpins: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("action pin tests: PASS")
}

type policyState struct {
	seenCheckout bool
	seenSetupGo  bool
}

func checkWorkflows(workflows string) error {
	paths, err := workflowPaths(workflows)
	if err != nil {
		return err
	}

	var state policyState
	for _, path := range paths {
		if err := checkWorkflow(path, &state); err != nil {
			return err
		}
	}

	if !state.seenCheckout {
		return fmt.Errorf("no actions/checkout reference found")
	}
	if !state.seenSetupGo {
		return fmt.Errorf("no actions/setup-go reference found")
	}

	return nil
}

func workflowPaths(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		if extension == ".yml" || extension == ".yaml" {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(paths)

	return paths, nil
}

func checkWorkflow(path string, state *policyState) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return fmt.Errorf("%s: parse workflow: %w", path, err)
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("%s: workflow must contain exactly one YAML document", path)
		}

		return fmt.Errorf("%s: parse workflow: %w", path, err)
	}

	return checkDocument(path, &document, state)
}

func checkDocument(path string, document *yaml.Node, state *policyState) error {
	if document.Kind != yaml.DocumentNode || len(document.Content) != 1 {
		return nodeError(path, document, "workflow must contain exactly one YAML document")
	}

	root, err := mappingNode(document.Content[0])
	if err != nil {
		return nodeError(path, document.Content[0], "workflow root must be a mapping")
	}
	_, jobsNode, found, err := mappingEntry(root, "jobs")
	if err != nil {
		return nodeError(path, root, "cannot inspect jobs")
	}
	if !found {
		return nil
	}

	jobs, err := mappingNode(jobsNode)
	if err != nil {
		return nodeError(path, jobsNode, "jobs must be a mapping")
	}
	// GitHub Actions executes only reusable-job and step-level uses fields.
	// Deliberately ignore arbitrary data keys named "uses" in env or inputs.
	for index := 0; index < len(jobs.Content); index += 2 {
		if err := checkJob(path, jobs.Content[index+1], state); err != nil {
			return err
		}
	}

	return nil
}

func checkJob(path string, node *yaml.Node, state *policyState) error {
	job, err := mappingNode(node)
	if err != nil {
		return nodeError(path, node, "job must be a mapping")
	}

	_, usesNode, found, err := mappingEntry(job, "uses")
	if err != nil {
		return nodeError(path, job, "cannot inspect reusable workflow reference")
	}
	if found {
		if _, checkErr := checkUses(path, usesNode, state); checkErr != nil {
			return checkErr
		}
	}

	_, stepsNode, found, err := mappingEntry(job, "steps")
	if err != nil {
		return nodeError(path, job, "cannot inspect steps")
	}
	if !found {
		return nil
	}

	return checkSteps(path, stepsNode, state)
}

func checkSteps(path string, node *yaml.Node, state *policyState) error {
	steps, err := sequenceNode(node)
	if err != nil {
		return nodeError(path, node, "steps must be a sequence")
	}
	for _, stepNode := range steps.Content {
		step, err := mappingNode(stepNode)
		if err != nil {
			return nodeError(path, stepNode, "step must be a mapping")
		}

		_, usesNode, found, err := mappingEntry(step, "uses")
		if err != nil {
			return nodeError(path, step, "cannot inspect action reference")
		}
		if !found {
			continue
		}

		action, err := checkUses(path, usesNode, state)
		if err != nil {
			return err
		}
		if strings.EqualFold(action, "actions/checkout") {
			if err := checkCheckoutInputs(path, step); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkUses(path string, node *yaml.Node, state *policyState) (string, error) {
	actionRef, value, err := scalarValue(node)
	if err != nil {
		return "", nodeError(path, node, "uses value must be a string")
	}
	if strings.HasPrefix(actionRef, "./") {
		// Local actions execute checked-out repository code, so they are not
		// external supply-chain references and GitHub does not permit a ref.
		return "", nil
	}
	if strings.HasPrefix(actionRef, "docker://") {
		return "", nodeError(path, value, "Docker actions are not allowed")
	}

	action, ref, hasRef := strings.Cut(actionRef, "@")
	if !hasRef || action == "" {
		return "", nodeError(path, value, "external action has no ref: %s", actionRef)
	}
	if !fullCommitSHA.MatchString(ref) {
		return "", nodeError(path, value, "external action is not pinned to a full commit: %s", actionRef)
	}

	switch strings.ToLower(action) {
	case "actions/checkout":
		state.seenCheckout = true
		if actionRef != "actions/checkout@"+checkoutSHA {
			return "", nodeError(path, value, "expected actions/checkout@%s", checkoutSHA)
		}
		if value.LineComment != "# "+checkoutVersion {
			return "", nodeError(path, value, "expected # %s annotation", checkoutVersion)
		}
	case "actions/setup-go":
		state.seenSetupGo = true
		if actionRef != "actions/setup-go@"+setupGoSHA {
			return "", nodeError(path, value, "expected actions/setup-go@%s", setupGoSHA)
		}
		if value.LineComment != "# "+setupGoVersion {
			return "", nodeError(path, value, "expected # %s annotation", setupGoVersion)
		}
	}

	return action, nil
}

func checkCheckoutInputs(path string, step *yaml.Node) error {
	_, withNode, found, err := mappingEntry(step, "with")
	if err != nil {
		return nodeError(path, step, "cannot inspect checkout inputs")
	}
	if !found {
		return nil
	}

	with, err := mappingNode(withNode)
	if err != nil {
		return nodeError(path, withNode, "checkout inputs must be a mapping")
	}
	unsafeKey, _, found, err := mappingEntry(with, "allow-unsafe-pr-checkout")
	if err != nil {
		return nodeError(path, with, "cannot inspect checkout inputs")
	}
	if found {
		return nodeError(path, unsafeKey, "allow-unsafe-pr-checkout must not be configured")
	}

	return nil
}

func mappingEntry(mapping *yaml.Node, name string) (*yaml.Node, *yaml.Node, bool, error) {
	mapping, err := mappingNode(mapping)
	if err != nil {
		return nil, nil, false, err
	}
	for index := 0; index < len(mapping.Content); index += 2 {
		key, value := mapping.Content[index], mapping.Content[index+1]
		keyValue, _, err := scalarValue(key)
		if err == nil && keyValue == name {
			return key, value, true, nil
		}
	}

	return nil, nil, false, nil
}

func mappingNode(node *yaml.Node) (*yaml.Node, error) {
	node, err := dereferenceNode(node)
	if err != nil || node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("not a mapping")
	}

	return node, nil
}

func sequenceNode(node *yaml.Node) (*yaml.Node, error) {
	node, err := dereferenceNode(node)
	if err != nil || node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("not a sequence")
	}

	return node, nil
}

func scalarValue(node *yaml.Node) (string, *yaml.Node, error) {
	node, err := dereferenceNode(node)
	if err != nil || node.Kind != yaml.ScalarNode {
		return "", node, fmt.Errorf("not a scalar")
	}

	return node.Value, node, nil
}

func dereferenceNode(node *yaml.Node) (*yaml.Node, error) {
	seen := make(map[*yaml.Node]struct{})
	for node != nil && node.Kind == yaml.AliasNode {
		if _, exists := seen[node]; exists {
			return nil, fmt.Errorf("cyclic alias")
		}
		seen[node] = struct{}{}
		node = node.Alias
	}
	if node == nil {
		return nil, fmt.Errorf("nil node")
	}

	return node, nil
}

func nodeError(path string, node *yaml.Node, format string, args ...any) error {
	return fmt.Errorf("%s:%d: %s", path, node.Line, fmt.Sprintf(format, args...))
}
