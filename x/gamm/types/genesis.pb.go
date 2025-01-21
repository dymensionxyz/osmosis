// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/gamm/v1beta1/genesis.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the gamm module's genesis state.
type GenesisState struct {
	Pools []*types.Any `protobuf:"bytes,1,rep,name=pools,proto3" json:"pools,omitempty"`
	// will be renamed to next_pool_id in an upcoming version
	NextPoolNumber uint64 `protobuf:"varint,2,opt,name=next_pool_number,json=nextPoolNumber,proto3" json:"next_pool_number,omitempty"`
	Params         Params `protobuf:"bytes,3,opt,name=params,proto3" json:"params"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc3b6373232d6d98, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetPools() []*types.Any {
	if m != nil {
		return m.Pools
	}
	return nil
}

func (m *GenesisState) GetNextPoolNumber() uint64 {
	if m != nil {
		return m.NextPoolNumber
	}
	return 0
}

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

type Params struct {
	PoolCreationFee      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=pool_creation_fee,json=poolCreationFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"pool_creation_fee" yaml:"pool_creation_fee"`
	EnableGlobalPoolFees bool                                     `protobuf:"varint,2,opt,name=enable_global_pool_fees,json=enableGlobalPoolFees,proto3" json:"enable_global_pool_fees,omitempty"`
	GlobalFees           GlobalFees                               `protobuf:"bytes,3,opt,name=global_fees,json=globalFees,proto3" json:"global_fees" yaml:"global_fees"`
	TakerFee             cosmossdk_io_math.LegacyDec              `protobuf:"bytes,4,opt,name=taker_fee,json=takerFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"taker_fee" yaml:"taker_fee"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc3b6373232d6d98, []int{1}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetPoolCreationFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.PoolCreationFee
	}
	return nil
}

func (m *Params) GetEnableGlobalPoolFees() bool {
	if m != nil {
		return m.EnableGlobalPoolFees
	}
	return false
}

func (m *Params) GetGlobalFees() GlobalFees {
	if m != nil {
		return m.GlobalFees
	}
	return GlobalFees{}
}

type GlobalFees struct {
	SwapFee cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=swap_fee,json=swapFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"swap_fee" yaml:"swap_fee"`
	ExitFee cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=exit_fee,json=exitFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"exit_fee" yaml:"exit_fee"`
}

func (m *GlobalFees) Reset()         { *m = GlobalFees{} }
func (m *GlobalFees) String() string { return proto.CompactTextString(m) }
func (*GlobalFees) ProtoMessage()    {}
func (*GlobalFees) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc3b6373232d6d98, []int{2}
}
func (m *GlobalFees) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GlobalFees) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GlobalFees.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GlobalFees) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalFees.Merge(m, src)
}
func (m *GlobalFees) XXX_Size() int {
	return m.Size()
}
func (m *GlobalFees) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalFees.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalFees proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GenesisState)(nil), "dymensionxyz.dymension.gamm.v1beta1.GenesisState")
	proto.RegisterType((*Params)(nil), "dymensionxyz.dymension.gamm.v1beta1.Params")
	proto.RegisterType((*GlobalFees)(nil), "dymensionxyz.dymension.gamm.v1beta1.GlobalFees")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/gamm/v1beta1/genesis.proto", fileDescriptor_cc3b6373232d6d98)
}

var fileDescriptor_cc3b6373232d6d98 = []byte{
	// 588 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0x8d, 0xfb, 0xf7, 0xa5, 0xd3, 0x4f, 0xb4, 0x58, 0x95, 0x48, 0x83, 0xe4, 0x44, 0x66, 0x63,
	0x09, 0x65, 0x86, 0x14, 0x75, 0xc3, 0x0e, 0x37, 0xa4, 0x2a, 0xaa, 0x50, 0x65, 0x76, 0x6c, 0xac,
	0xb1, 0x73, 0xe3, 0x5a, 0xb1, 0x3d, 0x51, 0xc6, 0x29, 0x31, 0x4f, 0x81, 0xc4, 0x43, 0x20, 0xb1,
	0xee, 0x8a, 0x27, 0xa8, 0x58, 0xa0, 0x8a, 0x15, 0x62, 0x11, 0x50, 0xf2, 0x06, 0x7d, 0x02, 0x34,
	0x3f, 0x31, 0x11, 0xdd, 0x44, 0xac, 0x92, 0xe3, 0x73, 0xcf, 0x39, 0x73, 0x4f, 0x26, 0x46, 0xed,
	0x5e, 0x91, 0x42, 0xc6, 0x63, 0x96, 0x4d, 0x8a, 0x77, 0xa4, 0x04, 0x24, 0xa2, 0x69, 0x4a, 0x2e,
	0xdb, 0x01, 0xe4, 0xb4, 0x4d, 0x22, 0xc8, 0x80, 0xc7, 0x1c, 0x0f, 0x47, 0x2c, 0x67, 0xe6, 0xa3,
	0x65, 0x09, 0x2e, 0x01, 0x16, 0x12, 0xac, 0x25, 0xf5, 0xfd, 0x88, 0x45, 0x4c, 0xce, 0x13, 0xf1,
	0x4d, 0x49, 0xeb, 0x07, 0x11, 0x63, 0x51, 0x02, 0x44, 0xa2, 0x60, 0xdc, 0x27, 0x34, 0x2b, 0x16,
	0x54, 0xc8, 0x78, 0xca, 0xb8, 0xaf, 0x34, 0x0a, 0x68, 0xca, 0x52, 0x88, 0x04, 0x94, 0x43, 0x79,
	0xa6, 0x90, 0xc5, 0x99, 0xe2, 0xed, 0xcf, 0x06, 0xfa, 0xff, 0x44, 0x1d, 0xf1, 0x75, 0x4e, 0x73,
	0x30, 0x8f, 0xd0, 0xe6, 0x90, 0xb1, 0x84, 0xd7, 0x8c, 0xe6, 0xba, 0xb3, 0x73, 0xb8, 0x8f, 0x55,
	0x2c, 0x5e, 0xc4, 0xe2, 0xe7, 0x59, 0xe1, 0x6e, 0x7f, 0xb9, 0x6a, 0x6d, 0x9e, 0x33, 0x96, 0x9c,
	0x7a, 0x6a, 0xda, 0x74, 0xd0, 0x5e, 0x06, 0x93, 0xdc, 0x17, 0xc8, 0xcf, 0xc6, 0x69, 0x00, 0xa3,
	0xda, 0x5a, 0xd3, 0x70, 0x36, 0xbc, 0x7b, 0xe2, 0xb9, 0x98, 0x7d, 0x25, 0x9f, 0x9a, 0xa7, 0x68,
	0x6b, 0x48, 0x47, 0x34, 0xe5, 0xb5, 0xf5, 0xa6, 0xe1, 0xec, 0x1c, 0x3e, 0xc6, 0x2b, 0x74, 0x82,
	0xcf, 0xa5, 0xc4, 0xdd, 0xb8, 0x9e, 0x36, 0x2a, 0x9e, 0x36, 0xb0, 0x3f, 0xae, 0xa3, 0x2d, 0x45,
	0x98, 0x1f, 0x0c, 0x74, 0x5f, 0x66, 0x87, 0x23, 0xa0, 0x79, 0xcc, 0x32, 0xbf, 0x0f, 0xa0, 0x77,
	0x38, 0xc0, 0xba, 0x12, 0x51, 0x42, 0xe9, 0x78, 0xcc, 0xe2, 0xcc, 0x3d, 0x13, 0x7e, 0xb7, 0xd3,
	0x46, 0xad, 0xa0, 0x69, 0xf2, 0xcc, 0xbe, 0xe3, 0x60, 0x7f, 0xfa, 0xd9, 0x70, 0xa2, 0x38, 0xbf,
	0x18, 0x07, 0x38, 0x64, 0xa9, 0xee, 0x56, 0x7f, 0xb4, 0x78, 0x6f, 0x40, 0xf2, 0x62, 0x08, 0x5c,
	0x9a, 0x71, 0x6f, 0x57, 0xe8, 0x8f, 0xb5, 0xbc, 0x0b, 0xa2, 0xcc, 0x07, 0x90, 0xd1, 0x20, 0x01,
	0x3f, 0x4a, 0x58, 0x40, 0x13, 0x55, 0x4f, 0x1f, 0x80, 0xcb, 0x72, 0xaa, 0xde, 0xbe, 0xa2, 0x4f,
	0x24, 0x2b, 0x4a, 0xea, 0x02, 0x70, 0x33, 0x41, 0x3b, 0x7a, 0x5e, 0x8e, 0xaa, 0x9e, 0xc8, 0x4a,
	0x3d, 0x29, 0x27, 0xe1, 0xe2, 0xd6, 0xf5, 0x6e, 0xa6, 0xda, 0x6d, 0xc9, 0xd1, 0xf6, 0x50, 0x54,
	0xce, 0x99, 0x01, 0xda, 0xce, 0xe9, 0x00, 0x46, 0xb2, 0xb1, 0x8d, 0xa6, 0xe1, 0x6c, 0xbb, 0x2f,
	0x84, 0xf4, 0xc7, 0xb4, 0xf1, 0x50, 0x2d, 0xca, 0x7b, 0x03, 0x1c, 0x33, 0x92, 0xd2, 0xfc, 0x02,
	0x9f, 0x41, 0x44, 0xc3, 0xa2, 0x03, 0xe1, 0xed, 0xb4, 0xb1, 0xa7, 0x9c, 0x4b, 0xb5, 0xfd, 0xed,
	0xaa, 0x85, 0x74, 0xd7, 0x1d, 0x08, 0xbd, 0xaa, 0x64, 0xba, 0x00, 0xf6, 0x57, 0x03, 0xa1, 0x3f,
	0x47, 0x33, 0x7d, 0x54, 0xe5, 0x6f, 0xe9, 0x50, 0xff, 0x46, 0x22, 0xb1, 0xb3, 0x5a, 0xe2, 0xae,
	0x4a, 0x5c, 0x88, 0xff, 0x0e, 0xfc, 0x4f, 0x10, 0xa2, 0x78, 0x1f, 0x55, 0x61, 0x12, 0xe7, 0x32,
	0x60, 0xed, 0x1f, 0x02, 0x16, 0xe2, 0x3b, 0x01, 0x82, 0xe8, 0x02, 0xb8, 0x2f, 0xaf, 0x67, 0x96,
	0x71, 0x33, 0xb3, 0x8c, 0x5f, 0x33, 0xcb, 0x78, 0x3f, 0xb7, 0x2a, 0x37, 0x73, 0xab, 0xf2, 0x7d,
	0x6e, 0x55, 0xde, 0x3c, 0x59, 0xba, 0x2e, 0x52, 0x18, 0xf3, 0x56, 0x42, 0x03, 0xbe, 0x00, 0xe4,
	0xb2, 0x7d, 0x44, 0x26, 0xea, 0x25, 0x21, 0x2f, 0x4f, 0xb0, 0x25, 0xff, 0x5b, 0x4f, 0x7f, 0x07,
	0x00, 0x00, 0xff, 0xff, 0x38, 0x2a, 0xd5, 0x7a, 0x50, 0x04, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.NextPoolNumber != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.NextPoolNumber))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TakerFee.Size()
		i -= size
		if _, err := m.TakerFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.GlobalFees.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.EnableGlobalPoolFees {
		i--
		if m.EnableGlobalPoolFees {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.PoolCreationFee) > 0 {
		for iNdEx := len(m.PoolCreationFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolCreationFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *GlobalFees) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GlobalFees) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GlobalFees) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ExitFee.Size()
		i -= size
		if _, err := m.ExitFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.SwapFee.Size()
		i -= size
		if _, err := m.SwapFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.NextPoolNumber != 0 {
		n += 1 + sovGenesis(uint64(m.NextPoolNumber))
	}
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PoolCreationFee) > 0 {
		for _, e := range m.PoolCreationFee {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.EnableGlobalPoolFees {
		n += 2
	}
	l = m.GlobalFees.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.TakerFee.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *GlobalFees) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.SwapFee.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.ExitFee.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pools = append(m.Pools, &types.Any{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextPoolNumber", wireType)
			}
			m.NextPoolNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextPoolNumber |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolCreationFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolCreationFee = append(m.PoolCreationFee, types1.Coin{})
			if err := m.PoolCreationFee[len(m.PoolCreationFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableGlobalPoolFees", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.EnableGlobalPoolFees = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalFees", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GlobalFees.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakerFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakerFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GlobalFees) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GlobalFees: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GlobalFees: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SwapFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SwapFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExitFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExitFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
