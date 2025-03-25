// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package organizations

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	awstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwvalidators "github.com/hashicorp/terraform-provider-aws/internal/framework/validators"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_organizations_account_parent", name="Account Parent")
func newResourceAccountParent(_ context.Context) (resource.ResourceWithConfigure, error) {
	r := &resourceAccountParent{}

	return r, nil
}

const (
	ResNameAccountParent = "Account Parent"
)

type resourceAccountParent struct {
	framework.ResourceWithConfigure
}

func (r *resourceAccountParent) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			names.AttrAccountID: schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					fwvalidators.AWSAccountID(),
				},
			},
			names.AttrParentID: schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.Any(
						fwvalidators.AWSOrganizationRootID(),
						fwvalidators.AWSOrganizationOUID(),
					),
				},
			},
		},
	}
}

func (r *resourceAccountParent) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var plan resourceAccountParentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentParentAccountID, err := FindParentAccountID(ctx, conn, plan.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionCreating, ResNameAccountParent, plan.AccountID.String(), err),
			err.Error(),
		)
	}

	input := organizations.MoveAccountInput{
		AccountId:           flex.StringFromFramework(ctx, plan.AccountID),
		SourceParentId:      currentParentAccountID,
		DestinationParentId: flex.StringFromFramework(ctx, plan.ParentID),
	}

	_, err = conn.MoveAccount(ctx, &input)
	if err != nil && !errs.IsA[*awstypes.DuplicateAccountException](err) {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionCreating, ResNameAccountParent, plan.AccountID.String(), err),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceAccountParent) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var state resourceAccountParentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current parent ID
	parentID, err := FindParentAccountID(ctx, conn, state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionReading, ResNameAccountParent, state.AccountID.String(), err),
			err.Error(),
		)
		return
	}

	// Set the attributes
	state.ParentID = types.StringPointerValue(parentID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceAccountParent) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Meta().OrganizationsClient(ctx)

	var plan, state resourceAccountParentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diff, d := flex.Diff(ctx, plan, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diff.HasChanges() {
		input := organizations.MoveAccountInput{
			AccountId:           flex.StringFromFramework(ctx, plan.AccountID),
			SourceParentId:      flex.StringFromFramework(ctx, state.ParentID),
			DestinationParentId: flex.StringFromFramework(ctx, plan.ParentID),
		}

		_, err := conn.MoveAccount(ctx, &input)
		if err != nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.Organizations, create.ErrActionUpdating, ResNameAccountParent, plan.AccountID.String(), err),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceAccountParent) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceAccountParentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.Meta().OrganizationsClient(ctx)

	roots, err := findRoots(ctx, conn, &organizations.ListRootsInput{})
	if err != nil {
		create.ProblemStandardMessage(names.Organizations, create.ErrActionDeleting, ResNameAccountParent, data.AccountID.String(), err)
	}

	input := organizations.MoveAccountInput{
		AccountId:           flex.StringFromFramework(ctx, data.AccountID),
		SourceParentId:      flex.StringFromFramework(ctx, data.ParentID),
		DestinationParentId: roots[0].Id,
	}

	_, err = conn.MoveAccount(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Organizations, create.ErrActionUpdating, ResNameAccountParent, data.AccountID.String(), err),
			err.Error(),
		)
		return
	}
}

func (r *resourceAccountParent) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root(names.AttrAccountID), req, resp)
}

type resourceAccountParentModel struct {
	AccountID types.String `tfsdk:"account_id"`
	ParentID  types.String `tfsdk:"parent_id"`
}
