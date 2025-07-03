package generator

import (
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/Mad-Pixels/go-dyno/internal/utils/fs"
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

	if err = fs.IsFileOrError(cb.config.schemaPath); err != nil {
		return nil, err
	}
	if err = fs.IsDirOrCreate(cb.config.outputDir); err != nil {
		return nil, err
	}

	if cb.config.packageName != nil {
		safe := conv.ToSafeName(*cb.config.packageName)
		cb.config.packageName = &safe

		logger.Log.Info().
			Str("safe", safe).
			Str("value", *cb.config.packageName).
			Msg("Initialize custom package name")
	}
	if cb.config.fileName != nil {
		safe := fs.AddFileExt(conv.ToSafeName(*cb.config.fileName), ".go")
		cb.config.fileName = &safe

		logger.Log.Info().
			Str("safe", safe).
			Str("value", *cb.config.fileName).
			Msg("Initialize custom filename")
	}

	logger.Log.Debug().Object("data", cb.config).Msg("Generator config initialized")
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

// WithDryRun enables/disables dry run mode.
func (cb *ConfigBuilder) WithDryRun(flag bool) *ConfigBuilder {
	cb.config.dryRun = flag
	return cb
}
