// Code generated by protoc-gen-go.
// source: jobbo.proto
// DO NOT EDIT!

/*
Package jobbo is a generated protocol buffer package.

It is generated from these files:
	jobbo.proto

It has these top-level messages:
	RunJobRequest
	RunJobReply
	ListInstancesRequest
	ListInstancesReply
	RemoveInstanceRequest
	RemoveInstanceReply
	GetInstanceLogsRequest
	GetInstanceLogsReply
	Instance
	InstanceState
	Job
	Metadata
	Resources
	Container
	Command
*/
package jobbo

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The Jobbo state of a task
type TaskStatus int32

const (
	TaskStatus_UNKNOWN   TaskStatus = 0
	TaskStatus_PENDING   TaskStatus = 1
	TaskStatus_SCHEDULED TaskStatus = 2
	TaskStatus_RUNNING   TaskStatus = 3
	TaskStatus_FAILED    TaskStatus = 4
	TaskStatus_FINISHED  TaskStatus = 5
)

var TaskStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "PENDING",
	2: "SCHEDULED",
	3: "RUNNING",
	4: "FAILED",
	5: "FINISHED",
}
var TaskStatus_value = map[string]int32{
	"UNKNOWN":   0,
	"PENDING":   1,
	"SCHEDULED": 2,
	"RUNNING":   3,
	"FAILED":    4,
	"FINISHED":  5,
}

