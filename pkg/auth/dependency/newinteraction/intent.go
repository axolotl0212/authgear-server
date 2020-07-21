package newinteraction

import (
	"reflect"
)

type Intent interface {
	InstantiateRootNode(ctx *Context, graph *Graph) (Node, error)
	DeriveEdges(ctx *Context, graph *Graph, node Node) ([]Edge, error)
}

type IntentFactory func() Intent

var intentRegistry = map[string]IntentFactory{}

func RegisterIntent(intent Intent) {
	intentType := reflect.TypeOf(intent).Elem()

	intentKind := intentType.Name()
	factory := IntentFactory(func() Intent {
		return reflect.New(intentType).Interface().(Intent)
	})
	intentRegistry[intentKind] = factory
}

func IntentKind(intent Intent) string {
	intentType := reflect.TypeOf(intent).Elem()
	return intentType.Name()
}

func InstantiateIntent(kind string) Intent {
	factory, ok := intentRegistry[kind]
	if !ok {
		panic("interaction: unknown intent kind: " + kind)
	}
	return factory()
}
