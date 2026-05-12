package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := &trinoProvider{
		version: "test",
		name:    "trino",
	}

	resp := &provider.MetadataResponse{}

	p.Metadata(
		context.Background(),
		provider.MetadataRequest{},
		resp,
	)

	if resp.TypeName != "trino" {
		t.Fatalf("expected type name trino, got %q", resp.TypeName)
	}

	if resp.Version != "test" {
		t.Fatalf("expected version test, got %q", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	p := &trinoProvider{}

	resp := &provider.SchemaResponse{}

	p.Schema(
		context.Background(),
		provider.SchemaRequest{},
		resp,
	)

	attrs := resp.Schema.Attributes

	requiredAttrs := []string{
		"host",
		"port",
		"user",
		"password",
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
		"catalog",
		"schema_name",
		"http_scheme",
		"path_to_pem",
		"file_name_pem",
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
}

func TestProviderSchemaPasswordSensitive(t *testing.T) {
	p := &trinoProvider{}

	resp := &provider.SchemaResponse{}

	p.Schema(
		context.Background(),
		provider.SchemaRequest{},
		resp,
	)

	attr, ok := resp.Schema.Attributes["password"]
	if !ok {
		t.Fatalf("expected password attribute to exist")
	}

	if !attr.IsSensitive() {
		t.Fatalf("expected password attribute to be sensitive")
	}
}

func TestProviderDataSources(t *testing.T) {
	p := &trinoProvider{}

	dataSources := p.DataSources(context.Background())

	if len(dataSources) != 0 {
		t.Fatalf("expected no data sources, got %d", len(dataSources))
	}
}

func TestProviderResources(t *testing.T) {
	p := &trinoProvider{}

	resources := p.Resources(context.Background())

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	res := resources[0]()
	if res == nil {
		t.Fatalf("expected resource, got nil")
	}

	_, ok := res.(*schemaResource)
	if !ok {
		t.Fatalf("expected *schemaResource, got %T", res)
	}
}

func TestNewProvider(t *testing.T) {
	factory := New("test", "trino")
	p := factory()

	if p == nil {
		t.Fatalf("expected provider, got nil")
	}

	trinoP, ok := p.(*trinoProvider)
	if !ok {
		t.Fatalf("expected *trinoProvider, got %T", p)
	}

	if trinoP.version != "test" {
		t.Fatalf("expected version test, got %q", trinoP.version)
	}
}
