// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        v3.17.3
// source: mail.proto

package mailbox

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Mail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageID  string     `protobuf:"bytes,1,opt,name=MessageID,proto3" json:"MessageID,omitempty"` // Unique ID number for this person.
	Subject    string     `protobuf:"bytes,2,opt,name=Subject,proto3" json:"Subject,omitempty"`
	To         []*Address `protobuf:"bytes,3,rep,name=to,proto3" json:"to,omitempty"`
	From       []*Address `protobuf:"bytes,4,rep,name=from,proto3" json:"from,omitempty"`
	Text       []*Body    `protobuf:"bytes,5,rep,name=text,proto3" json:"text,omitempty"`
	Attachment *Body      `protobuf:"bytes,6,opt,name=attachment,proto3" json:"attachment,omitempty"`
}

func (x *Mail) Reset() {
	*x = Mail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mail) ProtoMessage() {}

func (x *Mail) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mail.ProtoReflect.Descriptor instead.
func (*Mail) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{0}
}

func (x *Mail) GetMessageID() string {
	if x != nil {
		return x.MessageID
	}
	return ""
}

func (x *Mail) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Mail) GetTo() []*Address {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *Mail) GetFrom() []*Address {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Mail) GetText() []*Body {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *Mail) GetAttachment() *Body {
	if x != nil {
		return x.Attachment
	}
	return nil
}

type Body struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContentType string `protobuf:"bytes,1,opt,name=contentType,proto3" json:"contentType,omitempty"`
	MainBody    []byte `protobuf:"bytes,2,opt,name=mainBody,proto3" json:"mainBody,omitempty"`
}

func (x *Body) Reset() {
	*x = Body{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{1}
}

func (x *Body) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

func (x *Body) GetMainBody() []byte {
	if x != nil {
		return x.MainBody
	}
	return nil
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{2}
}

func (x *Address) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Address) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type SendMailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgID string `protobuf:"bytes,1,opt,name=MsgID,proto3" json:"MsgID,omitempty"` // Unique ID number for this person.
}

func (x *SendMailResponse) Reset() {
	*x = SendMailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMailResponse) ProtoMessage() {}

func (x *SendMailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMailResponse.ProtoReflect.Descriptor instead.
func (*SendMailResponse) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{3}
}

func (x *SendMailResponse) GetMsgID() string {
	if x != nil {
		return x.MsgID
	}
	return ""
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{4}
}

func (x *Server) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type UploadFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to File:
	//	*UploadFile_Name
	//	*UploadFile_ContentType
	//	*UploadFile_Content
	File isUploadFile_File `protobuf_oneof:"File"`
}

func (x *UploadFile) Reset() {
	*x = UploadFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFile) ProtoMessage() {}

func (x *UploadFile) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFile.ProtoReflect.Descriptor instead.
func (*UploadFile) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{5}
}

func (m *UploadFile) GetFile() isUploadFile_File {
	if m != nil {
		return m.File
	}
	return nil
}

func (x *UploadFile) GetName() string {
	if x, ok := x.GetFile().(*UploadFile_Name); ok {
		return x.Name
	}
	return ""
}

func (x *UploadFile) GetContentType() string {
	if x, ok := x.GetFile().(*UploadFile_ContentType); ok {
		return x.ContentType
	}
	return ""
}

func (x *UploadFile) GetContent() []byte {
	if x, ok := x.GetFile().(*UploadFile_Content); ok {
		return x.Content
	}
	return nil
}

type isUploadFile_File interface {
	isUploadFile_File()
}

type UploadFile_Name struct {
	Name string `protobuf:"bytes,1,opt,name=Name,proto3,oneof"`
}

type UploadFile_ContentType struct {
	ContentType string `protobuf:"bytes,2,opt,name=ContentType,proto3,oneof"`
}

type UploadFile_Content struct {
	Content []byte `protobuf:"bytes,3,opt,name=Content,proto3,oneof"`
}

func (*UploadFile_Name) isUploadFile_File() {}

func (*UploadFile_ContentType) isUploadFile_File() {}

func (*UploadFile_Content) isUploadFile_File() {}

type UploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileID string `protobuf:"bytes,1,opt,name=FileID,proto3" json:"FileID,omitempty"`
}

func (x *UploadResponse) Reset() {
	*x = UploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResponse) ProtoMessage() {}

func (x *UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResponse.ProtoReflect.Descriptor instead.
func (*UploadResponse) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{6}
}

func (x *UploadResponse) GetFileID() string {
	if x != nil {
		return x.FileID
	}
	return ""
}

