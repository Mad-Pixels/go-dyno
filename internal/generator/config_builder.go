package generator

import (
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

// ConfigBuilder provides an interface for building generator config objects.
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder creates a new configuration builder.
func NewConfigBuilder(schemaPath, outputDir string) *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			schemaPath: schemaPath,
			outputDir:  outputDir,
			scope:      ScopeAll,
		},
	}
}

// Build validates and returns the finalized config.
func (cb *ConfigBuilder) Build() (*Config, error) {
	var err error

	if err = utils.IsFileOrError(cb.config.schemaPath); err != nil {
		return nil, err
	}
	if err = utils.IsDirOrCreate(cb.config.outputDir); err != nil {
		return nil, err
	}

	if cb.config.packageName != nil {
		safe := utils.ToSafeName(*cb.config.packageName)
		cb.config.packageName = &safe
	}
	if cb.config.fileName != nil {
		safe := utils.AddFileExt(utils.ToSafeName(*cb.config.fileName), ".go")
		cb.config.fileName = &safe
	}

	if cb.config.verbose {
		logger.Log.Debug().Any("config", cb.config).Msg("Generator config was built")
	}
	return cb.config, nil
}

// WithScope overwrite codegen scope type.
func (cb *ConfigBuilder) WithScope(scope scope) *ConfigBuilder {
	cb.config.scope = scope
	return cb
}

// WithPackageName overwrite generated package name.
func (cb *ConfigBuilder) WithPackageName(name string) *ConfigBuilder {
	cb.config.packageName = &name
	return cb
}

// WithFileName overwrite generated file name.
func (cb *ConfigBuilder) WithFileName(name string) *ConfigBuilder {
	cb.config.fileName = &name
	return cb
}

// WithVerbose enables/disables verbose output.
func (cb *ConfigBuilder) WithVerbose(flag bool) *ConfigBuilder {
	cb.config.verbose = flag
	return cb
}

// WithDryRun enables/disables dry run mode.
func (cb *ConfigBuilder) WithDryRun(flag bool) *ConfigBuilder {
	cb.config.dryRun = flag
	return cb
}
