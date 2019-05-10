// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lvmd/proto/lvmd.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type LogicalVolume struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SizeGb               uint64   `protobuf:"varint,2,opt,name=size_gb,json=sizeGb,proto3" json:"size_gb,omitempty"`
	DevMajor             uint32   `protobuf:"varint,3,opt,name=dev_major,json=devMajor,proto3" json:"dev_major,omitempty"`
	DevMinor             uint32   `protobuf:"varint,4,opt,name=dev_minor,json=devMinor,proto3" json:"dev_minor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogicalVolume) Reset()         { *m = LogicalVolume{} }
func (m *LogicalVolume) String() string { return proto.CompactTextString(m) }
func (*LogicalVolume) ProtoMessage()    {}
func (*LogicalVolume) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{1}
}

func (m *LogicalVolume) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogicalVolume.Unmarshal(m, b)
}
func (m *LogicalVolume) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogicalVolume.Marshal(b, m, deterministic)
}
func (m *LogicalVolume) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogicalVolume.Merge(m, src)
}
func (m *LogicalVolume) XXX_Size() int {
	return xxx_messageInfo_LogicalVolume.Size(m)
}
func (m *LogicalVolume) XXX_DiscardUnknown() {
	xxx_messageInfo_LogicalVolume.DiscardUnknown(m)
}

var xxx_messageInfo_LogicalVolume proto.InternalMessageInfo

func (m *LogicalVolume) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LogicalVolume) GetSizeGb() uint64 {
	if m != nil {
		return m.SizeGb
	}
	return 0
}

func (m *LogicalVolume) GetDevMajor() uint32 {
	if m != nil {
		return m.DevMajor
	}
	return 0
}

func (m *LogicalVolume) GetDevMinor() uint32 {
	if m != nil {
		return m.DevMinor
	}
	return 0
}

type CreateLVRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SizeGb               uint64   `protobuf:"varint,2,opt,name=size_gb,json=sizeGb,proto3" json:"size_gb,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateLVRequest) Reset()         { *m = CreateLVRequest{} }
func (m *CreateLVRequest) String() string { return proto.CompactTextString(m) }
func (*CreateLVRequest) ProtoMessage()    {}
func (*CreateLVRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{2}
}

