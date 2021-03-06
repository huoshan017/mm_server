// Code generated by protoc-gen-go. DO NOT EDIT.
// source: db_rpcsvr.proto

package db

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type PlayerStageTotalScoreHistoryTopData struct {
	Rank                 *int32   `protobuf:"varint,1,opt,name=Rank" json:"Rank,omitempty"`
	Score                *int32   `protobuf:"varint,2,opt,name=Score" json:"Score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerStageTotalScoreHistoryTopData) Reset()         { *m = PlayerStageTotalScoreHistoryTopData{} }
func (m *PlayerStageTotalScoreHistoryTopData) String() string { return proto.CompactTextString(m) }
func (*PlayerStageTotalScoreHistoryTopData) ProtoMessage()    {}
func (*PlayerStageTotalScoreHistoryTopData) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{0}
}

func (m *PlayerStageTotalScoreHistoryTopData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerStageTotalScoreHistoryTopData.Unmarshal(m, b)
}
func (m *PlayerStageTotalScoreHistoryTopData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerStageTotalScoreHistoryTopData.Marshal(b, m, deterministic)
}
func (m *PlayerStageTotalScoreHistoryTopData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerStageTotalScoreHistoryTopData.Merge(m, src)
}
func (m *PlayerStageTotalScoreHistoryTopData) XXX_Size() int {
	return xxx_messageInfo_PlayerStageTotalScoreHistoryTopData.Size(m)
}
func (m *PlayerStageTotalScoreHistoryTopData) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerStageTotalScoreHistoryTopData.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerStageTotalScoreHistoryTopData proto.InternalMessageInfo

func (m *PlayerStageTotalScoreHistoryTopData) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *PlayerStageTotalScoreHistoryTopData) GetScore() int32 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

type PlayerStageTotalScoreStage struct {
	Id                   *int32   `protobuf:"varint,1,opt,name=Id" json:"Id,omitempty"`
	TopScore             *int32   `protobuf:"varint,2,opt,name=TopScore" json:"TopScore,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerStageTotalScoreStage) Reset()         { *m = PlayerStageTotalScoreStage{} }
func (m *PlayerStageTotalScoreStage) String() string { return proto.CompactTextString(m) }
func (*PlayerStageTotalScoreStage) ProtoMessage()    {}
func (*PlayerStageTotalScoreStage) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{1}
}

func (m *PlayerStageTotalScoreStage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerStageTotalScoreStage.Unmarshal(m, b)
}
func (m *PlayerStageTotalScoreStage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerStageTotalScoreStage.Marshal(b, m, deterministic)
}
func (m *PlayerStageTotalScoreStage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerStageTotalScoreStage.Merge(m, src)
}
func (m *PlayerStageTotalScoreStage) XXX_Size() int {
	return xxx_messageInfo_PlayerStageTotalScoreStage.Size(m)
}
func (m *PlayerStageTotalScoreStage) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerStageTotalScoreStage.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerStageTotalScoreStage proto.InternalMessageInfo

func (m *PlayerStageTotalScoreStage) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *PlayerStageTotalScoreStage) GetTopScore() int32 {
	if m != nil && m.TopScore != nil {
		return *m.TopScore
	}
	return 0
}