type ServerList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Server `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
}

func (x *ServerList) Reset() {
	*x = ServerList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mail_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerList) ProtoMessage() {}

func (x *ServerList) ProtoReflect() protoreflect.Message {
	mi := &file_mail_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerList.ProtoReflect.Descriptor instead.
func (*ServerList) Descriptor() ([]byte, []int) {
	return file_mail_proto_rawDescGZIP(), []int{7}
}

func (x *ServerList) GetItems() []*Server {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_mail_proto protoreflect.FileDescriptor

var file_mail_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x61,
	0x69, 0x6c, 0x62, 0x6f, 0x78, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xd8, 0x01, 0x0a, 0x04, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x20, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x24, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x21, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x62, 0x6f, 0x78, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x2d,
	0x0a, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x42, 0x6f, 0x64,
	0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x44, 0x0a,
	0x04, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x69, 0x6e, 0x42,
	0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x61, 0x69, 0x6e, 0x42,
	0x6f, 0x64, 0x79, 0x22, 0x37, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x28, 0x0a, 0x10,
	0x53, 0x65, 0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x4d, 0x73, 0x67, 0x49, 0x44, 0x22, 0x1c, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x6a, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52,
	0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65,
	0x22, 0x28, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x22, 0x33, 0x0a, 0x0a, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x05, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f,
	0x78, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x32,
	0xe3, 0x01, 0x0a, 0x07, 0x4d, 0x61, 0x69, 0x6c, 0x42, 0x6f, 0x78, 0x12, 0x32, 0x0a, 0x04, 0x53,
	0x65, 0x6e, 0x64, 0x12, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x4d, 0x61,
	0x69, 0x6c, 0x1a, 0x19, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x3b, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x06,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x1a, 0x17, 0x2e, 0x6d, 0x61,
	0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x2b, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x0f, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x1a, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x62, 0x6f, 0x78, 0x2e, 0x4d, 0x61, 0x69,
	0x6c, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x62,
	0x6f, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mail_proto_rawDescOnce sync.Once
	file_mail_proto_rawDescData = file_mail_proto_rawDesc
)

func file_mail_proto_rawDescGZIP() []byte {
	file_mail_proto_rawDescOnce.Do(func() {
		file_mail_proto_rawDescData = protoimpl.X.CompressGZIP(file_mail_proto_rawDescData)
	})
	return file_mail_proto_rawDescData
}

var file_mail_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_mail_proto_goTypes = []interface{}{
	(*Mail)(nil),             // 0: mailbox.Mail
	(*Body)(nil),             // 1: mailbox.Body
	(*Address)(nil),          // 2: mailbox.Address
	(*SendMailResponse)(nil), // 3: mailbox.SendMailResponse
	(*Server)(nil),           // 4: mailbox.Server
	(*UploadFile)(nil),       // 5: mailbox.UploadFile
	(*UploadResponse)(nil),   // 6: mailbox.UploadResponse
	(*ServerList)(nil),       // 7: mailbox.ServerList
	(*empty.Empty)(nil),      // 8: google.protobuf.Empty
}
var file_mail_proto_depIdxs = []int32{
	2, // 0: mailbox.Mail.to:type_name -> mailbox.Address
	2, // 1: mailbox.Mail.from:type_name -> mailbox.Address
	1, // 2: mailbox.Mail.text:type_name -> mailbox.Body
	1, // 3: mailbox.Mail.attachment:type_name -> mailbox.Body
	4, // 4: mailbox.ServerList.Items:type_name -> mailbox.Server
	0, // 5: mailbox.MailBox.Send:input_type -> mailbox.Mail
	8, // 6: mailbox.MailBox.ListServer:input_type -> google.protobuf.Empty
	5, // 7: mailbox.MailBox.Upload:input_type -> mailbox.UploadFile
	4, // 8: mailbox.MailBox.Watch:input_type -> mailbox.Server
	3, // 9: mailbox.MailBox.Send:output_type -> mailbox.SendMailResponse
	7, // 10: mailbox.MailBox.ListServer:output_type -> mailbox.ServerList
	6, // 11: mailbox.MailBox.Upload:output_type -> mailbox.UploadResponse
	0, // 12: mailbox.MailBox.Watch:output_type -> mailbox.Mail
	9, // [9:13] is the sub-list for method output_type
	5, // [5:9] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_mail_proto_init() }
func file_mail_proto_init() {
	if File_mail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Body); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMailResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadFile); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mail_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_mail_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*UploadFile_Name)(nil),
		(*UploadFile_ContentType)(nil),
		(*UploadFile_Content)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mail_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mail_proto_goTypes,
		DependencyIndexes: file_mail_proto_depIdxs,
		MessageInfos:      file_mail_proto_msgTypes,
	}.Build()
	File_mail_proto = out.File
	file_mail_proto_rawDesc = nil
	file_mail_proto_goTypes = nil
	file_mail_proto_depIdxs = nil
}
