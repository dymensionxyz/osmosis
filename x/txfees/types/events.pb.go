// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/txfees/v1beta1/events.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
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

type EventChargeFee struct {
	Payer    string `protobuf:"bytes,1,opt,name=payer,proto3" json:"payer,omitempty"`
	TakerFee string `protobuf:"bytes,2,opt,name=taker_fee,json=takerFee,proto3" json:"taker_fee,omitempty"`
	// Beneficiary is the address that will receive the fee. Optional: may be empty.
	Beneficiary        string `protobuf:"bytes,3,opt,name=beneficiary,proto3" json:"beneficiary,omitempty"`
	BeneficiaryRevenue string `protobuf:"bytes,4,opt,name=beneficiary_revenue,json=beneficiaryRevenue,proto3" json:"beneficiary_revenue,omitempty"`
	CommunityPool      bool   `protobuf:"varint,5,opt,name=community_pool,json=communityPool,proto3" json:"community_pool,omitempty"`
}

func (m *EventChargeFee) Reset()         { *m = EventChargeFee{} }
func (m *EventChargeFee) String() string { return proto.CompactTextString(m) }
func (*EventChargeFee) ProtoMessage()    {}
func (*EventChargeFee) Descriptor() ([]byte, []int) {
	return fileDescriptor_fdb570c08d9ae603, []int{0}
}
func (m *EventChargeFee) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventChargeFee) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventChargeFee.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventChargeFee) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventChargeFee.Merge(m, src)
}
func (m *EventChargeFee) XXX_Size() int {
	return m.Size()
}
func (m *EventChargeFee) XXX_DiscardUnknown() {
	xxx_messageInfo_EventChargeFee.DiscardUnknown(m)
}

var xxx_messageInfo_EventChargeFee proto.InternalMessageInfo

func (m *EventChargeFee) GetPayer() string {
	if m != nil {
		return m.Payer
	}
	return ""
}

func (m *EventChargeFee) GetTakerFee() string {
	if m != nil {
		return m.TakerFee
	}
	return ""
}

func (m *EventChargeFee) GetBeneficiary() string {
	if m != nil {
		return m.Beneficiary
	}
	return ""
}

func (m *EventChargeFee) GetBeneficiaryRevenue() string {
	if m != nil {
		return m.BeneficiaryRevenue
	}
	return ""
}

func (m *EventChargeFee) GetCommunityPool() bool {
	if m != nil {
		return m.CommunityPool
	}
	return false
}

func init() {
	proto.RegisterType((*EventChargeFee)(nil), "dymensionxyz.dymension.txfees.v1beta1.EventChargeFee")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/txfees/v1beta1/events.proto", fileDescriptor_fdb570c08d9ae603)
}