type PlayerStageTotalScoreStageList struct {
	List                 []*PlayerStageTotalScoreStage `protobuf:"bytes,1,rep,name=List" json:"List,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *PlayerStageTotalScoreStageList) Reset()         { *m = PlayerStageTotalScoreStageList{} }
func (m *PlayerStageTotalScoreStageList) String() string { return proto.CompactTextString(m) }
func (*PlayerStageTotalScoreStageList) ProtoMessage()    {}
func (*PlayerStageTotalScoreStageList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{2}
}

func (m *PlayerStageTotalScoreStageList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerStageTotalScoreStageList.Unmarshal(m, b)
}
func (m *PlayerStageTotalScoreStageList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerStageTotalScoreStageList.Marshal(b, m, deterministic)
}
func (m *PlayerStageTotalScoreStageList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerStageTotalScoreStageList.Merge(m, src)
}
func (m *PlayerStageTotalScoreStageList) XXX_Size() int {
	return xxx_messageInfo_PlayerStageTotalScoreStageList.Size(m)
}
func (m *PlayerStageTotalScoreStageList) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerStageTotalScoreStageList.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerStageTotalScoreStageList proto.InternalMessageInfo

func (m *PlayerStageTotalScoreStageList) GetList() []*PlayerStageTotalScoreStage {
	if m != nil {
		return m.List
	}
	return nil
}

type PlayerCharmHistoryTopData struct {
	Rank                 *int32   `protobuf:"varint,1,opt,name=Rank" json:"Rank,omitempty"`
	Charm                *int32   `protobuf:"varint,2,opt,name=Charm" json:"Charm,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerCharmHistoryTopData) Reset()         { *m = PlayerCharmHistoryTopData{} }
func (m *PlayerCharmHistoryTopData) String() string { return proto.CompactTextString(m) }
func (*PlayerCharmHistoryTopData) ProtoMessage()    {}
func (*PlayerCharmHistoryTopData) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{3}
}

func (m *PlayerCharmHistoryTopData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerCharmHistoryTopData.Unmarshal(m, b)
}
func (m *PlayerCharmHistoryTopData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerCharmHistoryTopData.Marshal(b, m, deterministic)
}
func (m *PlayerCharmHistoryTopData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerCharmHistoryTopData.Merge(m, src)
}
func (m *PlayerCharmHistoryTopData) XXX_Size() int {
	return xxx_messageInfo_PlayerCharmHistoryTopData.Size(m)
}
func (m *PlayerCharmHistoryTopData) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerCharmHistoryTopData.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerCharmHistoryTopData proto.InternalMessageInfo

func (m *PlayerCharmHistoryTopData) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *PlayerCharmHistoryTopData) GetCharm() int32 {
	if m != nil && m.Charm != nil {
		return *m.Charm
	}
	return 0
}

