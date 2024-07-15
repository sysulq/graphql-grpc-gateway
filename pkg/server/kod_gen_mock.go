// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/server/kod_gen_interface.go
//
// Generated by this command:
//
//	mockgen -source pkg/server/kod_gen_interface.go -destination pkg/server/kod_gen_mock.go -package server
//

// Package server is a generated GoMock package.
package server

import (
	context "context"
	http "net/http"
	reflect "reflect"

	desc "github.com/jhump/protoreflect/desc"
	graphql "github.com/nautilus/graphql"
	generator "github.com/sysulq/graphql-gateway/pkg/generator"
	ast "github.com/vektah/gqlparser/v2/ast"
	gomock "go.uber.org/mock/gomock"
	protoadapt "google.golang.org/protobuf/protoadapt"
)

// MockConfigComponent is a mock of ConfigComponent interface.
type MockConfigComponent struct {
	ctrl     *gomock.Controller
	recorder *MockConfigComponentMockRecorder
}

// MockConfigComponentMockRecorder is the mock recorder for MockConfigComponent.
type MockConfigComponentMockRecorder struct {
	mock *MockConfigComponent
}

// NewMockConfigComponent creates a new mock instance.
func NewMockConfigComponent(ctrl *gomock.Controller) *MockConfigComponent {
	mock := &MockConfigComponent{ctrl: ctrl}
	mock.recorder = &MockConfigComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigComponent) EXPECT() *MockConfigComponentMockRecorder {
	return m.recorder
}

// Config mocks base method.
func (m *MockConfigComponent) Config() *Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*Config)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockConfigComponentMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockConfigComponent)(nil).Config))
}

// MockServerComponent is a mock of ServerComponent interface.
type MockServerComponent struct {
	ctrl     *gomock.Controller
	recorder *MockServerComponentMockRecorder
}

// MockServerComponentMockRecorder is the mock recorder for MockServerComponent.
type MockServerComponentMockRecorder struct {
	mock *MockServerComponent
}

// NewMockServerComponent creates a new mock instance.
func NewMockServerComponent(ctrl *gomock.Controller) *MockServerComponent {
	mock := &MockServerComponent{ctrl: ctrl}
	mock.recorder = &MockServerComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerComponent) EXPECT() *MockServerComponentMockRecorder {
	return m.recorder
}

// BuildServer mocks base method.
func (m *MockServerComponent) BuildServer() (http.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildServer")
	ret0, _ := ret[0].(http.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildServer indicates an expected call of BuildServer.
func (mr *MockServerComponentMockRecorder) BuildServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildServer", reflect.TypeOf((*MockServerComponent)(nil).BuildServer))
}

// MockCaller is a mock of Caller interface.
type MockCaller struct {
	ctrl     *gomock.Controller
	recorder *MockCallerMockRecorder
}

// MockCallerMockRecorder is the mock recorder for MockCaller.
type MockCallerMockRecorder struct {
	mock *MockCaller
}

// NewMockCaller creates a new mock instance.
func NewMockCaller(ctrl *gomock.Controller) *MockCaller {
	mock := &MockCaller{ctrl: ctrl}
	mock.recorder = &MockCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCaller) EXPECT() *MockCallerMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockCaller) Call(ctx context.Context, rpc *desc.MethodDescriptor, message protoadapt.MessageV1) (protoadapt.MessageV1, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, rpc, message)
	ret0, _ := ret[0].(protoadapt.MessageV1)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call.
func (mr *MockCallerMockRecorder) Call(ctx, rpc, message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockCaller)(nil).Call), ctx, rpc, message)
}

// GetDescs mocks base method.
func (m *MockCaller) GetDescs() []*desc.FileDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDescs")
	ret0, _ := ret[0].([]*desc.FileDescriptor)
	return ret0
}

// GetDescs indicates an expected call of GetDescs.
func (mr *MockCallerMockRecorder) GetDescs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescs", reflect.TypeOf((*MockCaller)(nil).GetDescs))
}

// MockQueryer is a mock of Queryer interface.
type MockQueryer struct {
	ctrl     *gomock.Controller
	recorder *MockQueryerMockRecorder
}

// MockQueryerMockRecorder is the mock recorder for MockQueryer.
type MockQueryerMockRecorder struct {
	mock *MockQueryer
}

// NewMockQueryer creates a new mock instance.
func NewMockQueryer(ctrl *gomock.Controller) *MockQueryer {
	mock := &MockQueryer{ctrl: ctrl}
	mock.recorder = &MockQueryerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryer) EXPECT() *MockQueryerMockRecorder {
	return m.recorder
}

// Query mocks base method.
func (m *MockQueryer) Query(ctx context.Context, input *graphql.QueryInput, result any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", ctx, input, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockQueryerMockRecorder) Query(ctx, input, result any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockQueryer)(nil).Query), ctx, input, result)
}

