// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package appmesh

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	endpoints_sdkv1 "github.com/aws/aws-sdk-go/aws/endpoints"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	appmesh_sdkv1 "github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceGatewayRoute,
			TypeName: "aws_appmesh_gateway_route",
			Name:     "Gateway Route",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceMesh,
			TypeName: "aws_appmesh_mesh",
			Name:     "Service Mesh",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceRoute,
			TypeName: "aws_appmesh_route",
			Name:     "Route",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceVirtualGateway,
			TypeName: "aws_appmesh_virtual_gateway",
			Name:     "Virtual Gateway",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceVirtualNode,
			TypeName: "aws_appmesh_virtual_node",
			Name:     "Virtual Node",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceVirtualRouter,
			TypeName: "aws_appmesh_virtual_router",
			Name:     "Virtual Router",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceVirtualService,
			TypeName: "aws_appmesh_virtual_service",
			Name:     "Virtual Service",
			Tags:     &types.ServicePackageResourceTags{},
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceGatewayRoute,
			TypeName: "aws_appmesh_gateway_route",
			Name:     "Gateway Route",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceMesh,
			TypeName: "aws_appmesh_mesh",
			Name:     "Service Mesh",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceRoute,
			TypeName: "aws_appmesh_route",
			Name:     "Route",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceVirtualGateway,
			TypeName: "aws_appmesh_virtual_gateway",
			Name:     "Virtual Gateway",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceVirtualNode,
			TypeName: "aws_appmesh_virtual_node",
			Name:     "Virtual Node",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceVirtualRouter,
			TypeName: "aws_appmesh_virtual_router",
			Name:     "Virtual Router",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceVirtualService,
			TypeName: "aws_appmesh_virtual_service",
			Name:     "Virtual Service",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.AppMesh
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*appmesh_sdkv1.AppMesh, error) {
	sess := config[names.AttrSession].(*session_sdkv1.Session)

	cfg := aws_sdkv1.Config{}

	if endpoint := config[names.AttrEndpoint].(string); endpoint != "" {
		tflog.Debug(ctx, "setting endpoint", map[string]any{
			"tf_aws.endpoint": endpoint,
		})
		cfg.Endpoint = aws_sdkv1.String(endpoint)

		if sess.Config.UseFIPSEndpoint == endpoints_sdkv1.FIPSEndpointStateEnabled {
			tflog.Debug(ctx, "endpoint set, ignoring UseFIPSEndpoint setting")
			cfg.UseFIPSEndpoint = endpoints_sdkv1.FIPSEndpointStateDisabled
		}
	}

	return appmesh_sdkv1.New(sess.Copy(&cfg)), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
