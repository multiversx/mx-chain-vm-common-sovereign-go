package eth

import ethAbi "github.com/ethereum/go-ethereum/accounts/abi"

const (
	Size8  = 8
	Size16 = 16
	Size32 = 32
	Size64 = 64

	BeginTuple = "("
	EndTuple   = ")"
	BeginArray = "["
	EndArray   = "]"

	Uint8   = "uint8"
	Uint16  = "uint16"
	Uint32  = "uint32"
	Uint64  = "uint64"
	Uint256 = "uint256"
	Int8    = "int8"
	Int16   = "int16"
	Int32   = "int32"
	Int64   = "int64"
	Int256  = "int256"
	Bool    = "bool"
	Bytes   = "bytes"
	Address = "address"
	String  = "string"
	Tuple   = "tuple"
)

var Uint8Type, _ = ethAbi.NewType(Uint8, Uint8, []ethAbi.ArgumentMarshaling{})

func ExtractElem(ethArgument *ethAbi.Type) *ethAbi.Type {
	switch ethArgument.T {
	case ethAbi.FunctionTy, ethAbi.FixedBytesTy:
		return &Uint8Type
	default:
		return ethArgument.Elem
	}
}
