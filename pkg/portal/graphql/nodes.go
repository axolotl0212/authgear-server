package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

type NodeResolver func(ctx context.Context, id string) (interface{}, error)

var resolvers = map[string]NodeResolver{}
var nodeTypes = map[reflect.Type]*graphql.Object{}

var nodeDefs = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
	IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
		resolvedID := relay.FromGlobalID(id)
		if resolvedID == nil {
			return nil, errors.New("invalid ID")
		}
		resolver, ok := resolvers[resolvedID.Type]
		if !ok {
			return nil, fmt.Errorf("unknown node type: %s", resolvedID.Type)
		}
		return resolver(ctx, resolvedID.ID)
	},
	TypeResolve: func(params graphql.ResolveTypeParams) *graphql.Object {
		objType, ok := nodeTypes[reflect.TypeOf(params.Value)]
		if !ok {
			panic(fmt.Sprintf("graphql: unknown value type: %T", params.Value))
		}
		return objType
	},
})

func node(schema *graphql.Object, modelType interface{}, resolver NodeResolver) *graphql.Object {
	resolvers[schema.Name()] = resolver
	nodeTypes[reflect.TypeOf(modelType)] = schema
	return schema
}
