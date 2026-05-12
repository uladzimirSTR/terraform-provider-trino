package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestNewSchemaResource(t *testing.T) {
	res := NewSchemaResource()

	if res == nil {
		t.Fatalf("expected resource, got nil")
	}

	_, ok := res.(*schemaResource)
	if !ok {
		t.Fatalf("expected *schemaResource, got %T", res)
	}
}

func TestSchemaResourceMetadata(t *testing.T) {
	res := &schemaResource{}

	resp := &resource.MetadataResponse{}

	res.Metadata(
		context.Background(),
		resource.MetadataRequest{
			ProviderTypeName: "trino",
		},
		resp,
	)

	if resp.TypeName != "trino_schema" {
		t.Fatalf("expected type name trino_schema, got %q", resp.TypeName)
	}
}

func TestSchemaResourceSchema(t *testing.T) {
	res := &schemaResource{}

	resp := &resource.SchemaResponse{}

	res.Schema(
		context.Background(),
		resource.SchemaRequest{},
		resp,
	)

	attrs := resp.Schema.Attributes

	requiredAttrs := []string{
		"catalog",
		"name",
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

func TestSchemaResourceConfigureNilProviderData(t *testing.T) {
	res := &schemaResource{}

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

func TestSchemaResourceConfigureUnexpectedProviderDataType(t *testing.T) {
	res := &schemaResource{}

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
