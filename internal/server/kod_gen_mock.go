//go:build !ignoreKodGen

// Code generated by MockGen. DO NOT EDIT.
// Source: internal/server/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source internal/server/kod_gen_interface.go -destination internal/server/kod_gen_mock.go -package server -typed -build_constraint !ignoreKodGen
//

// Package server is a generated GoMock package.
package server

import (
	context "context"
	http "net/http"
	reflect "reflect"

	grpcdynamic "github.com/jhump/protoreflect/v2/grpcdynamic"
	graphql "github.com/nautilus/graphql"
	ast "github.com/vektah/gqlparser/v2/ast"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// MockGateway is a mock of Gateway interface.
type MockGateway struct {
	ctrl     *gomock.Controller
	recorder *MockGatewayMockRecorder
	isgomock struct{}
}

// MockGatewayMockRecorder is the mock recorder for MockGateway.
type MockGatewayMockRecorder struct {
	mock *MockGateway
}

// NewMockGateway creates a new mock instance.
func NewMockGateway(ctrl *gomock.Controller) *MockGateway {
	mock := &MockGateway{ctrl: ctrl}
	mock.recorder = &MockGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGateway) EXPECT() *MockGatewayMockRecorder {
	return m.recorder
}

// BuildHTTPServer mocks base method.
func (m *MockGateway) BuildHTTPServer() (http.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildHTTPServer")
	ret0, _ := ret[0].(http.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildHTTPServer indicates an expected call of BuildHTTPServer.
func (mr *MockGatewayMockRecorder) BuildHTTPServer() *MockGatewayBuildHTTPServerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildHTTPServer", reflect.TypeOf((*MockGateway)(nil).BuildHTTPServer))
	return &MockGatewayBuildHTTPServerCall{Call: call}
}

// MockGatewayBuildHTTPServerCall wrap *gomock.Call
type MockGatewayBuildHTTPServerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGatewayBuildHTTPServerCall) Return(arg0 http.Handler, arg1 error) *MockGatewayBuildHTTPServerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGatewayBuildHTTPServerCall) Do(f func() (http.Handler, error)) *MockGatewayBuildHTTPServerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGatewayBuildHTTPServerCall) DoAndReturn(f func() (http.Handler, error)) *MockGatewayBuildHTTPServerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BuildServer mocks base method.
func (m *MockGateway) BuildServer() (http.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildServer")
	ret0, _ := ret[0].(http.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildServer indicates an expected call of BuildServer.
func (mr *MockGatewayMockRecorder) BuildServer() *MockGatewayBuildServerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildServer", reflect.TypeOf((*MockGateway)(nil).BuildServer))
	return &MockGatewayBuildServerCall{Call: call}
}

// MockGatewayBuildServerCall wrap *gomock.Call
type MockGatewayBuildServerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGatewayBuildServerCall) Return(arg0 http.Handler, arg1 error) *MockGatewayBuildServerCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGatewayBuildServerCall) Do(f func() (http.Handler, error)) *MockGatewayBuildServerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGatewayBuildServerCall) DoAndReturn(f func() (http.Handler, error)) *MockGatewayBuildServerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockGraphqlCaller is a mock of GraphqlCaller interface.
type MockGraphqlCaller struct {
	ctrl     *gomock.Controller
	recorder *MockGraphqlCallerMockRecorder
	isgomock struct{}
}

// MockGraphqlCallerMockRecorder is the mock recorder for MockGraphqlCaller.
type MockGraphqlCallerMockRecorder struct {
	mock *MockGraphqlCaller
}

// NewMockGraphqlCaller creates a new mock instance.
func NewMockGraphqlCaller(ctrl *gomock.Controller) *MockGraphqlCaller {
	mock := &MockGraphqlCaller{ctrl: ctrl}
	mock.recorder = &MockGraphqlCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphqlCaller) EXPECT() *MockGraphqlCallerMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockGraphqlCaller) Call(ctx context.Context, rpc protoreflect.MethodDescriptor, message proto.Message) (proto.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, rpc, message)
	ret0, _ := ret[0].(proto.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call.
func (mr *MockGraphqlCallerMockRecorder) Call(ctx, rpc, message any) *MockGraphqlCallerCallCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockGraphqlCaller)(nil).Call), ctx, rpc, message)
	return &MockGraphqlCallerCallCall{Call: call}
}

