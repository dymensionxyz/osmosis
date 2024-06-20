// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/txfees/v1beta1/feetoken.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

// FeeToken is a struct that specifies a coin denom, and pool ID pair.
// This marks the token as eligible for use as a tx fee asset in Osmosis.
// Its price in osmo is derived through looking at the provided pool ID.
// The pool ID must have base denom as one of its assets.
type FeeToken struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	PoolID uint64 `protobuf:"varint,2,opt,name=poolID,proto3" json:"poolID,omitempty" yaml:"pool_id"`
}

func (m *FeeToken) Reset()         { *m = FeeToken{} }
func (m *FeeToken) String() string { return proto.CompactTextString(m) }
func (*FeeToken) ProtoMessage()    {}
func (*FeeToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca4a790beba5662b, []int{0}
}
func (m *FeeToken) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeeToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeeToken.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeeToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeeToken.Merge(m, src)
}
func (m *FeeToken) XXX_Size() int {
	return m.Size()
}
func (m *FeeToken) XXX_DiscardUnknown() {
	xxx_messageInfo_FeeToken.DiscardUnknown(m)
}

var xxx_messageInfo_FeeToken proto.InternalMessageInfo

func (m *FeeToken) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *FeeToken) GetPoolID() uint64 {
	if m != nil {
		return m.PoolID
	}
	return 0
}

func init() {
	proto.RegisterType((*FeeToken)(nil), "dymensionxyz.dymension.txfees.v1beta1.FeeToken")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/txfees/v1beta1/feetoken.proto", fileDescriptor_ca4a790beba5662b)
}

var fileDescriptor_ca4a790beba5662b = []byte{
	// 252 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x49, 0xa9, 0xcc, 0x4d,
	0xcd, 0x2b, 0xce, 0xcc, 0xcf, 0xab, 0xa8, 0xac, 0xd2, 0x87, 0x73, 0xf4, 0x4b, 0x2a, 0xd2, 0x52,
	0x53, 0x8b, 0xf5, 0xcb, 0x0c, 0x93, 0x52, 0x4b, 0x12, 0x0d, 0xf5, 0xd3, 0x52, 0x53, 0x4b, 0xf2,
	0xb3, 0x53, 0xf3, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x54, 0x91, 0x75, 0xe9, 0xc1, 0x39,
	0x7a, 0x10, 0x5d, 0x7a, 0x50, 0x5d, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x1d, 0xfa, 0x20,
	0x16, 0x44, 0xb3, 0x52, 0x0a, 0x17, 0x87, 0x5b, 0x6a, 0x6a, 0x08, 0xc8, 0x38, 0x21, 0x35, 0x2e,
	0xd6, 0x94, 0xd4, 0xbc, 0xfc, 0x5c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x4e, 0x27, 0x81, 0x4f, 0xf7,
	0xe4, 0x79, 0x2a, 0x13, 0x73, 0x73, 0xac, 0x94, 0xc0, 0xc2, 0x4a, 0x41, 0x10, 0x69, 0x21, 0x2d,
	0x2e, 0xb6, 0x82, 0xfc, 0xfc, 0x1c, 0x4f, 0x17, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x16, 0x27, 0xa1,
	0x4f, 0xf7, 0xe4, 0xf9, 0x20, 0x0a, 0x41, 0xe2, 0xf1, 0x99, 0x29, 0x4a, 0x41, 0x50, 0x15, 0x56,
	0x2c, 0x2f, 0x16, 0xc8, 0x33, 0x3a, 0xf9, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3,
	0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c,
	0x43, 0x94, 0x51, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x7e, 0x71,
	0x6e, 0x7e, 0x71, 0x66, 0xb1, 0x6e, 0x4e, 0x62, 0x52, 0x31, 0x8c, 0xa3, 0x5f, 0x66, 0x68, 0xaa,
	0x5f, 0x01, 0x0b, 0x81, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0xd3, 0x8d, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x7b, 0x93, 0xa7, 0x35, 0x2f, 0x01, 0x00, 0x00,
}

func (this *FeeToken) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*FeeToken)
	if !ok {
		that2, ok := that.(FeeToken)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Denom != that1.Denom {
		return false
	}
	if this.PoolID != that1.PoolID {
		return false
	}
	return true
}
func (m *FeeToken) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeeToken) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeeToken) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PoolID != 0 {
		i = encodeVarintFeetoken(dAtA, i, uint64(m.PoolID))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintFeetoken(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintFeetoken(dAtA []byte, offset int, v uint64) int {
	offset -= sovFeetoken(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FeeToken) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovFeetoken(uint64(l))
	}
	if m.PoolID != 0 {
		n += 1 + sovFeetoken(uint64(m.PoolID))
	}
	return n
}

func sovFeetoken(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFeetoken(x uint64) (n int) {
	return sovFeetoken(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FeeToken) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeetoken
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
			return fmt.Errorf("proto: FeeToken: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeeToken: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeetoken
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
				return ErrInvalidLengthFeetoken
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeetoken
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolID", wireType)
			}
			m.PoolID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeetoken
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeetoken(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeetoken
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
func skipFeetoken(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFeetoken
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
					return 0, ErrIntOverflowFeetoken
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
					return 0, ErrIntOverflowFeetoken
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
				return 0, ErrInvalidLengthFeetoken
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFeetoken
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFeetoken
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFeetoken        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFeetoken          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFeetoken = fmt.Errorf("proto: unexpected end of group")
)
