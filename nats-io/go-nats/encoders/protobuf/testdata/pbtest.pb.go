// Code generated by protoc-gen-go.
// source: pbtest.proto
// DO NOT EDIT!

/*
Package testdata is a generated protocol buffer package.

It is generated from these files:
	pbtest.proto

It has these top-level messages:
	Person
*/
package testdata

import proto "protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type Person struct {
	Name     string             `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Age      int32              `protobuf:"varint,2,opt,name=age" json:"age,omitempty"`
	Address  string             `protobuf:"bytes,3,opt,name=address" json:"address,omitempty"`
	Children map[string]*Person `protobuf:"bytes,10,rep,name=children" json:"children,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Person) Reset()         { *m = Person{} }
func (m *Person) String() string { return proto.CompactTextString(m) }
func (*Person) ProtoMessage()    {}

func (m *Person) GetChildren() map[string]*Person {
	if m != nil {
		return m.Children
	}
	return nil
}

func init() {
}