// MockGraphqlCallerCallCall wrap *gomock.Call
type MockGraphqlCallerCallCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerCallCall) Return(arg0 proto.Message, arg1 error) *MockGraphqlCallerCallCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerCallCall) Do(f func(context.Context, protoreflect.MethodDescriptor, proto.Message) (proto.Message, error)) *MockGraphqlCallerCallCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerCallCall) DoAndReturn(f func(context.Context, protoreflect.MethodDescriptor, proto.Message) (proto.Message, error)) *MockGraphqlCallerCallCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockGraphqlCallerRegistry is a mock of GraphqlCallerRegistry interface.
type MockGraphqlCallerRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockGraphqlCallerRegistryMockRecorder
	isgomock struct{}
}

// MockGraphqlCallerRegistryMockRecorder is the mock recorder for MockGraphqlCallerRegistry.
type MockGraphqlCallerRegistryMockRecorder struct {
	mock *MockGraphqlCallerRegistry
}

// NewMockGraphqlCallerRegistry creates a new mock instance.
func NewMockGraphqlCallerRegistry(ctrl *gomock.Controller) *MockGraphqlCallerRegistry {
	mock := &MockGraphqlCallerRegistry{ctrl: ctrl}
	mock.recorder = &MockGraphqlCallerRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphqlCallerRegistry) EXPECT() *MockGraphqlCallerRegistryMockRecorder {
	return m.recorder
}

// FindMethodByName mocks base method.
func (m *MockGraphqlCallerRegistry) FindMethodByName(op ast.Operation, name string) protoreflect.MethodDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMethodByName", op, name)
	ret0, _ := ret[0].(protoreflect.MethodDescriptor)
	return ret0
}

// FindMethodByName indicates an expected call of FindMethodByName.
func (mr *MockGraphqlCallerRegistryMockRecorder) FindMethodByName(op, name any) *MockGraphqlCallerRegistryFindMethodByNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMethodByName", reflect.TypeOf((*MockGraphqlCallerRegistry)(nil).FindMethodByName), op, name)
	return &MockGraphqlCallerRegistryFindMethodByNameCall{Call: call}
}

// MockGraphqlCallerRegistryFindMethodByNameCall wrap *gomock.Call
type MockGraphqlCallerRegistryFindMethodByNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerRegistryFindMethodByNameCall) Return(arg0 protoreflect.MethodDescriptor) *MockGraphqlCallerRegistryFindMethodByNameCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerRegistryFindMethodByNameCall) Do(f func(ast.Operation, string) protoreflect.MethodDescriptor) *MockGraphqlCallerRegistryFindMethodByNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerRegistryFindMethodByNameCall) DoAndReturn(f func(ast.Operation, string) protoreflect.MethodDescriptor) *MockGraphqlCallerRegistryFindMethodByNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetCallerStub mocks base method.
func (m *MockGraphqlCallerRegistry) GetCallerStub(service string) *grpcdynamic.Stub {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCallerStub", service)
	ret0, _ := ret[0].(*grpcdynamic.Stub)
	return ret0
}

// GetCallerStub indicates an expected call of GetCallerStub.
func (mr *MockGraphqlCallerRegistryMockRecorder) GetCallerStub(service any) *MockGraphqlCallerRegistryGetCallerStubCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCallerStub", reflect.TypeOf((*MockGraphqlCallerRegistry)(nil).GetCallerStub), service)
	return &MockGraphqlCallerRegistryGetCallerStubCall{Call: call}
}