// MockRegistry is a mock of Registry interface.
type MockRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockRegistryMockRecorder
}

// MockRegistryMockRecorder is the mock recorder for MockRegistry.
type MockRegistryMockRecorder struct {
	mock *MockRegistry
}

// NewMockRegistry creates a new mock instance.
func NewMockRegistry(ctrl *gomock.Controller) *MockRegistry {
	mock := &MockRegistry{ctrl: ctrl}
	mock.recorder = &MockRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegistry) EXPECT() *MockRegistryMockRecorder {
	return m.recorder
}

// FindFieldByName mocks base method.
func (m *MockRegistry) FindFieldByName(msg desc.Descriptor, name string) *desc.FieldDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindFieldByName", msg, name)
	ret0, _ := ret[0].(*desc.FieldDescriptor)
	return ret0
}

// FindFieldByName indicates an expected call of FindFieldByName.
func (mr *MockRegistryMockRecorder) FindFieldByName(msg, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindFieldByName", reflect.TypeOf((*MockRegistry)(nil).FindFieldByName), msg, name)
}

// FindGraphqlFieldByProtoField mocks base method.
func (m *MockRegistry) FindGraphqlFieldByProtoField(msg *ast.Definition, name string) *ast.FieldDefinition {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindGraphqlFieldByProtoField", msg, name)
	ret0, _ := ret[0].(*ast.FieldDefinition)
	return ret0
}

// FindGraphqlFieldByProtoField indicates an expected call of FindGraphqlFieldByProtoField.
func (mr *MockRegistryMockRecorder) FindGraphqlFieldByProtoField(msg, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindGraphqlFieldByProtoField", reflect.TypeOf((*MockRegistry)(nil).FindGraphqlFieldByProtoField), msg, name)
}

// FindMethodByName mocks base method.
func (m *MockRegistry) FindMethodByName(op ast.Operation, name string) *desc.MethodDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMethodByName", op, name)
	ret0, _ := ret[0].(*desc.MethodDescriptor)
	return ret0
}

// FindMethodByName indicates an expected call of FindMethodByName.
func (mr *MockRegistryMockRecorder) FindMethodByName(op, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMethodByName", reflect.TypeOf((*MockRegistry)(nil).FindMethodByName), op, name)
}

// FindObjectByFullyQualifiedName mocks base method.
func (m *MockRegistry) FindObjectByFullyQualifiedName(fqn string) (*desc.MessageDescriptor, *ast.Definition) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindObjectByFullyQualifiedName", fqn)
	ret0, _ := ret[0].(*desc.MessageDescriptor)
	ret1, _ := ret[1].(*ast.Definition)
	return ret0, ret1
}

// FindObjectByFullyQualifiedName indicates an expected call of FindObjectByFullyQualifiedName.
func (mr *MockRegistryMockRecorder) FindObjectByFullyQualifiedName(fqn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindObjectByFullyQualifiedName", reflect.TypeOf((*MockRegistry)(nil).FindObjectByFullyQualifiedName), fqn)
}

// FindObjectByName mocks base method.
func (m *MockRegistry) FindObjectByName(name string) *desc.MessageDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindObjectByName", name)
	ret0, _ := ret[0].(*desc.MessageDescriptor)
	return ret0
}

// FindObjectByName indicates an expected call of FindObjectByName.
func (mr *MockRegistryMockRecorder) FindObjectByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindObjectByName", reflect.TypeOf((*MockRegistry)(nil).FindObjectByName), name)
}

// FindUnionFieldByMessageFQNAndName mocks base method.
func (m *MockRegistry) FindUnionFieldByMessageFQNAndName(fqn, name string) *desc.FieldDescriptor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUnionFieldByMessageFQNAndName", fqn, name)
	ret0, _ := ret[0].(*desc.FieldDescriptor)
	return ret0
}

// FindUnionFieldByMessageFQNAndName indicates an expected call of FindUnionFieldByMessageFQNAndName.
func (mr *MockRegistryMockRecorder) FindUnionFieldByMessageFQNAndName(fqn, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUnionFieldByMessageFQNAndName", reflect.TypeOf((*MockRegistry)(nil).FindUnionFieldByMessageFQNAndName), fqn, name)
}

// SchemaDescriptorList mocks base method.
func (m *MockRegistry) SchemaDescriptorList() generator.SchemaDescriptorList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SchemaDescriptorList")
	ret0, _ := ret[0].(generator.SchemaDescriptorList)
	return ret0
}

// SchemaDescriptorList indicates an expected call of SchemaDescriptorList.
func (mr *MockRegistryMockRecorder) SchemaDescriptorList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SchemaDescriptorList", reflect.TypeOf((*MockRegistry)(nil).SchemaDescriptorList))
}
