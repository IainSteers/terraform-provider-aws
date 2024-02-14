// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cloudfront

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	cloudfront_sdkv2 "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	cloudfront_sdkv1 "github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newKeyValueStoreResource,
			Name:    "Key Value Store",
		},
		{
			Factory: newResourceContinuousDeploymentPolicy,
			Name:    "Continuous Deployment Policy",
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceCachePolicy,
			TypeName: "aws_cloudfront_cache_policy",
		},
		{
			Factory:  DataSourceDistribution,
			TypeName: "aws_cloudfront_distribution",
		},
		{
			Factory:  DataSourceFunction,
			TypeName: "aws_cloudfront_function",
		},
		{
			Factory:  DataSourceLogDeliveryCanonicalUserID,
			TypeName: "aws_cloudfront_log_delivery_canonical_user_id",
		},
		{
			Factory:  DataSourceOriginAccessIdentities,
			TypeName: "aws_cloudfront_origin_access_identities",
		},
		{
			Factory:  DataSourceOriginAccessIdentity,
			TypeName: "aws_cloudfront_origin_access_identity",
		},
		{
			Factory:  DataSourceOriginRequestPolicy,
			TypeName: "aws_cloudfront_origin_request_policy",
		},
		{
			Factory:  DataSourceRealtimeLogConfig,
			TypeName: "aws_cloudfront_realtime_log_config",
		},
		{
			Factory:  DataSourceResponseHeadersPolicy,
			TypeName: "aws_cloudfront_response_headers_policy",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceCachePolicy,
			TypeName: "aws_cloudfront_cache_policy",
		},
		{
			Factory:  ResourceDistribution,
			TypeName: "aws_cloudfront_distribution",
			Name:     "Distribution",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceFieldLevelEncryptionConfig,
			TypeName: "aws_cloudfront_field_level_encryption_config",
		},
		{
			Factory:  ResourceFieldLevelEncryptionProfile,
			TypeName: "aws_cloudfront_field_level_encryption_profile",
		},
		{
			Factory:  ResourceFunction,
			TypeName: "aws_cloudfront_function",
		},
		{
			Factory:  ResourceKeyGroup,
			TypeName: "aws_cloudfront_key_group",
		},
		{
			Factory:  ResourceMonitoringSubscription,
			TypeName: "aws_cloudfront_monitoring_subscription",
		},
		{
			Factory:  ResourceOriginAccessControl,
			TypeName: "aws_cloudfront_origin_access_control",
		},
		{
			Factory:  ResourceOriginAccessIdentity,
			TypeName: "aws_cloudfront_origin_access_identity",
		},
		{
			Factory:  ResourceOriginRequestPolicy,
			TypeName: "aws_cloudfront_origin_request_policy",
		},
		{
			Factory:  ResourcePublicKey,
			TypeName: "aws_cloudfront_public_key",
		},
		{
			Factory:  ResourceRealtimeLogConfig,
			TypeName: "aws_cloudfront_realtime_log_config",
		},
		{
			Factory:  ResourceResponseHeadersPolicy,
			TypeName: "aws_cloudfront_response_headers_policy",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CloudFront
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*cloudfront_sdkv1.CloudFront, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return cloudfront_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*cloudfront_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return cloudfront_sdkv2.NewFromConfig(cfg, func(o *cloudfront_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
