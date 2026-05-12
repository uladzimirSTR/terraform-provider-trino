package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewTableResource(t *testing.T) {
	res := NewTableResource()

	if res == nil {
		t.Fatalf("expected resource, got nil")
	}

	_, ok := res.(*tableResource)
	if !ok {
		t.Fatalf("expected *tableResource, got %T", res)
	}
}

func TestTableResourceMetadata(t *testing.T) {
	res := &tableResource{}

	resp := &resource.MetadataResponse{}

	res.Metadata(
		context.Background(),
		resource.MetadataRequest{
			ProviderTypeName: "trino",
		},
		resp,
	)

	if resp.TypeName != "trino_table" {
		t.Fatalf("expected type name trino_table, got %q", resp.TypeName)
	}
}

func TestTableResourceSchema(t *testing.T) {
	res := &tableResource{}

	resp := &resource.SchemaResponse{}

	res.Schema(
		context.Background(),
		resource.SchemaRequest{},
		resp,
	)

	attrs := resp.Schema.Attributes

	requiredAttrs := []string{
		"catalog",
		"schema_name",
		"name",
		"columns",
	}

	for _, name := range requiredAttrs {
		attr, ok := attrs[name]
		if !ok {
			t.Fatalf("expected attribute %q to exist", name)
		}

		if !attr.IsRequired() {
			t.Fatalf("expected attribute %q to be required", name)
		}
	}

	optionalAttrs := []string{
		"location",
		"format",
	}

	for _, name := range optionalAttrs {
		attr, ok := attrs[name]
		if !ok {
			t.Fatalf("expected attribute %q to exist", name)
		}

		if !attr.IsOptional() {
			t.Fatalf("expected attribute %q to be optional", name)
		}
	}

	computedAttrs := []string{
		"id",
	}

	for _, name := range computedAttrs {
		attr, ok := attrs[name]
		if !ok {
			t.Fatalf("expected attribute %q to exist", name)
		}

		if !attr.IsComputed() {
			t.Fatalf("expected attribute %q to be computed", name)
		}
	}
}

func TestTableResourceSchemaColumnsNestedAttributes(t *testing.T) {
	res := &tableResource{}

	resp := &resource.SchemaResponse{}

	res.Schema(
		context.Background(),
		resource.SchemaRequest{},
		resp,
	)

	attr, ok := resp.Schema.Attributes["columns"]
	if !ok {
		t.Fatalf("expected columns attribute to exist")
	}

	if !attr.IsRequired() {
		t.Fatalf("expected columns attribute to be required")
	}
}

func TestTableResourceConfigureNilProviderData(t *testing.T) {
	res := &tableResource{}

	resp := &resource.ConfigureResponse{}

	res.Configure(
		context.Background(),
		resource.ConfigureRequest{
			ProviderData: nil,
		},
		resp,
	)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", resp.Diagnostics)
	}

	if res.client != nil {
		t.Fatalf("expected client to remain nil")
	}
}

func TestTableResourceConfigureUnexpectedProviderDataType(t *testing.T) {
	res := &tableResource{}

	resp := &resource.ConfigureResponse{}

	res.Configure(
		context.Background(),
		resource.ConfigureRequest{
			ProviderData: "not-a-client",
		},
		resp,
	)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics error, got none")
	}

	got := resp.Diagnostics.Errors()[0].Summary()

	if !strings.Contains(got, "Unexpected Provider Data Type") {
		t.Fatalf("expected unexpected provider data type error, got %q", got)
	}
}