// MockGraphqlCallerRegistryGetCallerStubCall wrap *gomock.Call
type MockGraphqlCallerRegistryGetCallerStubCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerRegistryGetCallerStubCall) Return(arg0 *grpcdynamic.Stub) *MockGraphqlCallerRegistryGetCallerStubCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerRegistryGetCallerStubCall) Do(f func(string) *grpcdynamic.Stub) *MockGraphqlCallerRegistryGetCallerStubCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerRegistryGetCallerStubCall) DoAndReturn(f func(string) *grpcdynamic.Stub) *MockGraphqlCallerRegistryGetCallerStubCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GraphQLSchema mocks base method.
func (m *MockGraphqlCallerRegistry) GraphQLSchema() *ast.Schema {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GraphQLSchema")
	ret0, _ := ret[0].(*ast.Schema)
	return ret0
}

// GraphQLSchema indicates an expected call of GraphQLSchema.
func (mr *MockGraphqlCallerRegistryMockRecorder) GraphQLSchema() *MockGraphqlCallerRegistryGraphQLSchemaCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GraphQLSchema", reflect.TypeOf((*MockGraphqlCallerRegistry)(nil).GraphQLSchema))
	return &MockGraphqlCallerRegistryGraphQLSchemaCall{Call: call}
}

// MockGraphqlCallerRegistryGraphQLSchemaCall wrap *gomock.Call
type MockGraphqlCallerRegistryGraphQLSchemaCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerRegistryGraphQLSchemaCall) Return(arg0 *ast.Schema) *MockGraphqlCallerRegistryGraphQLSchemaCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerRegistryGraphQLSchemaCall) Do(f func() *ast.Schema) *MockGraphqlCallerRegistryGraphQLSchemaCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerRegistryGraphQLSchemaCall) DoAndReturn(f func() *ast.Schema) *MockGraphqlCallerRegistryGraphQLSchemaCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Marshal mocks base method.
func (m *MockGraphqlCallerRegistry) Marshal(proto proto.Message, field *ast.Field) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", proto, field)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockGraphqlCallerRegistryMockRecorder) Marshal(proto, field any) *MockGraphqlCallerRegistryMarshalCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockGraphqlCallerRegistry)(nil).Marshal), proto, field)
	return &MockGraphqlCallerRegistryMarshalCall{Call: call}
}

// MockGraphqlCallerRegistryMarshalCall wrap *gomock.Call
type MockGraphqlCallerRegistryMarshalCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerRegistryMarshalCall) Return(arg0 any, arg1 error) *MockGraphqlCallerRegistryMarshalCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerRegistryMarshalCall) Do(f func(proto.Message, *ast.Field) (any, error)) *MockGraphqlCallerRegistryMarshalCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerRegistryMarshalCall) DoAndReturn(f func(proto.Message, *ast.Field) (any, error)) *MockGraphqlCallerRegistryMarshalCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Unmarshal mocks base method.
func (m *MockGraphqlCallerRegistry) Unmarshal(desc protoreflect.MessageDescriptor, field *ast.Field, vars map[string]any) (proto.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", desc, field, vars)
	ret0, _ := ret[0].(proto.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockGraphqlCallerRegistryMockRecorder) Unmarshal(desc, field, vars any) *MockGraphqlCallerRegistryUnmarshalCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockGraphqlCallerRegistry)(nil).Unmarshal), desc, field, vars)
	return &MockGraphqlCallerRegistryUnmarshalCall{Call: call}
}

// MockGraphqlCallerRegistryUnmarshalCall wrap *gomock.Call
type MockGraphqlCallerRegistryUnmarshalCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlCallerRegistryUnmarshalCall) Return(arg0 proto.Message, arg1 error) *MockGraphqlCallerRegistryUnmarshalCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlCallerRegistryUnmarshalCall) Do(f func(protoreflect.MessageDescriptor, *ast.Field, map[string]any) (proto.Message, error)) *MockGraphqlCallerRegistryUnmarshalCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlCallerRegistryUnmarshalCall) DoAndReturn(f func(protoreflect.MessageDescriptor, *ast.Field, map[string]any) (proto.Message, error)) *MockGraphqlCallerRegistryUnmarshalCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockGraphqlReflection is a mock of GraphqlReflection interface.
type MockGraphqlReflection struct {
	ctrl     *gomock.Controller
	recorder *MockGraphqlReflectionMockRecorder
	isgomock struct{}
}

