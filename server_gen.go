package remote

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *BeginTransactionRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Writable":
			z.Writable, err = dc.ReadBool()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z BeginTransactionRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Writable"
	err = en.Append(0x81, 0xa8, 0x57, 0x72, 0x69, 0x74, 0x61, 0x62, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.Writable)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BeginTransactionRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Writable"
	o = append(o, 0x81, 0xa8, 0x57, 0x72, 0x69, 0x74, 0x61, 0x62, 0x6c, 0x65)
	o = msgp.AppendBool(o, z.Writable)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BeginTransactionRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Writable":
			z.Writable, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z BeginTransactionRequest) Msgsize() (s int) {
	s = 1 + 9 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BeginTransactionResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z BeginTransactionResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "ContextID"
	err = en.Append(0x81, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ContextID)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BeginTransactionResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "ContextID"
	o = append(o, 0x81, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ContextID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BeginTransactionResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z BeginTransactionResponse) Msgsize() (s int) {
	s = 1 + 10 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BucketRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Key":
			z.Key, err = dc.ReadBytes(z.Key)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *BucketRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ContextID"
	err = en.Append(0x82, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ContextID)
	if err != nil {
		return
	}
	// write "Key"
	err = en.Append(0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Key)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BucketRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ContextID"
	o = append(o, 0x82, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ContextID)
	// string "Key"
	o = append(o, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, z.Key)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BucketRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Key":
			z.Key, bts, err = msgp.ReadBytesBytes(bts, z.Key)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *BucketRequest) Msgsize() (s int) {
	s = 1 + 10 + msgp.Uint64Size + 4 + msgp.BytesPrefixSize + len(z.Key)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BucketResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "BucketContextID":
			z.BucketContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z BucketResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "BucketID"
	err = en.Append(0x82, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.BucketID)
	if err != nil {
		return
	}
	// write "BucketContextID"
	err = en.Append(0xaf, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.BucketContextID)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BucketResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "BucketID"
	o = append(o, 0x82, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.BucketID)
	// string "BucketContextID"
	o = append(o, 0xaf, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.BucketContextID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BucketResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "BucketContextID":
			z.BucketContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z BucketResponse) Msgsize() (s int) {
	s = 1 + 9 + msgp.Uint64Size + 16 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BucketStatsRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "BucketID":
			z.BucketID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z BucketStatsRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ContextID"
	err = en.Append(0x82, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ContextID)
	if err != nil {
		return
	}
	// write "BucketID"
	err = en.Append(0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.BucketID)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z BucketStatsRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ContextID"
	o = append(o, 0x82, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ContextID)
	// string "BucketID"
	o = append(o, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.BucketID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BucketStatsRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ContextID":
			z.ContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "BucketID":
			z.BucketID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z BucketStatsRequest) Msgsize() (s int) {
	s = 1 + 10 + msgp.Uint64Size + 9 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *GetReqeust) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ContextID":
			z.ContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Key":
			z.Key, err = dc.ReadBytes(z.Key)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *GetReqeust) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "BucketID"
	err = en.Append(0x83, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.BucketID)
	if err != nil {
		return
	}
	// write "ContextID"
	err = en.Append(0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ContextID)
	if err != nil {
		return
	}
	// write "Key"
	err = en.Append(0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Key)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GetReqeust) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "BucketID"
	o = append(o, 0x83, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.BucketID)
	// string "ContextID"
	o = append(o, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ContextID)
	// string "Key"
	o = append(o, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, z.Key)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GetReqeust) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ContextID":
			z.ContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Key":
			z.Key, bts, err = msgp.ReadBytesBytes(bts, z.Key)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *GetReqeust) Msgsize() (s int) {
	s = 1 + 9 + msgp.Uint64Size + 10 + msgp.Uint64Size + 4 + msgp.BytesPrefixSize + len(z.Key)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *GetResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Val":
			z.Val, err = dc.ReadBytes(z.Val)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *GetResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Val"
	err = en.Append(0x81, 0xa3, 0x56, 0x61, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Val)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GetResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Val"
	o = append(o, 0x81, 0xa3, 0x56, 0x61, 0x6c)
	o = msgp.AppendBytes(o, z.Val)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GetResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Val":
			z.Val, bts, err = msgp.ReadBytesBytes(bts, z.Val)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *GetResponse) Msgsize() (s int) {
	s = 1 + 4 + msgp.BytesPrefixSize + len(z.Val)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PutReqeust) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "ContextID":
			z.ContextID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Key":
			z.Key, err = dc.ReadBytes(z.Key)
			if err != nil {
				return
			}
		case "Val":
			z.Val, err = dc.ReadBytes(z.Val)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PutReqeust) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "BucketID"
	err = en.Append(0x84, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.BucketID)
	if err != nil {
		return
	}
	// write "ContextID"
	err = en.Append(0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ContextID)
	if err != nil {
		return
	}
	// write "Key"
	err = en.Append(0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Key)
	if err != nil {
		return
	}
	// write "Val"
	err = en.Append(0xa3, 0x56, 0x61, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Val)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PutReqeust) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "BucketID"
	o = append(o, 0x84, 0xa8, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.BucketID)
	// string "ContextID"
	o = append(o, 0xa9, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ContextID)
	// string "Key"
	o = append(o, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, z.Key)
	// string "Val"
	o = append(o, 0xa3, 0x56, 0x61, 0x6c)
	o = msgp.AppendBytes(o, z.Val)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PutReqeust) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "BucketID":
			z.BucketID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "ContextID":
			z.ContextID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Key":
			z.Key, bts, err = msgp.ReadBytesBytes(bts, z.Key)
			if err != nil {
				return
			}
		case "Val":
			z.Val, bts, err = msgp.ReadBytesBytes(bts, z.Val)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *PutReqeust) Msgsize() (s int) {
	s = 1 + 9 + msgp.Uint64Size + 10 + msgp.Uint64Size + 4 + msgp.BytesPrefixSize + len(z.Key) + 4 + msgp.BytesPrefixSize + len(z.Val)
	return
}