type PlayerCatOuqiCat struct {
	CatId                *int32   `protobuf:"varint,1,opt,name=CatId" json:"CatId,omitempty"`
	Ouqi                 *int32   `protobuf:"varint,2,opt,name=Ouqi" json:"Ouqi,omitempty"`
	UpdateTime           *int32   `protobuf:"varint,3,opt,name=UpdateTime" json:"UpdateTime,omitempty"`
	HistoryTopRank       *int32   `protobuf:"varint,4,opt,name=HistoryTopRank" json:"HistoryTopRank,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerCatOuqiCat) Reset()         { *m = PlayerCatOuqiCat{} }
func (m *PlayerCatOuqiCat) String() string { return proto.CompactTextString(m) }
func (*PlayerCatOuqiCat) ProtoMessage()    {}
func (*PlayerCatOuqiCat) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{4}
}

func (m *PlayerCatOuqiCat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerCatOuqiCat.Unmarshal(m, b)
}
func (m *PlayerCatOuqiCat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerCatOuqiCat.Marshal(b, m, deterministic)
}
func (m *PlayerCatOuqiCat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerCatOuqiCat.Merge(m, src)
}
func (m *PlayerCatOuqiCat) XXX_Size() int {
	return xxx_messageInfo_PlayerCatOuqiCat.Size(m)
}
func (m *PlayerCatOuqiCat) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerCatOuqiCat.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerCatOuqiCat proto.InternalMessageInfo

func (m *PlayerCatOuqiCat) GetCatId() int32 {
	if m != nil && m.CatId != nil {
		return *m.CatId
	}
	return 0
}

func (m *PlayerCatOuqiCat) GetOuqi() int32 {
	if m != nil && m.Ouqi != nil {
		return *m.Ouqi
	}
	return 0
}

func (m *PlayerCatOuqiCat) GetUpdateTime() int32 {
	if m != nil && m.UpdateTime != nil {
		return *m.UpdateTime
	}
	return 0
}

func (m *PlayerCatOuqiCat) GetHistoryTopRank() int32 {
	if m != nil && m.HistoryTopRank != nil {
		return *m.HistoryTopRank
	}
	return 0
}

type PlayerCatOuqiCatList struct {
	List                 []*PlayerCatOuqiCat `protobuf:"bytes,1,rep,name=List" json:"List,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *PlayerCatOuqiCatList) Reset()         { *m = PlayerCatOuqiCatList{} }
func (m *PlayerCatOuqiCatList) String() string { return proto.CompactTextString(m) }
func (*PlayerCatOuqiCatList) ProtoMessage()    {}
func (*PlayerCatOuqiCatList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{5}
}

func (m *PlayerCatOuqiCatList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerCatOuqiCatList.Unmarshal(m, b)
}
func (m *PlayerCatOuqiCatList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerCatOuqiCatList.Marshal(b, m, deterministic)
}
func (m *PlayerCatOuqiCatList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerCatOuqiCatList.Merge(m, src)
}
func (m *PlayerCatOuqiCatList) XXX_Size() int {
	return xxx_messageInfo_PlayerCatOuqiCatList.Size(m)
}
func (m *PlayerCatOuqiCatList) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerCatOuqiCatList.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerCatOuqiCatList proto.InternalMessageInfo

func (m *PlayerCatOuqiCatList) GetList() []*PlayerCatOuqiCat {
	if m != nil {
		return m.List
	}
	return nil
}

type PlayerBeZanedHistoryTopData struct {
	Rank                 *int32   `protobuf:"varint,1,opt,name=Rank" json:"Rank,omitempty"`
	Zaned                *int32   `protobuf:"varint,2,opt,name=Zaned" json:"Zaned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerBeZanedHistoryTopData) Reset()         { *m = PlayerBeZanedHistoryTopData{} }
func (m *PlayerBeZanedHistoryTopData) String() string { return proto.CompactTextString(m) }
func (*PlayerBeZanedHistoryTopData) ProtoMessage()    {}
func (*PlayerBeZanedHistoryTopData) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d32edf7ae298984, []int{6}
}

func (m *PlayerBeZanedHistoryTopData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerBeZanedHistoryTopData.Unmarshal(m, b)
}
func (m *PlayerBeZanedHistoryTopData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerBeZanedHistoryTopData.Marshal(b, m, deterministic)
}
func (m *PlayerBeZanedHistoryTopData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerBeZanedHistoryTopData.Merge(m, src)
}
func (m *PlayerBeZanedHistoryTopData) XXX_Size() int {
	return xxx_messageInfo_PlayerBeZanedHistoryTopData.Size(m)
}
func (m *PlayerBeZanedHistoryTopData) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerBeZanedHistoryTopData.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerBeZanedHistoryTopData proto.InternalMessageInfo

func (m *PlayerBeZanedHistoryTopData) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *PlayerBeZanedHistoryTopData) GetZaned() int32 {
	if m != nil && m.Zaned != nil {
		return *m.Zaned
	}
	return 0
}

func init() {
	proto.RegisterType((*PlayerStageTotalScoreHistoryTopData)(nil), "db.PlayerStageTotalScoreHistoryTopData")
	proto.RegisterType((*PlayerStageTotalScoreStage)(nil), "db.PlayerStageTotalScoreStage")
	proto.RegisterType((*PlayerStageTotalScoreStageList)(nil), "db.PlayerStageTotalScoreStageList")
	proto.RegisterType((*PlayerCharmHistoryTopData)(nil), "db.PlayerCharmHistoryTopData")
	proto.RegisterType((*PlayerCatOuqiCat)(nil), "db.PlayerCatOuqiCat")
	proto.RegisterType((*PlayerCatOuqiCatList)(nil), "db.PlayerCatOuqiCatList")
	proto.RegisterType((*PlayerBeZanedHistoryTopData)(nil), "db.PlayerBeZanedHistoryTopData")
}

func init() { proto.RegisterFile("db_rpcsvr.proto", fileDescriptor_6d32edf7ae298984) }

var fileDescriptor_6d32edf7ae298984 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x49, 0x9a, 0x82, 0x8c, 0x50, 0x65, 0xc9, 0x21, 0x56, 0x08, 0x65, 0x05, 0xc9, 0x29,
	0x87, 0xfe, 0x02, 0xb1, 0x8a, 0x2d, 0x08, 0x95, 0x34, 0x5e, 0xbc, 0xc8, 0xa4, 0xbb, 0x68, 0xb0,
	0xed, 0xc6, 0xcd, 0x28, 0xf4, 0xee, 0x0f, 0x97, 0x4c, 0xd2, 0xa6, 0x16, 0x0b, 0x3d, 0xed, 0xbe,
	0x99, 0x6f, 0x1e, 0x33, 0x0f, 0xce, 0x54, 0xf6, 0x6a, 0x8b, 0x79, 0xf9, 0x6d, 0xe3, 0xc2, 0x1a,
	0x32, 0xc2, 0x55, 0x99, 0x9c, 0xc2, 0xd5, 0xd3, 0x02, 0xd7, 0xda, 0xce, 0x08, 0xdf, 0x74, 0x6a,
	0x08, 0x17, 0xb3, 0xb9, 0xb1, 0x7a, 0x9c, 0x97, 0x64, 0xec, 0x3a, 0x35, 0xc5, 0x1d, 0x12, 0x0a,
	0x01, 0x5e, 0x82, 0xab, 0x8f, 0xc0, 0x19, 0x38, 0x51, 0x37, 0xe1, 0xbf, 0xf0, 0xa1, 0xcb, 0x68,
	0xe0, 0x72, 0xb1, 0x16, 0x72, 0x0c, 0xfd, 0x7f, 0x0d, 0x59, 0x8a, 0x1e, 0xb8, 0x13, 0xd5, 0xb8,
	0xb8, 0x13, 0x25, 0xfa, 0x70, 0x92, 0x9a, 0x62, 0xd7, 0x66, 0xab, 0x65, 0x0a, 0xe1, 0x61, 0xa7,
	0xc7, 0xbc, 0x24, 0x31, 0x04, 0xaf, 0x7a, 0x03, 0x67, 0xd0, 0x89, 0x4e, 0x87, 0x61, 0xac, 0xb2,
	0xf8, 0xf0, 0x44, 0xc2, 0xac, 0xbc, 0x87, 0x8b, 0x9a, 0x19, 0xbd, 0xa3, 0x5d, 0x1e, 0x77, 0x26,
	0xa3, 0x9b, 0x33, 0x59, 0xc8, 0x1f, 0x07, 0xce, 0x1b, 0x1f, 0xa4, 0xe9, 0xd7, 0x67, 0x3e, 0x42,
	0x62, 0x14, 0x69, 0x7b, 0x60, 0x2d, 0x2a, 0xd3, 0x0a, 0x68, 0xe6, 0xf9, 0x2f, 0x42, 0x80, 0xe7,
	0x42, 0x21, 0xe9, 0x34, 0x5f, 0xea, 0xa0, 0xc3, 0x9d, 0x9d, 0x8a, 0xb8, 0x86, 0x5e, 0xbb, 0x1a,
	0xaf, 0xe4, 0x31, 0xb3, 0x57, 0x95, 0x37, 0xe0, 0xef, 0x6f, 0xc1, 0xc9, 0x44, 0x7f, 0x92, 0xf1,
	0xdb, 0x64, 0x5a, 0xae, 0xc9, 0xe3, 0x01, 0x2e, 0xeb, 0xce, 0xad, 0x7e, 0xc1, 0x95, 0x56, 0xc7,
	0x25, 0xc2, 0xe8, 0x26, 0x11, 0x16, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x60, 0x49, 0x30, 0x28,
	0x5f, 0x02, 0x00, 0x00,
}
