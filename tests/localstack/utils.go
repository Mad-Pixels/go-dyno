package localstack

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// LocalStackConfig contains configuration for connecting to LocalStack
type LocalStackConfig struct {
	Endpoint string
	Region   string
	Timeout  time.Duration
}

// DefaultLocalStackConfig returns default settings for LocalStack
func DefaultLocalStackConfig() LocalStackConfig {
	return LocalStackConfig{
		Endpoint: "http://localhost:4566",
		Region:   "us-east-1",
		Timeout:  30 * time.Second,
	}
}

// ConnectToLocalStack creates a connection to LocalStack DynamoDB
func ConnectToLocalStack(t *testing.T, cfg LocalStackConfig) *dynamodb.Client {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "test",
				SecretAccessKey: "test",
				SessionToken:    "",
			}, nil
		})),
	)
	require.NoError(t, err, "Failed to load AWS configuration")

	client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})
	_, err = client.ListTables(ctx, &dynamodb.ListTablesInput{})
	require.NoError(t, err, "Failed to connect to LocalStack DynamoDB")
	return client
}

// ConnectToLocalStackWithRetry creates a connection to LocalStack DynamoDB with retry attempts
func ConnectToLocalStackWithRetry(t *testing.T, cfg LocalStackConfig, maxRetries int) *dynamodb.Client {
	t.Helper()

	var client *dynamodb.Client
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		awsCfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "test",
					SecretAccessKey: "test",
					SessionToken:    "",
				}, nil
			})),
		)
		if err != nil {
			cancel()
			lastErr = fmt.Errorf("failed to load AWS configuration: %w", err)
			continue
		}

		client = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
		_, err = client.ListTables(ctx, &dynamodb.ListTablesInput{})
		cancel()
		if err == nil {
			return client
		}

		lastErr = err
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	require.NoError(t, lastErr, "Failed to connect to LocalStack DynamoDB after %d attempts", maxRetries)
	return client
}

// TestContext creates a context with timeout for tests
func TestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// GeneratedPackage contains information about a generated package
type GeneratedPackage struct {
	Name string // Package name (e.g., "test_blog_posts")
	Path string // Package path (e.g., "./generated/test_blog_posts")
}

// GetGeneratedPackages returns a list of all generated packages for testing
func GetGeneratedPackages() []GeneratedPackage {
	return []GeneratedPackage{
		{
			Name: "test_blog_posts",
			Path: "./generated/test_blog_posts",
		},
		{
			Name: "test_table_simple",
			Path: "./generated/test_table_simple",
		},
	}
}

// WaitForLocalStack waits until LocalStack becomes available
func WaitForLocalStack(t *testing.T, cfg LocalStackConfig, maxWaitTime time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			require.Fail(t, "LocalStack did not become available within %v", maxWaitTime)
			return
		case <-ticker.C:
			awsCfg, err := config.LoadDefaultConfig(context.Background(),
				config.WithRegion(cfg.Region),
				config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     "test",
						SecretAccessKey: "test",
					}, nil
				})),
			)
			if err != nil {
				continue
			}

			client := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(cfg.Endpoint)
			})
			if _, err := client.ListTables(context.Background(), &dynamodb.ListTablesInput{}); err == nil {
				return
			}
		}
	}
}