func (m *CreateLVRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateLVRequest.Unmarshal(m, b)
}
func (m *CreateLVRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateLVRequest.Marshal(b, m, deterministic)
}
func (m *CreateLVRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateLVRequest.Merge(m, src)
}
func (m *CreateLVRequest) XXX_Size() int {
	return xxx_messageInfo_CreateLVRequest.Size(m)
}
func (m *CreateLVRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateLVRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateLVRequest proto.InternalMessageInfo

func (m *CreateLVRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateLVRequest) GetSizeGb() uint64 {
	if m != nil {
		return m.SizeGb
	}
	return 0
}

type CreateLVResponse struct {
	CommandOutput        string         `protobuf:"bytes,1,opt,name=command_output,json=commandOutput,proto3" json:"command_output,omitempty"`
	Volume               *LogicalVolume `protobuf:"bytes,2,opt,name=volume,proto3" json:"volume,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateLVResponse) Reset()         { *m = CreateLVResponse{} }
func (m *CreateLVResponse) String() string { return proto.CompactTextString(m) }
func (*CreateLVResponse) ProtoMessage()    {}
func (*CreateLVResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{3}
}

func (m *CreateLVResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateLVResponse.Unmarshal(m, b)
}
func (m *CreateLVResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateLVResponse.Marshal(b, m, deterministic)
}
func (m *CreateLVResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateLVResponse.Merge(m, src)
}
func (m *CreateLVResponse) XXX_Size() int {
	return xxx_messageInfo_CreateLVResponse.Size(m)
}
func (m *CreateLVResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateLVResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateLVResponse proto.InternalMessageInfo

func (m *CreateLVResponse) GetCommandOutput() string {
	if m != nil {
		return m.CommandOutput
	}
	return ""
}

func (m *CreateLVResponse) GetVolume() *LogicalVolume {
	if m != nil {
		return m.Volume
	}
	return nil
}

type RemoveLVRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveLVRequest) Reset()         { *m = RemoveLVRequest{} }
func (m *RemoveLVRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveLVRequest) ProtoMessage()    {}
func (*RemoveLVRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{4}
}

func (m *RemoveLVRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveLVRequest.Unmarshal(m, b)
}
func (m *RemoveLVRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveLVRequest.Marshal(b, m, deterministic)
}
func (m *RemoveLVRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveLVRequest.Merge(m, src)
}
func (m *RemoveLVRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveLVRequest.Size(m)
}
func (m *RemoveLVRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveLVRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveLVRequest proto.InternalMessageInfo

func (m *RemoveLVRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RemoveLVResponse struct {
	CommandOutput        string   `protobuf:"bytes,1,opt,name=command_output,json=commandOutput,proto3" json:"command_output,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveLVResponse) Reset()         { *m = RemoveLVResponse{} }
func (m *RemoveLVResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveLVResponse) ProtoMessage()    {}
func (*RemoveLVResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{5}
}

func (m *RemoveLVResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveLVResponse.Unmarshal(m, b)
}
func (m *RemoveLVResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveLVResponse.Marshal(b, m, deterministic)
}
func (m *RemoveLVResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveLVResponse.Merge(m, src)
}
func (m *RemoveLVResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveLVResponse.Size(m)
}
func (m *RemoveLVResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveLVResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveLVResponse proto.InternalMessageInfo

func (m *RemoveLVResponse) GetCommandOutput() string {
	if m != nil {
		return m.CommandOutput
	}
	return ""
}

type ResizeLVRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	SizeGb               uint64   `protobuf:"varint,2,opt,name=size_gb,json=sizeGb,proto3" json:"size_gb,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResizeLVRequest) Reset()         { *m = ResizeLVRequest{} }
func (m *ResizeLVRequest) String() string { return proto.CompactTextString(m) }
func (*ResizeLVRequest) ProtoMessage()    {}
func (*ResizeLVRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{6}
}

func (m *ResizeLVRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResizeLVRequest.Unmarshal(m, b)
}
func (m *ResizeLVRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResizeLVRequest.Marshal(b, m, deterministic)
}
func (m *ResizeLVRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResizeLVRequest.Merge(m, src)
}
func (m *ResizeLVRequest) XXX_Size() int {
	return xxx_messageInfo_ResizeLVRequest.Size(m)
}
func (m *ResizeLVRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResizeLVRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResizeLVRequest proto.InternalMessageInfo

func (m *ResizeLVRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ResizeLVRequest) GetSizeGb() uint64 {
	if m != nil {
		return m.SizeGb
	}
	return 0
}

type ResizeLVResponse struct {
	CommandOutput        string   `protobuf:"bytes,1,opt,name=command_output,json=commandOutput,proto3" json:"command_output,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResizeLVResponse) Reset()         { *m = ResizeLVResponse{} }
func (m *ResizeLVResponse) String() string { return proto.CompactTextString(m) }
func (*ResizeLVResponse) ProtoMessage()    {}
func (*ResizeLVResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{7}
}

func (m *ResizeLVResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResizeLVResponse.Unmarshal(m, b)
}
func (m *ResizeLVResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResizeLVResponse.Marshal(b, m, deterministic)
}
func (m *ResizeLVResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResizeLVResponse.Merge(m, src)
}
func (m *ResizeLVResponse) XXX_Size() int {
	return xxx_messageInfo_ResizeLVResponse.Size(m)
}
func (m *ResizeLVResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ResizeLVResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ResizeLVResponse proto.InternalMessageInfo

func (m *ResizeLVResponse) GetCommandOutput() string {
	if m != nil {
		return m.CommandOutput
	}
	return ""
}

type GetLVListResponse struct {
	Volumes              []*LogicalVolume `protobuf:"bytes,1,rep,name=volumes,proto3" json:"volumes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *GetLVListResponse) Reset()         { *m = GetLVListResponse{} }
func (m *GetLVListResponse) String() string { return proto.CompactTextString(m) }
func (*GetLVListResponse) ProtoMessage()    {}
func (*GetLVListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{8}
}

func (m *GetLVListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLVListResponse.Unmarshal(m, b)
}
func (m *GetLVListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLVListResponse.Marshal(b, m, deterministic)
}
func (m *GetLVListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLVListResponse.Merge(m, src)
}
func (m *GetLVListResponse) XXX_Size() int {
	return xxx_messageInfo_GetLVListResponse.Size(m)
}
func (m *GetLVListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLVListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetLVListResponse proto.InternalMessageInfo

func (m *GetLVListResponse) GetVolumes() []*LogicalVolume {
	if m != nil {
		return m.Volumes
	}
	return nil
}

type GetFreeBytesResponse struct {
	FreeBytes            uint64   `protobuf:"varint,1,opt,name=free_bytes,json=freeBytes,proto3" json:"free_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFreeBytesResponse) Reset()         { *m = GetFreeBytesResponse{} }
func (m *GetFreeBytesResponse) String() string { return proto.CompactTextString(m) }
func (*GetFreeBytesResponse) ProtoMessage()    {}
func (*GetFreeBytesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb3c510e545f3bbd, []int{9}
}

func (m *GetFreeBytesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFreeBytesResponse.Unmarshal(m, b)
}
func (m *GetFreeBytesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFreeBytesResponse.Marshal(b, m, deterministic)
}
func (m *GetFreeBytesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFreeBytesResponse.Merge(m, src)
}
func (m *GetFreeBytesResponse) XXX_Size() int {
	return xxx_messageInfo_GetFreeBytesResponse.Size(m)
}
func (m *GetFreeBytesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFreeBytesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetFreeBytesResponse proto.InternalMessageInfo

func (m *GetFreeBytesResponse) GetFreeBytes() uint64 {
	if m != nil {
		return m.FreeBytes
	}
	return 0
}

func init() {
	proto.RegisterType((*Empty)(nil), "proto.Empty")
	proto.RegisterType((*LogicalVolume)(nil), "proto.LogicalVolume")
	proto.RegisterType((*CreateLVRequest)(nil), "proto.CreateLVRequest")
	proto.RegisterType((*CreateLVResponse)(nil), "proto.CreateLVResponse")
	proto.RegisterType((*RemoveLVRequest)(nil), "proto.RemoveLVRequest")
	proto.RegisterType((*RemoveLVResponse)(nil), "proto.RemoveLVResponse")
	proto.RegisterType((*ResizeLVRequest)(nil), "proto.ResizeLVRequest")
	proto.RegisterType((*ResizeLVResponse)(nil), "proto.ResizeLVResponse")
	proto.RegisterType((*GetLVListResponse)(nil), "proto.GetLVListResponse")
	proto.RegisterType((*GetFreeBytesResponse)(nil), "proto.GetFreeBytesResponse")
}

func init() { proto.RegisterFile("lvmd/proto/lvmd.proto", fileDescriptor_cb3c510e545f3bbd) }

var fileDescriptor_cb3c510e545f3bbd = []byte{
	// 417 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcf, 0xab, 0xda, 0x40,
	0x10, 0xc7, 0x49, 0x5f, 0x9e, 0x9a, 0xe9, 0xb3, 0xef, 0x75, 0xb1, 0x35, 0x28, 0x85, 0x10, 0x10,
	0x3c, 0x14, 0x05, 0xa5, 0x07, 0x11, 0x7a, 0xa8, 0xb4, 0xb9, 0xa4, 0x14, 0xb6, 0x90, 0x6b, 0x48,
	0xcc, 0x28, 0x29, 0x6e, 0x36, 0x4d, 0x36, 0x01, 0xdb, 0xff, 0xae, 0x7f, 0x59, 0xc9, 0x9a, 0x1f,
	0x6a, 0x68, 0xc1, 0x77, 0x72, 0x67, 0x76, 0xbe, 0xf3, 0xdd, 0xf9, 0x8c, 0x81, 0x37, 0x87, 0x9c,
	0x05, 0xf3, 0x38, 0xe1, 0x82, 0xcf, 0x8b, 0xe3, 0x4c, 0x1e, 0xc9, 0xbd, 0xfc, 0x31, 0xbb, 0x70,
	0xff, 0x99, 0xc5, 0xe2, 0x68, 0xe6, 0xd0, 0xb7, 0xf9, 0x3e, 0xdc, 0x7a, 0x07, 0x87, 0x1f, 0x32,
	0x86, 0x84, 0x80, 0x1a, 0x79, 0x0c, 0x75, 0xc5, 0x50, 0xa6, 0x1a, 0x95, 0x67, 0x32, 0x84, 0x6e,
	0x1a, 0xfe, 0x42, 0x77, 0xef, 0xeb, 0x2f, 0x0c, 0x65, 0xaa, 0xd2, 0x4e, 0x11, 0x5a, 0x3e, 0x19,
	0x83, 0x16, 0x60, 0xee, 0x32, 0xef, 0x07, 0x4f, 0xf4, 0x3b, 0x43, 0x99, 0xf6, 0x69, 0x2f, 0xc0,
	0xfc, 0x6b, 0x11, 0xd7, 0x97, 0x61, 0xc4, 0x13, 0x5d, 0x6d, 0x2e, 0x8b, 0xd8, 0xfc, 0x08, 0x8f,
	0x9b, 0x04, 0x3d, 0x81, 0xb6, 0x43, 0xf1, 0x67, 0x86, 0xa9, 0xb8, 0xc9, 0xd9, 0xdc, 0xc3, 0x53,
	0xa3, 0x4f, 0x63, 0x1e, 0xa5, 0x48, 0x26, 0xf0, 0x6a, 0xcb, 0x19, 0xf3, 0xa2, 0xc0, 0xe5, 0x99,
	0x88, 0x33, 0x51, 0xb6, 0xea, 0x97, 0xd9, 0x6f, 0x32, 0x49, 0xde, 0x43, 0x27, 0x97, 0xb3, 0xca,
	0x96, 0x2f, 0x17, 0x83, 0x13, 0x9a, 0xd9, 0x05, 0x07, 0x5a, 0xd6, 0x98, 0x13, 0x78, 0xa4, 0xc8,
	0x78, 0xfe, 0xff, 0x87, 0x9a, 0x2b, 0x78, 0x6a, 0xca, 0x6e, 0x7a, 0x4f, 0x81, 0x82, 0x62, 0x31,
	0xd6, 0x33, 0x51, 0x48, 0xeb, 0x4a, 0x7f, 0x9b, 0xf5, 0x06, 0x5e, 0x5b, 0x28, 0x6c, 0xc7, 0x0e,
	0x53, 0x51, 0x6b, 0x67, 0xd0, 0x3d, 0xcd, 0x9e, 0xea, 0x8a, 0x71, 0xf7, 0x4f, 0x40, 0x55, 0x91,
	0xf9, 0x01, 0x06, 0x16, 0x8a, 0x2f, 0x09, 0xe2, 0xa7, 0xa3, 0xc0, 0xb4, 0xee, 0xf3, 0x0e, 0x60,
	0x97, 0x20, 0xba, 0x7e, 0x91, 0x95, 0xfe, 0x2a, 0xd5, 0x76, 0x55, 0xd9, 0xe2, 0x8f, 0x02, 0x9a,
	0xed, 0x7c, 0xc7, 0x24, 0x0f, 0xb7, 0x48, 0xd6, 0xd0, 0xab, 0xf6, 0x49, 0xde, 0x96, 0x7e, 0x57,
	0x7f, 0x90, 0xd1, 0xb0, 0x95, 0x2f, 0x9d, 0xd6, 0xd0, 0xab, 0xe0, 0xd7, 0xe2, 0xab, 0xa5, 0xd5,
	0xe2, 0xd6, 0x96, 0xa4, 0xf8, 0x84, 0xef, 0x4c, 0x7c, 0xb1, 0x8f, 0x33, 0xf1, 0x25, 0xe7, 0xc5,
	0x6f, 0xd0, 0x1c, 0xab, 0x9a, 0x61, 0x09, 0x5a, 0x4d, 0x93, 0x3c, 0x94, 0x12, 0xf9, 0x99, 0x8d,
	0xf4, 0x32, 0x6a, 0xd3, 0x5e, 0xc1, 0xc3, 0x39, 0xbd, 0x2b, 0xdd, 0xb8, 0xd1, 0xb5, 0x00, 0xfb,
	0x1d, 0x79, 0xb7, 0xfc, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x08, 0xc5, 0xce, 0xeb, 0x03, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LVServiceClient is the client API for LVService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LVServiceClient interface {
	CreateLV(ctx context.Context, in *CreateLVRequest, opts ...grpc.CallOption) (*CreateLVResponse, error)
	RemoveLV(ctx context.Context, in *RemoveLVRequest, opts ...grpc.CallOption) (*RemoveLVResponse, error)
	ResizeLV(ctx context.Context, in *ResizeLVRequest, opts ...grpc.CallOption) (*ResizeLVResponse, error)
}

type lVServiceClient struct {
	cc *grpc.ClientConn
}

func NewLVServiceClient(cc *grpc.ClientConn) LVServiceClient {
	return &lVServiceClient{cc}
}

func (c *lVServiceClient) CreateLV(ctx context.Context, in *CreateLVRequest, opts ...grpc.CallOption) (*CreateLVResponse, error) {
	out := new(CreateLVResponse)
	err := c.cc.Invoke(ctx, "/proto.LVService/CreateLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lVServiceClient) RemoveLV(ctx context.Context, in *RemoveLVRequest, opts ...grpc.CallOption) (*RemoveLVResponse, error) {
	out := new(RemoveLVResponse)
	err := c.cc.Invoke(ctx, "/proto.LVService/RemoveLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lVServiceClient) ResizeLV(ctx context.Context, in *ResizeLVRequest, opts ...grpc.CallOption) (*ResizeLVResponse, error) {
	out := new(ResizeLVResponse)
	err := c.cc.Invoke(ctx, "/proto.LVService/ResizeLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LVServiceServer is the server API for LVService service.
type LVServiceServer interface {
	CreateLV(context.Context, *CreateLVRequest) (*CreateLVResponse, error)
	RemoveLV(context.Context, *RemoveLVRequest) (*RemoveLVResponse, error)
	ResizeLV(context.Context, *ResizeLVRequest) (*ResizeLVResponse, error)
}

func RegisterLVServiceServer(s *grpc.Server, srv LVServiceServer) {
	s.RegisterService(&_LVService_serviceDesc, srv)
}

func _LVService_CreateLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).CreateLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/CreateLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).CreateLV(ctx, req.(*CreateLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LVService_RemoveLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).RemoveLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/RemoveLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).RemoveLV(ctx, req.(*RemoveLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LVService_ResizeLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResizeLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).ResizeLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/ResizeLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).ResizeLV(ctx, req.(*ResizeLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LVService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LVService",
	HandlerType: (*LVServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLV",
			Handler:    _LVService_CreateLV_Handler,
		},
		{
			MethodName: "RemoveLV",
			Handler:    _LVService_RemoveLV_Handler,
		},
		{
			MethodName: "ResizeLV",
			Handler:    _LVService_ResizeLV_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lvmd/proto/lvmd.proto",
}

// VGServiceClient is the client API for VGService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VGServiceClient interface {
	GetLVList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetLVListResponse, error)
	GetFreeBytes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetFreeBytesResponse, error)
}

type vGServiceClient struct {
	cc *grpc.ClientConn
}

func NewVGServiceClient(cc *grpc.ClientConn) VGServiceClient {
	return &vGServiceClient{cc}
}

func (c *vGServiceClient) GetLVList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetLVListResponse, error) {
	out := new(GetLVListResponse)
	err := c.cc.Invoke(ctx, "/proto.VGService/GetLVList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vGServiceClient) GetFreeBytes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetFreeBytesResponse, error) {
	out := new(GetFreeBytesResponse)
	err := c.cc.Invoke(ctx, "/proto.VGService/GetFreeBytes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VGServiceServer is the server API for VGService service.
type VGServiceServer interface {
	GetLVList(context.Context, *Empty) (*GetLVListResponse, error)
	GetFreeBytes(context.Context, *Empty) (*GetFreeBytesResponse, error)
}

func RegisterVGServiceServer(s *grpc.Server, srv VGServiceServer) {
	s.RegisterService(&_VGService_serviceDesc, srv)
}

func _VGService_GetLVList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VGServiceServer).GetLVList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.VGService/GetLVList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VGServiceServer).GetLVList(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _VGService_GetFreeBytes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VGServiceServer).GetFreeBytes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.VGService/GetFreeBytes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VGServiceServer).GetFreeBytes(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _VGService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.VGService",
	HandlerType: (*VGServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLVList",
			Handler:    _VGService_GetLVList_Handler,
		},
		{
			MethodName: "GetFreeBytes",
			Handler:    _VGService_GetFreeBytes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lvmd/proto/lvmd.proto",
}