// MockGraphqlReflectionMockRecorder is the mock recorder for MockGraphqlReflection.
type MockGraphqlReflectionMockRecorder struct {
	mock *MockGraphqlReflection
}

// NewMockGraphqlReflection creates a new mock instance.
func NewMockGraphqlReflection(ctrl *gomock.Controller) *MockGraphqlReflection {
	mock := &MockGraphqlReflection{ctrl: ctrl}
	mock.recorder = &MockGraphqlReflectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphqlReflection) EXPECT() *MockGraphqlReflectionMockRecorder {
	return m.recorder
}

// ListPackages mocks base method.
func (m *MockGraphqlReflection) ListPackages(ctx context.Context, cc grpc.ClientConnInterface) ([]protoreflect.FileDescriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPackages", ctx, cc)
	ret0, _ := ret[0].([]protoreflect.FileDescriptor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPackages indicates an expected call of ListPackages.
func (mr *MockGraphqlReflectionMockRecorder) ListPackages(ctx, cc any) *MockGraphqlReflectionListPackagesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPackages", reflect.TypeOf((*MockGraphqlReflection)(nil).ListPackages), ctx, cc)
	return &MockGraphqlReflectionListPackagesCall{Call: call}
}

// MockGraphqlReflectionListPackagesCall wrap *gomock.Call
type MockGraphqlReflectionListPackagesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlReflectionListPackagesCall) Return(arg0 []protoreflect.FileDescriptor, arg1 error) *MockGraphqlReflectionListPackagesCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlReflectionListPackagesCall) Do(f func(context.Context, grpc.ClientConnInterface) ([]protoreflect.FileDescriptor, error)) *MockGraphqlReflectionListPackagesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlReflectionListPackagesCall) DoAndReturn(f func(context.Context, grpc.ClientConnInterface) ([]protoreflect.FileDescriptor, error)) *MockGraphqlReflectionListPackagesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockGraphqlQueryer is a mock of GraphqlQueryer interface.
type MockGraphqlQueryer struct {
	ctrl     *gomock.Controller
	recorder *MockGraphqlQueryerMockRecorder
	isgomock struct{}
}

// MockGraphqlQueryerMockRecorder is the mock recorder for MockGraphqlQueryer.
type MockGraphqlQueryerMockRecorder struct {
	mock *MockGraphqlQueryer
}

// NewMockGraphqlQueryer creates a new mock instance.
func NewMockGraphqlQueryer(ctrl *gomock.Controller) *MockGraphqlQueryer {
	mock := &MockGraphqlQueryer{ctrl: ctrl}
	mock.recorder = &MockGraphqlQueryerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraphqlQueryer) EXPECT() *MockGraphqlQueryerMockRecorder {
	return m.recorder
}

// Query mocks base method.
func (m *MockGraphqlQueryer) Query(ctx context.Context, input *graphql.QueryInput, result any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", ctx, input, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockGraphqlQueryerMockRecorder) Query(ctx, input, result any) *MockGraphqlQueryerQueryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockGraphqlQueryer)(nil).Query), ctx, input, result)
	return &MockGraphqlQueryerQueryCall{Call: call}
}

// MockGraphqlQueryerQueryCall wrap *gomock.Call
type MockGraphqlQueryerQueryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockGraphqlQueryerQueryCall) Return(arg0 error) *MockGraphqlQueryerQueryCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockGraphqlQueryerQueryCall) Do(f func(context.Context, *graphql.QueryInput, any) error) *MockGraphqlQueryerQueryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockGraphqlQueryerQueryCall) DoAndReturn(f func(context.Context, *graphql.QueryInput, any) error) *MockGraphqlQueryerQueryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockHttpUpstreamInvoker is a mock of HttpUpstreamInvoker interface.
type MockHttpUpstreamInvoker struct {
	ctrl     *gomock.Controller
	recorder *MockHttpUpstreamInvokerMockRecorder
	isgomock struct{}
}