func (x TaskStatus) String() string {
	return proto.EnumName(TaskStatus_name, int32(x))
}
func (TaskStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The request message containing the job to be created.
type RunJobRequest struct {
	Job *Job `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
}

func (m *RunJobRequest) Reset()                    { *m = RunJobRequest{} }
func (m *RunJobRequest) String() string            { return proto.CompactTextString(m) }
func (*RunJobRequest) ProtoMessage()               {}
func (*RunJobRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RunJobRequest) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

// The response message containing the job that was created (metadata will be updated)
type RunJobReply struct {
	Job *Job `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
}

func (m *RunJobReply) Reset()                    { *m = RunJobReply{} }
func (m *RunJobReply) String() string            { return proto.CompactTextString(m) }
func (*RunJobReply) ProtoMessage()               {}
func (*RunJobReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RunJobReply) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

// The jobs returned may be filtered by passing various properties here
type ListInstancesRequest struct {
	Status TaskStatus `protobuf:"varint,1,opt,name=status,enum=jobbo.TaskStatus" json:"status,omitempty"`
	Name   string     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Uid    string     `protobuf:"bytes,3,opt,name=uid" json:"uid,omitempty"`
}

func (m *ListInstancesRequest) Reset()                    { *m = ListInstancesRequest{} }
func (m *ListInstancesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListInstancesRequest) ProtoMessage()               {}
func (*ListInstancesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type ListInstancesReply struct {
	Instances []*Instance `protobuf:"bytes,1,rep,name=instances" json:"instances,omitempty"`
}

func (m *ListInstancesReply) Reset()                    { *m = ListInstancesReply{} }
func (m *ListInstancesReply) String() string            { return proto.CompactTextString(m) }
func (*ListInstancesReply) ProtoMessage()               {}
func (*ListInstancesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ListInstancesReply) GetInstances() []*Instance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type RemoveInstanceRequest struct {
	Uid    string     `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
	Status TaskStatus `protobuf:"varint,2,opt,name=status,enum=jobbo.TaskStatus" json:"status,omitempty"`
}

func (m *RemoveInstanceRequest) Reset()                    { *m = RemoveInstanceRequest{} }
func (m *RemoveInstanceRequest) String() string            { return proto.CompactTextString(m) }
func (*RemoveInstanceRequest) ProtoMessage()               {}
func (*RemoveInstanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type RemoveInstanceReply struct {
	Instance *Instance `protobuf:"bytes,1,opt,name=instance" json:"instance,omitempty"`
}

func (m *RemoveInstanceReply) Reset()                    { *m = RemoveInstanceReply{} }
func (m *RemoveInstanceReply) String() string            { return proto.CompactTextString(m) }
func (*RemoveInstanceReply) ProtoMessage()               {}
func (*RemoveInstanceReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RemoveInstanceReply) GetInstance() *Instance {
	if m != nil {
		return m.Instance
	}
	return nil
}

type GetInstanceLogsRequest struct {
	Uid  string `protobuf:"bytes,1,opt,name=uid" json:"uid,omitempty"`
	Tail uint32 `protobuf:"varint,2,opt,name=tail" json:"tail,omitempty"`
}

func (m *GetInstanceLogsRequest) Reset()                    { *m = GetInstanceLogsRequest{} }
func (m *GetInstanceLogsRequest) String() string            { return proto.CompactTextString(m) }
func (*GetInstanceLogsRequest) ProtoMessage()               {}
func (*GetInstanceLogsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type GetInstanceLogsReply struct {
	Stderr []byte `protobuf:"bytes,1,opt,name=stderr,proto3" json:"stderr,omitempty"`
	Stdout []byte `protobuf:"bytes,2,opt,name=stdout,proto3" json:"stdout,omitempty"`
}

func (m *GetInstanceLogsReply) Reset()                    { *m = GetInstanceLogsReply{} }
func (m *GetInstanceLogsReply) String() string            { return proto.CompactTextString(m) }
func (*GetInstanceLogsReply) ProtoMessage()               {}
func (*GetInstanceLogsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type Instance struct {
	Job   *Job           `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
	State *InstanceState `protobuf:"bytes,2,opt,name=state" json:"state,omitempty"`
}

func (m *Instance) Reset()                    { *m = Instance{} }
func (m *Instance) String() string            { return proto.CompactTextString(m) }
func (*Instance) ProtoMessage()               {}
func (*Instance) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Instance) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

func (m *Instance) GetState() *InstanceState {
	if m != nil {
		return m.State
	}
	return nil
}

// The mutable state stored alongside the immutable job definition
type InstanceState struct {
	Status        TaskStatus `protobuf:"varint,1,opt,name=status,enum=jobbo.TaskStatus" json:"status,omitempty"`
	FailureReason string     `protobuf:"bytes,2,opt,name=failure_reason,json=failureReason" json:"failure_reason,omitempty"`
}

func (m *InstanceState) Reset()                    { *m = InstanceState{} }
func (m *InstanceState) String() string            { return proto.CompactTextString(m) }
func (*InstanceState) ProtoMessage()               {}
func (*InstanceState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

// Defines a job. Immutable once returned from RunJob
type Job struct {
	Metadata  *Metadata  `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	Resources *Resources `protobuf:"bytes,2,opt,name=resources" json:"resources,omitempty"`
	Command   *Command   `protobuf:"bytes,3,opt,name=command" json:"command,omitempty"`
	Container *Container `protobuf:"bytes,4,opt,name=container" json:"container,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Job) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Job) GetResources() *Resources {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *Job) GetCommand() *Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *Job) GetContainer() *Container {
	if m != nil {
		return m.Container
	}
	return nil
}

type Metadata struct {
	Name         string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Uid          string                     `protobuf:"bytes,2,opt,name=uid" json:"uid,omitempty"`
	CreationTime *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=creation_time,json=creationTime" json:"creation_time,omitempty"`
	Labels       map[string]string          `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Metadata) Reset()                    { *m = Metadata{} }
func (m *Metadata) String() string            { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()               {}
func (*Metadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Metadata) GetCreationTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.CreationTime
	}
	return nil
}

func (m *Metadata) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type Resources struct {
	Mem  float32 `protobuf:"fixed32,1,opt,name=mem" json:"mem,omitempty"`
	Cpus float32 `protobuf:"fixed32,2,opt,name=cpus" json:"cpus,omitempty"`
	Disk float32 `protobuf:"fixed32,3,opt,name=disk" json:"disk,omitempty"`
}

func (m *Resources) Reset()                    { *m = Resources{} }
func (m *Resources) String() string            { return proto.CompactTextString(m) }
func (*Resources) ProtoMessage()               {}
func (*Resources) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type Container struct {
	Image          string                   `protobuf:"bytes,1,opt,name=image" json:"image,omitempty"`
	PortMappings   []*Container_PortMapping `protobuf:"bytes,2,rep,name=port_mappings,json=portMappings" json:"port_mappings,omitempty"`
	ForcePullImage bool                     `protobuf:"varint,3,opt,name=force_pull_image,json=forcePullImage" json:"force_pull_image,omitempty"`
	Volumes        []*Container_Volume      `protobuf:"bytes,4,rep,name=volumes" json:"volumes,omitempty"`
}

func (m *Container) Reset()                    { *m = Container{} }
func (m *Container) String() string            { return proto.CompactTextString(m) }
func (*Container) ProtoMessage()               {}
func (*Container) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Container) GetPortMappings() []*Container_PortMapping {
	if m != nil {
		return m.PortMappings
	}
	return nil
}

func (m *Container) GetVolumes() []*Container_Volume {
	if m != nil {
		return m.Volumes
	}
	return nil
}

type Container_PortMapping struct {
	HostPort      uint32 `protobuf:"varint,1,opt,name=host_port,json=hostPort" json:"host_port,omitempty"`
	ContainerPort uint32 `protobuf:"varint,2,opt,name=container_port,json=containerPort" json:"container_port,omitempty"`
	// Protocol to expose as (ie: tcp, udp).
	Protocol string `protobuf:"bytes,3,opt,name=protocol" json:"protocol,omitempty"`
}

func (m *Container_PortMapping) Reset()                    { *m = Container_PortMapping{} }
func (m *Container_PortMapping) String() string            { return proto.CompactTextString(m) }
func (*Container_PortMapping) ProtoMessage()               {}
func (*Container_PortMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13, 0} }

type Container_Volume struct {
	ReadWrite bool `protobuf:"varint,1,opt,name=read_write,json=readWrite" json:"read_write,omitempty"`
	// Path pointing to a directory or file in the container. If the
	// path is a relative path, it is relative to the container work
	// directory. If the path is an absolute path, that path must
	// already exist.
	ContainerPath string `protobuf:"bytes,2,opt,name=container_path,json=containerPath" json:"container_path,omitempty"`
	// Absolute path pointing to a directory or file on the host or a
	// path relative to the container work directory.
	HostPath string `protobuf:"bytes,3,opt,name=host_path,json=hostPath" json:"host_path,omitempty"`
}

func (m *Container_Volume) Reset()                    { *m = Container_Volume{} }
func (m *Container_Volume) String() string            { return proto.CompactTextString(m) }
func (*Container_Volume) ProtoMessage()               {}
func (*Container_Volume) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13, 1} }

type Command struct {
	Uris        []*Command_URI    `protobuf:"bytes,1,rep,name=uris" json:"uris,omitempty"`
	Environment map[string]string `protobuf:"bytes,2,rep,name=environment" json:"environment,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Shell       bool              `protobuf:"varint,3,opt,name=shell" json:"shell,omitempty"`
	Value       string            `protobuf:"bytes,4,opt,name=value" json:"value,omitempty"`
	Arguments   []string          `protobuf:"bytes,5,rep,name=arguments" json:"arguments,omitempty"`
	User        string            `protobuf:"bytes,6,opt,name=user" json:"user,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *Command) GetUris() []*Command_URI {
	if m != nil {
		return m.Uris
	}
	return nil
}

func (m *Command) GetEnvironment() map[string]string {
	if m != nil {
		return m.Environment
	}
	return nil
}

type Command_URI struct {
	Value      string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	Executable bool   `protobuf:"varint,2,opt,name=executable" json:"executable,omitempty"`
	// In case the fetched file is recognized as an archive, extract
	// its contents into the sandbox. Note that a cached archive is
	// not copied from the cache to the sandbox in case extraction
	// originates from an archive in the cache.
	Extract bool `protobuf:"varint,3,opt,name=extract" json:"extract,omitempty"`
	// If this field is "true", the fetcher cache will be used. If not,
	// fetching bypasses the cache and downloads directly into the
	// sandbox directory, no matter whether a suitable cache file is
	// available or not. The former directs the fetcher to download to
	// the file cache, then copy from there to the sandbox. Subsequent
	// fetch attempts with the same URI will omit downloading and copy
	// from the cache as long as the file is resident there. Cache files
	// may get evicted at any time, which then leads to renewed
	// downloading. See also "docs/fetcher.md" and
	// "docs/fetcher-cache-internals.md".
	Cache bool `protobuf:"varint,4,opt,name=cache" json:"cache,omitempty"`
	// The fetcher's default behavior is to use the URI string's basename to
	// name the local copy. If this field is provided, the local copy will be
	// named with its value instead. If there is a directory component (which
	// must be a relative path), the local copy will be stored in that
	// subdirectory inside the sandbox.
	OutputFile string `protobuf:"bytes,5,opt,name=output_file,json=outputFile" json:"output_file,omitempty"`
}

func (m *Command_URI) Reset()                    { *m = Command_URI{} }
func (m *Command_URI) String() string            { return proto.CompactTextString(m) }
func (*Command_URI) ProtoMessage()               {}
func (*Command_URI) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14, 0} }

func init() {
	proto.RegisterType((*RunJobRequest)(nil), "jobbo.RunJobRequest")
	proto.RegisterType((*RunJobReply)(nil), "jobbo.RunJobReply")
	proto.RegisterType((*ListInstancesRequest)(nil), "jobbo.ListInstancesRequest")
	proto.RegisterType((*ListInstancesReply)(nil), "jobbo.ListInstancesReply")
	proto.RegisterType((*RemoveInstanceRequest)(nil), "jobbo.RemoveInstanceRequest")
	proto.RegisterType((*RemoveInstanceReply)(nil), "jobbo.RemoveInstanceReply")
	proto.RegisterType((*GetInstanceLogsRequest)(nil), "jobbo.GetInstanceLogsRequest")
	proto.RegisterType((*GetInstanceLogsReply)(nil), "jobbo.GetInstanceLogsReply")
	proto.RegisterType((*Instance)(nil), "jobbo.Instance")
	proto.RegisterType((*InstanceState)(nil), "jobbo.InstanceState")
	proto.RegisterType((*Job)(nil), "jobbo.Job")
	proto.RegisterType((*Metadata)(nil), "jobbo.Metadata")
	proto.RegisterType((*Resources)(nil), "jobbo.Resources")
	proto.RegisterType((*Container)(nil), "jobbo.Container")
	proto.RegisterType((*Container_PortMapping)(nil), "jobbo.Container.PortMapping")
	proto.RegisterType((*Container_Volume)(nil), "jobbo.Container.Volume")
	proto.RegisterType((*Command)(nil), "jobbo.Command")
	proto.RegisterType((*Command_URI)(nil), "jobbo.Command.URI")
	proto.RegisterEnum("jobbo.TaskStatus", TaskStatus_name, TaskStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Jobbo service

type JobboClient interface {
	// Creates a new job, and adds it to the list of pending jobs
	RunJob(ctx context.Context, in *RunJobRequest, opts ...grpc.CallOption) (*RunJobReply, error)
	// Lists instances under management
	ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesReply, error)
	// Removes an instance
	RemoveInstance(ctx context.Context, in *RemoveInstanceRequest, opts ...grpc.CallOption) (*RemoveInstanceReply, error)
	GetInstanceLogs(ctx context.Context, in *GetInstanceLogsRequest, opts ...grpc.CallOption) (*GetInstanceLogsReply, error)
}

type jobboClient struct {
	cc *grpc.ClientConn
}

func NewJobboClient(cc *grpc.ClientConn) JobboClient {
	return &jobboClient{cc}
}

func (c *jobboClient) RunJob(ctx context.Context, in *RunJobRequest, opts ...grpc.CallOption) (*RunJobReply, error) {
	out := new(RunJobReply)
	err := grpc.Invoke(ctx, "/jobbo.Jobbo/RunJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobboClient) ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesReply, error) {
	out := new(ListInstancesReply)
	err := grpc.Invoke(ctx, "/jobbo.Jobbo/ListInstances", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobboClient) RemoveInstance(ctx context.Context, in *RemoveInstanceRequest, opts ...grpc.CallOption) (*RemoveInstanceReply, error) {
	out := new(RemoveInstanceReply)
	err := grpc.Invoke(ctx, "/jobbo.Jobbo/RemoveInstance", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobboClient) GetInstanceLogs(ctx context.Context, in *GetInstanceLogsRequest, opts ...grpc.CallOption) (*GetInstanceLogsReply, error) {
	out := new(GetInstanceLogsReply)
	err := grpc.Invoke(ctx, "/jobbo.Jobbo/GetInstanceLogs", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Jobbo service

type JobboServer interface {
	// Creates a new job, and adds it to the list of pending jobs
	RunJob(context.Context, *RunJobRequest) (*RunJobReply, error)
	// Lists instances under management
	ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesReply, error)
	// Removes an instance
	RemoveInstance(context.Context, *RemoveInstanceRequest) (*RemoveInstanceReply, error)
	GetInstanceLogs(context.Context, *GetInstanceLogsRequest) (*GetInstanceLogsReply, error)
}

func RegisterJobboServer(s *grpc.Server, srv JobboServer) {
	s.RegisterService(&_Jobbo_serviceDesc, srv)
}

func _Jobbo_RunJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobboServer).RunJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobbo.Jobbo/RunJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobboServer).RunJob(ctx, req.(*RunJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobbo_ListInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobboServer).ListInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobbo.Jobbo/ListInstances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobboServer).ListInstances(ctx, req.(*ListInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobbo_RemoveInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobboServer).RemoveInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobbo.Jobbo/RemoveInstance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobboServer).RemoveInstance(ctx, req.(*RemoveInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobbo_GetInstanceLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInstanceLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobboServer).GetInstanceLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobbo.Jobbo/GetInstanceLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobboServer).GetInstanceLogs(ctx, req.(*GetInstanceLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Jobbo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "jobbo.Jobbo",
	HandlerType: (*JobboServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunJob",
			Handler:    _Jobbo_RunJob_Handler,
		},
		{
			MethodName: "ListInstances",
			Handler:    _Jobbo_ListInstances_Handler,
		},
		{
			MethodName: "RemoveInstance",
			Handler:    _Jobbo_RemoveInstance_Handler,
		},
		{
			MethodName: "GetInstanceLogs",
			Handler:    _Jobbo_GetInstanceLogs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("jobbo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1117 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x56, 0xdd, 0x72, 0xdb, 0x44,
	0x14, 0xae, 0xfc, 0x17, 0xfb, 0x38, 0x4e, 0xcd, 0x12, 0x8a, 0x51, 0xd2, 0x26, 0xa3, 0x19, 0x98,
	0xd0, 0x4e, 0xd5, 0xc1, 0xe5, 0x02, 0xb8, 0x28, 0x93, 0x1f, 0xa7, 0x75, 0x48, 0xdd, 0xcc, 0x26,
	0xa1, 0x37, 0xcc, 0x78, 0xd6, 0xf2, 0xc6, 0x56, 0x23, 0x69, 0x55, 0x69, 0x15, 0x9a, 0x77, 0xe0,
	0x15, 0xb8, 0xe0, 0x8a, 0x97, 0xe0, 0x65, 0x78, 0x13, 0x66, 0xff, 0x24, 0x5b, 0x49, 0xe8, 0xf4,
	0x6e, 0xcf, 0x77, 0xbe, 0xf3, 0xbb, 0xe7, 0xac, 0x04, 0xed, 0x77, 0x6c, 0x32, 0x61, 0x6e, 0x9c,
	0x30, 0xce, 0x50, 0x5d, 0x0a, 0xf6, 0xd6, 0x8c, 0xb1, 0x59, 0x40, 0x9f, 0x49, 0x70, 0x92, 0x5d,
	0x3c, 0xe3, 0x7e, 0x48, 0x53, 0x4e, 0xc2, 0x58, 0xf1, 0x9c, 0xa7, 0xd0, 0xc1, 0x59, 0x74, 0xc4,
	0x26, 0x98, 0xbe, 0xcf, 0x68, 0xca, 0xd1, 0x26, 0x54, 0xdf, 0xb1, 0x49, 0xcf, 0xda, 0xb6, 0x76,
	0xda, 0x7d, 0x70, 0x95, 0x4f, 0xa1, 0x17, 0xb0, 0xf3, 0x04, 0xda, 0x86, 0x1e, 0x07, 0xd7, 0x1f,
	0x21, 0xcf, 0x60, 0xfd, 0xd8, 0x4f, 0xf9, 0x30, 0x4a, 0x39, 0x89, 0x3c, 0x9a, 0x9a, 0x10, 0xdf,
	0x42, 0x23, 0xe5, 0x84, 0x67, 0xa9, 0x34, 0x5c, 0xeb, 0x7f, 0xa6, 0x0d, 0xcf, 0x48, 0x7a, 0x79,
	0x2a, 0x15, 0x58, 0x13, 0x10, 0x82, 0x5a, 0x44, 0x42, 0xda, 0xab, 0x6c, 0x5b, 0x3b, 0x2d, 0x2c,
	0xcf, 0xa8, 0x0b, 0xd5, 0xcc, 0x9f, 0xf6, 0xaa, 0x12, 0x12, 0x47, 0x67, 0x1f, 0x50, 0x29, 0x90,
	0x48, 0xee, 0x29, 0xb4, 0x7c, 0x83, 0xf4, 0xac, 0xed, 0xea, 0x4e, 0xbb, 0x7f, 0x5f, 0x47, 0x32,
	0x4c, 0x5c, 0x30, 0x9c, 0x33, 0xf8, 0x02, 0xd3, 0x90, 0x5d, 0xd1, 0x5c, 0xa9, 0xd3, 0xd5, 0xf1,
	0xac, 0x3c, 0xde, 0x42, 0x01, 0x95, 0x8f, 0x14, 0xe0, 0xec, 0xc1, 0xe7, 0x65, 0xaf, 0x22, 0xb7,
	0x27, 0xd0, 0x34, 0x91, 0x75, 0xf7, 0x6e, 0xa4, 0x96, 0x13, 0x9c, 0x17, 0xf0, 0xe0, 0x25, 0xcd,
	0xab, 0x3b, 0x66, 0xb3, 0xf4, 0xee, 0xd4, 0x10, 0xd4, 0x38, 0xf1, 0x03, 0x99, 0x58, 0x07, 0xcb,
	0xb3, 0x73, 0x08, 0xeb, 0x37, 0xec, 0x45, 0x12, 0x0f, 0x44, 0x19, 0x53, 0x9a, 0x24, 0xd2, 0xc1,
	0x2a, 0xd6, 0x92, 0xc6, 0x59, 0xc6, 0xa5, 0x17, 0x85, 0xb3, 0x8c, 0x3b, 0x67, 0xd0, 0x34, 0x4e,
	0xfe, 0xff, 0xe6, 0xd1, 0x63, 0xa8, 0x8b, 0xfa, 0xd5, 0xbd, 0xb5, 0xfb, 0xeb, 0xa5, 0xda, 0x44,
	0x8f, 0x28, 0x56, 0x14, 0x87, 0x40, 0x67, 0x09, 0xff, 0x94, 0xf1, 0xf8, 0x1a, 0xd6, 0x2e, 0x88,
	0x1f, 0x64, 0x09, 0x1d, 0x27, 0x94, 0xa4, 0x2c, 0xd2, 0x83, 0xd2, 0xd1, 0x28, 0x96, 0xa0, 0xf3,
	0x8f, 0x05, 0xd5, 0x23, 0x36, 0x11, 0x5d, 0x0f, 0x29, 0x27, 0x53, 0xc2, 0x49, 0xa9, 0xeb, 0xaf,
	0x35, 0x8c, 0x73, 0x02, 0x72, 0xa1, 0x95, 0xd0, 0x94, 0x65, 0x89, 0x18, 0x1f, 0x55, 0x47, 0x57,
	0xb3, 0xb1, 0xc1, 0x71, 0x41, 0x41, 0x3b, 0xb0, 0xe2, 0xb1, 0x30, 0x24, 0x91, 0x1a, 0xcd, 0x76,
	0x7f, 0x4d, 0xb3, 0xf7, 0x15, 0x8a, 0x8d, 0x5a, 0x78, 0xf6, 0x58, 0xc4, 0x89, 0x1f, 0xd1, 0xa4,
	0x57, 0x5b, 0xf2, 0xbc, 0x6f, 0x70, 0x5c, 0x50, 0x9c, 0x7f, 0x2d, 0x68, 0x9a, 0x04, 0xf3, 0x8d,
	0xb0, 0x6e, 0x6e, 0x44, 0xa5, 0x18, 0x83, 0x9f, 0xa1, 0xe3, 0x25, 0x94, 0x70, 0x9f, 0x45, 0x63,
	0xb1, 0xf2, 0x3a, 0x25, 0xdb, 0x55, 0xef, 0x81, 0x6b, 0xde, 0x03, 0xf7, 0xcc, 0xbc, 0x07, 0x78,
	0xd5, 0x18, 0x08, 0x08, 0x3d, 0x87, 0x46, 0x40, 0x26, 0x34, 0x48, 0x7b, 0x35, 0xb9, 0x39, 0x1b,
	0xa5, 0x46, 0xb9, 0xc7, 0x52, 0x3b, 0x88, 0x78, 0x72, 0x8d, 0x35, 0xd5, 0xfe, 0x11, 0xda, 0x0b,
	0xb0, 0x48, 0xeb, 0x92, 0x5e, 0x9b, 0xe9, 0xbc, 0xa4, 0xd7, 0x68, 0x1d, 0xea, 0x57, 0x24, 0xc8,
	0xcc, 0x3e, 0x2b, 0xe1, 0xa7, 0xca, 0x0f, 0x96, 0x33, 0x80, 0x56, 0xde, 0x55, 0x61, 0x18, 0xd2,
	0x50, 0x1a, 0x56, 0xb0, 0x38, 0x8a, 0xaa, 0xbd, 0x58, 0xef, 0x5b, 0x05, 0xcb, 0xb3, 0xc0, 0xa6,
	0x7e, 0x7a, 0x29, 0x4b, 0xab, 0x60, 0x79, 0x76, 0xfe, 0xac, 0x42, 0x2b, 0xef, 0xa1, 0x08, 0xe7,
	0x87, 0x64, 0x66, 0x9a, 0xa5, 0x04, 0xb4, 0x0b, 0x9d, 0x98, 0x25, 0x7c, 0x1c, 0x92, 0x38, 0xf6,
	0xa3, 0x99, 0x70, 0x2a, 0x2a, 0xdc, 0x2c, 0x5f, 0x81, 0x7b, 0xc2, 0x12, 0xfe, 0x5a, 0x91, 0xf0,
	0x6a, 0x5c, 0x08, 0xe2, 0xae, 0xbb, 0x17, 0x2c, 0xf1, 0xe8, 0x38, 0xce, 0x82, 0x60, 0xac, 0x62,
	0x88, 0x34, 0x9a, 0x78, 0x4d, 0xe2, 0x27, 0x59, 0x10, 0x0c, 0x65, 0xb0, 0xef, 0x60, 0xe5, 0x8a,
	0x05, 0x59, 0x48, 0x4d, 0x23, 0xbf, 0xbc, 0x11, 0xe6, 0x57, 0xa9, 0xc7, 0x86, 0x67, 0x87, 0xd0,
	0x5e, 0x88, 0x8c, 0x36, 0xa0, 0x35, 0x67, 0x29, 0x1f, 0x8b, 0x04, 0x64, 0x21, 0x1d, 0xdc, 0x14,
	0x80, 0xe0, 0x88, 0x05, 0xc8, 0xe7, 0x44, 0x31, 0xd4, 0xe2, 0x77, 0x72, 0x54, 0xd2, 0x6c, 0x68,
	0xca, 0x1b, 0xf7, 0x58, 0xa0, 0xdf, 0xcd, 0x5c, 0xb6, 0x2f, 0xa1, 0xa1, 0x32, 0x40, 0x0f, 0x01,
	0x12, 0x4a, 0xa6, 0xe3, 0xdf, 0x13, 0x9f, 0xab, 0x9e, 0x35, 0xc5, 0x80, 0x93, 0xe9, 0x5b, 0x01,
	0x94, 0x62, 0x11, 0x3e, 0x37, 0xcb, 0x56, 0xc4, 0x22, 0x7c, 0x5e, 0xe4, 0x2b, 0x18, 0x3a, 0x98,
	0xcc, 0x97, 0xf0, 0xb9, 0xf3, 0x57, 0x15, 0x56, 0xf4, 0x3e, 0xa0, 0x6f, 0xa0, 0x96, 0x25, 0xbe,
	0x79, 0x9a, 0xd1, 0xf2, 0xb6, 0xb8, 0xe7, 0x78, 0x88, 0xa5, 0x1e, 0xed, 0x42, 0x9b, 0x46, 0x57,
	0x7e, 0xc2, 0xa2, 0x90, 0x46, 0x5c, 0xdf, 0xd6, 0x56, 0x89, 0x3e, 0x28, 0x18, 0x6a, 0x26, 0x17,
	0x6d, 0xc4, 0x20, 0xa4, 0x73, 0x1a, 0x04, 0xfa, 0x92, 0x94, 0x50, 0x4c, 0x63, 0x6d, 0x61, 0x1a,
	0xd1, 0x26, 0xb4, 0x48, 0x32, 0xcb, 0x84, 0x5d, 0xda, 0xab, 0x6f, 0x57, 0x77, 0x5a, 0xb8, 0x00,
	0xc4, 0xd0, 0x65, 0x29, 0x4d, 0x7a, 0x0d, 0xb5, 0x7e, 0xe2, 0x6c, 0xff, 0x61, 0x41, 0xf5, 0x1c,
	0x0f, 0x0b, 0x7f, 0xd6, 0xa2, 0xbf, 0x47, 0x00, 0xf4, 0x03, 0xf5, 0x32, 0x4e, 0x26, 0x81, 0x1a,
	0xfc, 0x26, 0x5e, 0x40, 0x50, 0x0f, 0x56, 0xe8, 0x07, 0x9e, 0x10, 0x8f, 0xeb, 0xec, 0x8c, 0x28,
	0xfc, 0x79, 0xc4, 0x9b, 0xab, 0xfc, 0x9a, 0x58, 0x09, 0x68, 0x0b, 0xda, 0x2c, 0xe3, 0x71, 0xc6,
	0xc7, 0x17, 0x7e, 0x40, 0x7b, 0x75, 0x19, 0x0b, 0x14, 0x74, 0xe8, 0x07, 0xd4, 0x7e, 0x01, 0xdd,
	0x72, 0x37, 0x3e, 0x65, 0x15, 0x1f, 0xff, 0x06, 0x50, 0x3c, 0xb5, 0xa8, 0x0d, 0x2b, 0xe7, 0xa3,
	0x5f, 0x46, 0x6f, 0xde, 0x8e, 0xba, 0xf7, 0x84, 0x70, 0x32, 0x18, 0x1d, 0x0c, 0x47, 0x2f, 0xbb,
	0x16, 0xea, 0x40, 0xeb, 0x74, 0xff, 0xd5, 0xe0, 0xe0, 0xfc, 0x78, 0x70, 0xd0, 0xad, 0x08, 0x1d,
	0x3e, 0x1f, 0x8d, 0x84, 0xae, 0x8a, 0x00, 0x1a, 0x87, 0xbb, 0x43, 0xa1, 0xa8, 0xa1, 0x55, 0x68,
	0x1e, 0x0e, 0x47, 0xc3, 0xd3, 0x57, 0x83, 0x83, 0x6e, 0xbd, 0xff, 0x77, 0x05, 0xea, 0x47, 0xe2,
	0xea, 0xd0, 0xf7, 0xd0, 0x50, 0xff, 0x12, 0xc8, 0x7c, 0x1f, 0x96, 0xfe, 0x44, 0x6c, 0x54, 0x42,
	0xe3, 0xe0, 0xda, 0xb9, 0x87, 0x86, 0xd0, 0x59, 0xfa, 0xd6, 0x23, 0xf3, 0x32, 0xdd, 0xf6, 0xab,
	0x61, 0x7f, 0x75, 0xbb, 0x52, 0xb9, 0x3a, 0x86, 0xb5, 0xe5, 0x6f, 0x33, 0xda, 0xcc, 0x1f, 0xf8,
	0x5b, 0x7e, 0x04, 0x6c, 0xfb, 0x0e, 0xad, 0xf2, 0xf6, 0x06, 0xee, 0x97, 0xbe, 0xb2, 0xe8, 0xa1,
	0x36, 0xb8, 0xfd, 0xeb, 0x6d, 0x6f, 0xdc, 0xa5, 0x96, 0x0e, 0xf7, 0x5c, 0x78, 0x34, 0xf3, 0xf9,
	0x3c, 0x9b, 0xb8, 0x1e, 0x0b, 0xdd, 0xf7, 0xd9, 0xc4, 0xe7, 0x71, 0xc2, 0xa6, 0x99, 0xc7, 0x53,
	0x65, 0xb8, 0x07, 0xb2, 0x91, 0x27, 0x62, 0x93, 0x4f, 0xac, 0x49, 0x43, 0xae, 0xf4, 0xf3, 0xff,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x17, 0x5d, 0x49, 0x2a, 0x08, 0x0a, 0x00, 0x00,
}
