package common

import (
	"errors"
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/mock"
)

const EthComplexSignature = "function,address,uint56,bytes24,bool,(uint256,uint256),(uint256[],bool,bytes,address),address[],string,uint256[],bytes[],bool[],(uint256,int256)[],(uint256[],bool,bytes,address)[],uint256[2][3],(uint256[2],uint256[3])[],(uint256,bytes,bool,address)[3][2],(uint256,bytes,bool,(address,uint256)[]),(uint256,bytes,bool,(address,uint256)[])[2]"

const MvxComplexSignature1 = "Address,BigInt,bytes,bool,tuple<u64,i32>,tuple<List<u64>,bool,utf-8 string,Address>,List<Address>,List<BigInt>,List<bytes>,List<bool>,List<tuple<u64,i32>>,List<tuple<List<u64>,bool,utf-8 string,Address>>,array3<array2<BigInt>>,List<tuple<array2<u64>,array3<i32>>>,array2<array3<tuple<u64,bytes,bool>>>,tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>,array2<tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>>,Option<tuple<List<BigInt>,TokenIdentifier,bool>>,Option<tuple<List<BigInt>,TokenIdentifier,bool>>,optional<List<BigInt>>"
const MvxComplexSignature2 = "Address,BigInt,bytes,bool,tuple<u64,i32>,tuple<List<u64>,bool,utf-8 string,Address>,List<Address>,List<BigInt>,List<bytes>,List<bool>,List<tuple<u64,i32>>,List<tuple<List<u64>,bool,utf-8 string,Address>>,array3<array2<BigInt>>,List<tuple<array2<u64>,array3<i32>>>,array2<array3<tuple<u64,bytes,bool>>>,tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>,array2<tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>>,Option<tuple<List<BigInt>,TokenIdentifier,bool>>,variadic<List<BigInt>>"
const MvxComplexSignature3 = "u8,u16,u32,i8,i16,i64,Address,BigInt,bytes,bool,tuple<u64,i32>,tuple<List<u64>,bool,utf-8 string,Address>,List<Address>,List<BigInt>,List<bytes>,List<bool>,List<tuple<u64,i32>>,List<tuple<List<u64>,bool,utf-8 string,Address>>,array3<array2<BigInt>>,List<tuple<array2<u64>,array3<i32>>>,array2<array3<tuple<u64,bytes,bool>>>,tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>,array2<tuple<u64,bytes,bool,List<tuple<Address,BigInt>>>>,Option<tuple<List<BigInt>,TokenIdentifier,bool>>,multi<List<BigInt>,BigUint>"

var errNotImplemented = errors.New("not implemented")

func BuildTestEncodingContext() *EncodingContext {
	return BuildEncodingContext(&mock.AccountsStub{
		RequestAddressCalled: func(request *vmcommon.AddressRequest) (*vmcommon.AddressResponse, error) {
			err := vmcommon.ValidateAddressRequest(request)
			if err != nil {
				return nil, err
			}

			if request.SourceIdentifier == core.MVXAddressIdentifier {
				if request.RequestedIdentifier == core.ETHAddressIdentifier {
					ethereumAddress := request.SourceAddress[12:]
					return &vmcommon.AddressResponse{MultiversXAddress: request.SourceAddress, RequestedAddress: ethereumAddress}, nil
				}
			}

			if request.SourceIdentifier == core.ETHAddressIdentifier {
				if request.RequestedIdentifier == core.MVXAddressIdentifier {
					multiversXAddress := append(make([]byte, 12), request.SourceAddress...)
					return &vmcommon.AddressResponse{MultiversXAddress: multiversXAddress, RequestedAddress: multiversXAddress}, nil
				}
			}

			return nil, errNotImplemented
		},
	})
}