// MockHttpUpstreamInvokerMockRecorder is the mock recorder for MockHttpUpstreamInvoker.
type MockHttpUpstreamInvokerMockRecorder struct {
	mock *MockHttpUpstreamInvoker
}

// NewMockHttpUpstreamInvoker creates a new mock instance.
func NewMockHttpUpstreamInvoker(ctrl *gomock.Controller) *MockHttpUpstreamInvoker {
	mock := &MockHttpUpstreamInvoker{ctrl: ctrl}
	mock.recorder = &MockHttpUpstreamInvokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpUpstreamInvoker) EXPECT() *MockHttpUpstreamInvokerMockRecorder {
	return m.recorder
}

// Invoke mocks base method.
func (m *MockHttpUpstreamInvoker) Invoke(ctx context.Context, rw http.ResponseWriter, r *http.Request, upstream upstreamInfo, rpcPath string, pathNames []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Invoke", ctx, rw, r, upstream, rpcPath, pathNames)
}

// Invoke indicates an expected call of Invoke.
func (mr *MockHttpUpstreamInvokerMockRecorder) Invoke(ctx, rw, r, upstream, rpcPath, pathNames any) *MockHttpUpstreamInvokerInvokeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockHttpUpstreamInvoker)(nil).Invoke), ctx, rw, r, upstream, rpcPath, pathNames)
	return &MockHttpUpstreamInvokerInvokeCall{Call: call}
}

// MockHttpUpstreamInvokerInvokeCall wrap *gomock.Call
type MockHttpUpstreamInvokerInvokeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHttpUpstreamInvokerInvokeCall) Return() *MockHttpUpstreamInvokerInvokeCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHttpUpstreamInvokerInvokeCall) Do(f func(context.Context, http.ResponseWriter, *http.Request, upstreamInfo, string, []string)) *MockHttpUpstreamInvokerInvokeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHttpUpstreamInvokerInvokeCall) DoAndReturn(f func(context.Context, http.ResponseWriter, *http.Request, upstreamInfo, string, []string)) *MockHttpUpstreamInvokerInvokeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockHttpUpstream is a mock of HttpUpstream interface.
type MockHttpUpstream struct {
	ctrl     *gomock.Controller
	recorder *MockHttpUpstreamMockRecorder
	isgomock struct{}
}

// MockHttpUpstreamMockRecorder is the mock recorder for MockHttpUpstream.
type MockHttpUpstreamMockRecorder struct {
	mock *MockHttpUpstream
}

// NewMockHttpUpstream creates a new mock instance.
func NewMockHttpUpstream(ctrl *gomock.Controller) *MockHttpUpstream {
	mock := &MockHttpUpstream{ctrl: ctrl}
	mock.recorder = &MockHttpUpstreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpUpstream) EXPECT() *MockHttpUpstreamMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockHttpUpstream) Register(ctx context.Context, router *http.ServeMux) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", ctx, router)
}

// Register indicates an expected call of Register.
func (mr *MockHttpUpstreamMockRecorder) Register(ctx, router any) *MockHttpUpstreamRegisterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockHttpUpstream)(nil).Register), ctx, router)
	return &MockHttpUpstreamRegisterCall{Call: call}
}

// MockHttpUpstreamRegisterCall wrap *gomock.Call
type MockHttpUpstreamRegisterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockHttpUpstreamRegisterCall) Return() *MockHttpUpstreamRegisterCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockHttpUpstreamRegisterCall) Do(f func(context.Context, *http.ServeMux)) *MockHttpUpstreamRegisterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockHttpUpstreamRegisterCall) DoAndReturn(f func(context.Context, *http.ServeMux)) *MockHttpUpstreamRegisterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
