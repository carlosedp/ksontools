package pipeline

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/bryanl/woowoo/component"
	cmocks "github.com/bryanl/woowoo/component/mocks"
	"github.com/bryanl/woowoo/pipeline/mocks"
	appmocks "github.com/ksonnet/ksonnet/metadata/app/mocks"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPipeline_Namespaces(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {
		namespaces := []component.Namespace{}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)

		got, err := p.Namespaces()
		require.NoError(t, err)

		require.Equal(t, namespaces, got)
	})
}

func TestPipeline_EnvParameters(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {
		ns := component.NewNamespace(p.app, "/")
		namespaces := []component.Namespace{ns}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)
		c.On("Namespace", p.app, "/").Return(ns, nil)
		c.On("NSResolveParams", ns).Return("", nil)
		c.On("EnvParams", p.app, "default").Return("{}", nil)

		got, err := p.EnvParameters("/")
		require.NoError(t, err)

		require.Equal(t, "{ }\n", got)
	})
}

func TestPipeline_Components(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {
		cpnt := &cmocks.Component{}
		components := []component.Component{cpnt}

		ns := component.NewNamespace(p.app, "/")
		namespaces := []component.Namespace{ns}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)
		c.On("Namespace", p.app, "/").Return(ns, nil)
		c.On("NSResolveParams", ns).Return("", nil)
		c.On("EnvParams", p.app, "default").Return("{}", nil)
		c.On("Components", ns).Return(components, nil)

		got, err := p.Components(nil)
		require.NoError(t, err)

		require.Equal(t, components, got)
	})
}

func mockComponent(name string) *cmocks.Component {
	c := &cmocks.Component{}
	c.On("Name", true).Return(name)
	return c
}

func TestPipeline_Components_filtered(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {

		cpnt1 := mockComponent("cpnt1")
		cpnt2 := mockComponent("cpnt2")
		components := []component.Component{cpnt1, cpnt2}

		ns := component.NewNamespace(p.app, "/")
		namespaces := []component.Namespace{ns}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)
		c.On("Namespace", p.app, "/").Return(ns, nil)
		c.On("NSResolveParams", ns).Return("", nil)
		c.On("EnvParams", p.app, "default").Return("{}", nil)
		c.On("Components", ns).Return(components, nil)

		got, err := p.Components([]string{"cpnt1"})
		require.NoError(t, err)

		expected := []component.Component{cpnt1}

		require.Equal(t, expected, got)
	})
}

func TestPipeline_Objects(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {
		u := []*unstructured.Unstructured{
			{},
		}

		cpnt := &cmocks.Component{}
		cpnt.On("Objects", mock.Anything, "default").Return(u, nil)
		components := []component.Component{cpnt}

		ns := component.NewNamespace(p.app, "/")
		namespaces := []component.Namespace{ns}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)
		c.On("Namespace", p.app, "/").Return(ns, nil)
		c.On("NSResolveParams", ns).Return("", nil)
		c.On("EnvParams", p.app, "default").Return("{}", nil)
		c.On("Components", ns).Return(components, nil)

		got, err := p.Objects(nil)
		require.NoError(t, err)

		require.Equal(t, u, got)
	})
}

func TestPipeline_YAML(t *testing.T) {
	withPipeline(t, func(p *Pipeline, c *mocks.Component) {
		u := []*unstructured.Unstructured{
			{},
		}

		cpnt := &cmocks.Component{}
		cpnt.On("Objects", mock.Anything, "default").Return(u, nil)
		components := []component.Component{cpnt}

		ns := component.NewNamespace(p.app, "/")
		namespaces := []component.Namespace{ns}
		c.On("Namespaces", p.app, "default").Return(namespaces, nil)
		c.On("Namespace", p.app, "/").Return(ns, nil)
		c.On("NSResolveParams", ns).Return("", nil)
		c.On("EnvParams", p.app, "default").Return("{}", nil)
		c.On("Components", ns).Return(components, nil)

		r, err := p.YAML(nil)
		require.NoError(t, err)

		got, err := ioutil.ReadAll(r)
		require.NoError(t, err)

		expected := "---\n{}\n"

		require.Equal(t, expected, string(got))
	})
}

func Test_upgradeParams(t *testing.T) {
	in := `local params = import "../../components/params.libsonnet";`
	expected := `local params = std.extVar("__ksonnet/params");`

	got := upgradeParams("default", in)
	require.Equal(t, expected, got)
}

func withPipeline(t *testing.T, fn func(p *Pipeline, c *mocks.Component)) {
	app := &appmocks.App{}
	envName := "default"

	c := &mocks.Component{}

	p := New(app, envName, OverrideComponent(c))

	fn(p, c)
}
