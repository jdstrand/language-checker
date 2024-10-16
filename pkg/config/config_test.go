package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jdstrand/language-checker/pkg/rule"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("check-logger", func(t *testing.T) {
		out := &bytes.Buffer{}
		log.Logger = zerolog.New(out)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		// isolate config rules
		loaded, err := loadConfig("testdata/good.yaml")
		assert.NoError(t, err)
		configRules := make([]string, len(loaded.Rules))
		for i := range loaded.Rules {
			configRules[i] = fmt.Sprintf("%q", loaded.Rules[i].Name)
		}

		// isolate default rules
		defaultRules := make([]string, len(rule.DefaultRules))
		for i := range defaultRules {
			defaultRules[i] = fmt.Sprintf("%q", rule.DefaultRules[i].Name)
		}

		c, err := NewConfig("testdata/good.yaml", false)
		assert.NoError(t, err)
		enabledRules := make([]string, len(c.Rules))
		for i := range c.Rules {
			enabledRules[i] = fmt.Sprintf("%q", c.Rules[i].Name)
		}

		loadedRemoteConfigMsg := `{"level":"debug","filename":"testdata/good.yaml","message":"Adding custom ruleset from"}`
		loadedRemoteConfig := `{"level":"debug","filename":"testdata/good.yaml","message":"Adding custom ruleset from"}`
		loadedConfigMsg := `{"level":"debug","config":"testdata/good.yaml","message":"loaded config file"}`
		configRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"config rules"}`, strings.Join(configRules, ","))
		defaultRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"default rules"}`, strings.Join(defaultRules, ","))
		allRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"all enabled rules"}`, strings.Join(enabledRules, ","))
		assert.Equal(t,
			loadedRemoteConfigMsg+"\n"+loadedRemoteConfig+"\n"+loadedConfigMsg+"\n"+configRulesMsg+"\n"+defaultRulesMsg+"\n"+allRulesMsg+"\n",
			out.String())
	})

	t.Run("config-good", func(t *testing.T) {
		c, err := NewConfig("testdata/good.yaml", false)
		assert.NoError(t, err)

		expectedRules := []*rule.Rule{}
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule1",
			Terms:        []string{"rule1"},
			Alternatives: []string{"alt-rule1"},
			Severity:     rule.SevWarn,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule2",
			Terms:        []string{"rule2", "rule-2"},
			Alternatives: []string{"alt-rule2", "alt-rule-2"},
			Severity:     rule.SevError,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "whitelist",
			Terms:        []string{"rulewl", "rule-wl"},
			Alternatives: []string{"alt-rulewl", "alt-rule-wl"},
			Severity:     rule.SevError,
		})

		expected := &Config{
			Rules:       expectedRules,
			IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
		}
		expected.ConfigureRules(false)

		assert.EqualValues(t, expected.Rules, c.Rules)

		// check default config message
		assert.Equal(t, "No findings found.", c.GetSuccessExitMessage())
	})

	t.Run("config-empty-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("", false)
		assert.NoError(t, err)

		expectedEmpty := &Config{
			Rules:       rule.DefaultRules,
			IgnoreFiles: []string(nil),
		}
		assert.Equal(t, expectedEmpty, c)
	})

	t.Run("config-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/missing.yaml", false)
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("config-empty-success-message", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/empty-success-message.yaml", false)
		assert.NoError(t, err)

		// check default config message
		assert.Equal(t, "", c.GetSuccessExitMessage())
	})

	t.Run("config-custom-success-message", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/custom-success-message.yaml", false)
		assert.NoError(t, err)

		// check default config message
		assert.Equal(t, "this is a test", c.GetSuccessExitMessage())
	})

	t.Run("config-add-note-messaage", func(t *testing.T) {
		// Test when it is configured to add a note to the output message
		c, err := NewConfig("testdata/add-note-message.yaml", false)
		assert.NoError(t, err)

		// check global IncludeNote
		assert.Equal(t, true, c.IncludeNote)

		// check IncludeNote is set for rule2
		assert.Equal(t, true, *c.Rules[1].Options.IncludeNote)

		// check IncludeNote is not overridden for rule1
		assert.Equal(t, false, *c.Rules[0].Options.IncludeNote)
	})

	t.Run("config-dont-add-note-message", func(t *testing.T) {
		// Test when it is nott configured to add a note to the output message
		c, err := NewConfig("testdata/dont-add-note-message.yaml", false)
		assert.NoError(t, err)

		// check global IncludeNote
		assert.Equal(t, false, c.IncludeNote)

		// check IncludeNote is not set for rule2
		assert.Equal(t, false, *c.Rules[1].Options.IncludeNote)

		// check IncludeNote is not overridden for rule1
		assert.Equal(t, true, *c.Rules[0].Options.IncludeNote)
	})

	t.Run("disable-default-rules", func(t *testing.T) {
		c, err := NewConfig("testdata/good.yaml", true)
		assert.NoError(t, err)
		assert.Len(t, c.Rules, 3)

		c, err = NewConfig("testdata/good.yaml", false)
		assert.NoError(t, err)
		assert.Len(t, c.Rules, len(rule.DefaultRules)+2)
	})

	t.Run("config-exclude-a-category", func(t *testing.T) {
		// Test when configured to exclude a single category
		c, err := NewConfig("testdata/exclude-single-category.yaml", false)
		assert.NoError(t, err)

		expectedRules := []*rule.Rule{}
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule1",
			Terms:        []string{"rule1"},
			Alternatives: []string{"alt-rule1"},
			Severity:     rule.SevWarn,
			Options:      rule.Options{Categories: []string{"cat1"}},
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule3",
			Terms:        []string{"rule3", "rule-3"},
			Alternatives: []string{"alt-rule3", "alt-rule-3"},
			Severity:     rule.SevError,
		})

		expected := &Config{
			Rules:             expectedRules,
			ExcludeCategories: []string{"cat2"},
		}
		expected.ConfigureRules(false)

		assert.EqualValues(t, expected.Rules, c.Rules)
		assert.Equal(t, "No findings found.", c.GetSuccessExitMessage())
	})

	t.Run("config-exclude-multiple-categories", func(t *testing.T) {
		// Test when configured to exclude multiple categories
		c, err := NewConfig("testdata/exclude-multiple-categories.yaml", false)
		assert.NoError(t, err)

		expectedRules := []*rule.Rule{}
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule3",
			Terms:        []string{"rule3", "rule-3"},
			Alternatives: []string{"alt-rule3", "alt-rule-3"},
			Severity:     rule.SevError,
		})

		expected := &Config{
			Rules:             expectedRules,
			ExcludeCategories: []string{"cat1", "cat2"},
		}
		expected.ConfigureRules(false)

		assert.EqualValues(t, expected.Rules, c.Rules)
		assert.Equal(t, "No findings found.", c.GetSuccessExitMessage())
	})

	t.Run("load-config-with-bad-url", func(t *testing.T) {
		_, err := NewConfig("https://raw.githubusercontent.com/get-woke/woke/main/example", false)
		assert.Error(t, err)
	})

	t.Run("load-config-with-url", func(t *testing.T) {
		c, err := NewConfig("https://raw.githubusercontent.com/get-woke/woke/main/example.yaml", false)
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("load-remote-config-valid-url", func(t *testing.T) {
		c, err := loadRemoteConfig("https://raw.githubusercontent.com/get-woke/woke/main/example.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("load-remote-config-invalid-url", func(t *testing.T) {
		_, err := loadRemoteConfig("https://raw.githubusercontent.com/get-woke/woke/main/example")
		assert.Error(t, err)
	})
}
func Test_relative(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	assert.Equal(t, ".langcheck.yml", relative(filepath.Join(cwd, ".langcheck.yml")))
	assert.Equal(t, ".langcheck.yml", relative(".langcheck.yml"))
	assert.Equal(t, "dir/.langcheck.yml", relative("dir/.langcheck.yml"))
}

func Test_RemoveRule(t *testing.T) {
	configRules := []*rule.Rule{}
	configRules = append(configRules, &rule.Rule{
		Name:         "rule1",
		Terms:        []string{"rule1"},
		Alternatives: []string{"alt-rule1"},
		Severity:     rule.SevWarn,
	})
	configRules = append(configRules, &rule.Rule{
		Name:         "rule2",
		Terms:        []string{"rule2"},
		Alternatives: []string{"alt-rule2"},
		Severity:     rule.SevWarn,
	})
	configRules = append(configRules, &rule.Rule{
		Name:         "whitelist",
		Terms:        []string{"rulewl", "rule-wl"},
		Alternatives: []string{"alt-rulewl", "alt-rule-wl"},
		Severity:     rule.SevError,
	})

	expected := &Config{
		Rules:       configRules,
		IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
	}

	expected.RemoveRule(1)

	// validate that the second rule, rule2, was removed
	expectedRules := []*rule.Rule{configRules[0], configRules[2]}

	assert.EqualValues(t, expected.Rules, expectedRules)
}