var fileDescriptor_fdb570c08d9ae603 = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xb1, 0x4e, 0xeb, 0x30,
	0x18, 0x85, 0xeb, 0x7b, 0x29, 0x6a, 0x8d, 0xe8, 0x10, 0x3a, 0x44, 0x45, 0xb2, 0x2a, 0xa4, 0x4a,
	0x5d, 0x88, 0xd5, 0x22, 0x5e, 0x00, 0x04, 0x13, 0x03, 0xea, 0xc8, 0x52, 0xd9, 0xe1, 0x6f, 0x6a,
	0xd1, 0xf8, 0x8f, 0x6c, 0xb7, 0xaa, 0x79, 0x0a, 0x5e, 0x88, 0x9d, 0xb1, 0x23, 0x23, 0x6a, 0x5f,
	0x04, 0xc5, 0xa1, 0x51, 0xb6, 0x9c, 0xf3, 0xe5, 0xb3, 0xe5, 0x43, 0xa7, 0xaf, 0x3e, 0x07, 0x6d,
	0x15, 0xea, 0xad, 0x7f, 0xe7, 0x75, 0xe0, 0x6e, 0xbb, 0x00, 0xb0, 0x7c, 0x33, 0x91, 0xe0, 0xc4,
	0x84, 0xc3, 0x06, 0xb4, 0xb3, 0x49, 0x61, 0xd0, 0x61, 0x34, 0x6a, 0x3a, 0x49, 0x1d, 0x92, 0xca,
	0x49, 0xfe, 0x9c, 0x41, 0x3f, 0xc3, 0x0c, 0x83, 0xc1, 0xcb, 0xaf, 0x4a, 0x1e, 0xb0, 0x14, 0x6d,
	0x8e, 0x96, 0x4b, 0x61, 0xa1, 0x3e, 0x3e, 0x45, 0xa5, 0x2b, 0x7e, 0xf5, 0x49, 0x68, 0xef, 0xa1,
	0xbc, 0xed, 0x7e, 0x29, 0x4c, 0x06, 0x8f, 0x00, 0x51, 0x9f, 0xb6, 0x0b, 0xe1, 0xc1, 0xc4, 0x64,
	0x48, 0xc6, 0xdd, 0x59, 0x15, 0xa2, 0x4b, 0xda, 0x75, 0xe2, 0x0d, 0xcc, 0x7c, 0x01, 0x10, 0xff,
	0x0b, 0xa4, 0x13, 0x8a, 0x52, 0x19, 0xd2, 0x33, 0x09, 0x1a, 0x16, 0x2a, 0x55, 0xc2, 0xf8, 0xf8,
	0x7f, 0xc0, 0xcd, 0x2a, 0xe2, 0xf4, 0xa2, 0x11, 0xe7, 0xa6, 0x7c, 0xe1, 0x1a, 0xe2, 0x93, 0xf0,
	0x67, 0xd4, 0x40, 0xb3, 0x8a, 0x44, 0x23, 0xda, 0x4b, 0x31, 0xcf, 0xd7, 0x5a, 0x39, 0x3f, 0x2f,
	0x10, 0x57, 0x71, 0x7b, 0x48, 0xc6, 0x9d, 0xd9, 0x79, 0xdd, 0x3e, 0x23, 0xae, 0xee, 0x9e, 0xbe,
	0xf6, 0x8c, 0xec, 0xf6, 0x8c, 0xfc, 0xec, 0x19, 0xf9, 0x38, 0xb0, 0xd6, 0xee, 0xc0, 0x5a, 0xdf,
	0x07, 0xd6, 0x7a, 0x99, 0x66, 0xca, 0x2d, 0xd7, 0x32, 0x49, 0x31, 0xe7, 0x61, 0x03, 0x65, 0xaf,
	0x57, 0x42, 0xda, 0x63, 0xe0, 0x9b, 0xc9, 0x2d, 0xdf, 0x1e, 0x97, 0x77, 0xbe, 0x00, 0x2b, 0x4f,
	0xc3, 0x28, 0x37, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbb, 0x92, 0x02, 0x64, 0xa7, 0x01, 0x00,
	0x00,
}

func (m *EventChargeFee) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventChargeFee) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventChargeFee) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CommunityPool {
		i--
		if m.CommunityPool {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.BeneficiaryRevenue) > 0 {
		i -= len(m.BeneficiaryRevenue)
		copy(dAtA[i:], m.BeneficiaryRevenue)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.BeneficiaryRevenue)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Beneficiary) > 0 {
		i -= len(m.Beneficiary)
		copy(dAtA[i:], m.Beneficiary)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Beneficiary)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TakerFee) > 0 {
		i -= len(m.TakerFee)
		copy(dAtA[i:], m.TakerFee)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.TakerFee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Payer) > 0 {
		i -= len(m.Payer)
		copy(dAtA[i:], m.Payer)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Payer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventChargeFee) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Payer)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.TakerFee)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Beneficiary)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.BeneficiaryRevenue)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.CommunityPool {
		n += 2
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventChargeFee) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventChargeFee: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventChargeFee: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakerFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TakerFee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Beneficiary", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Beneficiary = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BeneficiaryRevenue", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BeneficiaryRevenue = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommunityPool", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
			m.CommunityPool = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
