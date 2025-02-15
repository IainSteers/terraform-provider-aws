// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package organizations

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	awstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @FrameworkResource("aws_organizations_account_invite", name="Account Invite")
func newResourceAccountInvite(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceAccountInvite{}

	r.SetDefaultCreateTimeout(30 * time.Minute)
	r.SetDefaultUpdateTimeout(30 * time.Minute)
	r.SetDefaultDeleteTimeout(30 * time.Minute)

	return r, nil
}

const (
	ResNameAccountInvite    = "Account Invite"
	InviteTargetTypeAccount = "ACCOUNT"
	InviteTargetTypeEmail   = "EMAIL"
)

type resourceAccountInvite struct {
	framework.ResourceWithConfigure
	framework.WithTimeouts
}

func (r *resourceAccountInvite) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "aws_organizations_account_invite"
}

func (r *resourceAccountInvite) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": framework.ARNAttributeComputedOnly(),
			"id":  framework.IDAttribute(),
			"notes": schema.StringAttribute{
				Optional: true,
			},
			names.AttrTags:    tftags.TagsAttribute(),
			names.AttrTagsAll: tftags.TagsAttributeComputedOnly(),
			"state": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"target": schema.ListNestedBlock{
				CustomType: fwtypes.NewListNestedObjectTypeOf[inviteTargetModel](ctx),
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(InviteTargetTypeAccount, InviteTargetTypeEmail),
							},
						},
					},
				},
			},
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func (r *resourceAccountInvite) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var plan resourceAccountInviteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var input organizations.InviteAccountToOrganizationInput
	resp.Diagnostics.Append(flex.Expand(ctx, plan, &input)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input.Tags = getTagsIn(ctx)

	out, err := conn.InviteAccountToOrganization(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionCreating, ResNameAccountInvite, *input.Target.Id, err),
			err.Error(),
		)
		return
	}
	if out == nil || out.Handshake == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionCreating, ResNameAccountInvite, *input.Target.Id, nil),
			errors.New("empty output").Error(),
		)
		return
	}

	resp.Diagnostics.Append(flex.Flatten(ctx, out.Handshake, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceAccountInvite) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var state resourceAccountInviteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handshakes are deleted 30 days after they enter a terminal state of ACCEPTED, DECLINED, or CANCELED
	maybeHandshake, err := findMaybeOrganizationsHandshakeByID(ctx, conn, state.ID.ValueString())
	if err != nil && !errs.IsA[*awstypes.HandshakeNotFoundException](err) {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionSetting, ResNameAccountInvite, state.ID.String(), err),
			err.Error(),
		)
		return
	}

	// Update the state if the handshake still exists
	if maybeHandshake.IsSome() {
		resp.Diagnostics.Append(flex.Flatten(ctx, maybeHandshake.MustUnwrap(), &state)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceAccountInvite) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceAccountInviteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceAccountInvite) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var state resourceAccountInviteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := organizations.CancelHandshakeInput{
		HandshakeId: state.ID.ValueStringPointer(),
	}

	_, err := conn.CancelHandshake(ctx, &input)
	if err != nil {
		// Handle cases where the handshake has aged out of the API or already reached a terminal state
		if errs.IsA[*awstypes.HandshakeNotFoundException](err) || errs.IsA[*awstypes.InvalidHandshakeTransitionException](err) {
			return
		}

		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionDeleting, ResNameAccountInvite, state.ID.String(), err),
			err.Error(),
		)
		return
	}
}

func (r *resourceAccountInvite) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *resourceAccountInvite) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	r.SetTagsAll(ctx, request, response)
}

// Try to find a handshake by ID and return an empty optional if none can be found
func findMaybeOrganizationsHandshakeByID(ctx context.Context, conn *organizations.Client, id string) (option.Option[*awstypes.Handshake], error) {
	out, err := findOrganizationsHandshakeByID(ctx, conn, id)
	if tfresource.NotFound(err) {
		return option.None[*awstypes.Handshake](), nil
	}
	if err != nil {
		return nil, err
	}
	return option.Some[*awstypes.Handshake](out), nil
}

func findOrganizationsHandshakeByID(ctx context.Context, conn *organizations.Client, id string) (*awstypes.Handshake, error) {
	in := &organizations.DescribeHandshakeInput{
		HandshakeId: aws.String(id),
	}

	out, err := conn.DescribeHandshake(ctx, in)
	if err != nil {
		if errs.IsA[*awstypes.HandshakeNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Handshake == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out.Handshake, nil
}

type resourceAccountInviteModel struct {
	ARN      types.String                                       `tfsdk:"arn"`
	Target   fwtypes.ListNestedObjectValueOf[inviteTargetModel] `tfsdk:"target"`
	Notes    types.String                                       `tfsdk:"notes"`
	ID       types.String                                       `tfsdk:"id"`
	State    types.String                                       `tfsdk:"state"`
	Tags     tftags.Map                                         `tfsdk:"tags"`
	TagsAll  tftags.Map                                         `tfsdk:"tags_all"`
	Timeouts timeouts.Value                                     `tfsdk:"timeouts"`
}

type inviteTargetModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}
