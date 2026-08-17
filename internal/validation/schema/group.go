package schema

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

type GroupCreateInput struct {
	Name string `form:"name" json:"name" validate:"required,min=1,max=100"`
	Type string `form:"type" json:"type" validate:"required,oneof=direct group"`
}

func SanitizeGroupCreate(input GroupCreateInput) GroupCreateInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.TrimSpace(strings.ToLower(input.Type))
	return input
}

func ValidateGroupCreate(input GroupCreateInput) validation.FieldErrors {
	return validation.ValidateStruct(input)
}

type GroupUpdateInput struct {
	Name        string `form:"name" json:"name" validate:"omitempty,min=1,max=100"`
	Description string `form:"description" json:"description" validate:"omitempty,max=500"`
}

func SanitizeGroupUpdate(input GroupUpdateInput) GroupUpdateInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	return input
}

func ValidateGroupUpdate(input GroupUpdateInput) validation.FieldErrors {
	return validation.ValidateStruct(input)
}

type GroupAddMemberInput struct {
	UserID string `form:"userID" json:"userID" validate:"required"`
}

func SanitizeGroupAddMember(input GroupAddMemberInput) GroupAddMemberInput {
	input.UserID = strings.TrimSpace(input.UserID)
	return input
}

func ValidateGroupAddMember(input GroupAddMemberInput) validation.FieldErrors {
	return validation.ValidateStruct(input)
}
