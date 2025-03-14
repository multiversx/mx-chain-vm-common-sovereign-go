package builtInFunctions

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestBaseComponentsHolder_addNFTToDestination(t *testing.T) {
	t.Parallel()

	t.Run("different shards should save liquidity to system account", func(t *testing.T) {
		t.Parallel()

		saveCalled := false
		addToLiquiditySystemAccCalled := false
		b := &baseComponentsHolder{
			esdtStorageHandler: &mock.ESDTNFTStorageHandlerStub{
				GetESDTNFTTokenOnDestinationCalled: func(_ vmcommon.UserAccountHandler, _ []byte, _ uint64) (*esdt.ESDigitalToken, bool, error) {
					return &esdt.ESDigitalToken{
						Value: big.NewInt(100),
					}, false, nil
				},
				SaveESDTNFTTokenCalled: func(_ []byte, _ vmcommon.UserAccountHandler, _ []byte, _ uint64, esdtData *esdt.ESDigitalToken, properties vmcommon.NftSaveArgs) ([]byte, error) {
					assert.Equal(t, big.NewInt(200), esdtData.Value)
					saveCalled = true
					return nil, nil
				},
				AddToLiquiditySystemAccCalled: func(esdtTokenKey []byte, _ uint32, nonce uint64, transferValue *big.Int, _ bool) error {
					assert.Equal(t, big.NewInt(100), transferValue)
					addToLiquiditySystemAccCalled = true
					return nil
				},
			},
			globalSettingsHandler: &mock.GlobalSettingsHandlerStub{
				IsPausedCalled: func(_ []byte) bool {
					return false
				},
			},
			shardCoordinator: &mock.ShardCoordinatorStub{
				SameShardCalled: func(_, _ []byte) bool {
					return false
				},
			},
			enableEpochsHandler: &mock.EnableEpochsHandlerStub{},
		}

		acc := &mock.UserAccountStub{}
		esdtDataToTransfer := &esdt.ESDigitalToken{
			Type:       0,
			Value:      big.NewInt(100),
			Properties: make([]byte, 0),
		}
		err := b.addNFTToDestination([]byte("sndAddr"), []byte("dstAddr"), acc, esdtDataToTransfer, []byte("esdtTokenKey"), 0, false, false)
		assert.Nil(t, err)
		assert.True(t, addToLiquiditySystemAccCalled)
		assert.True(t, saveCalled)
	})
}

func TestBaseComponentsHolder_getLatestEsdtData(t *testing.T) {
	t.Parallel()

	t.Run("flag disabled should return transfer esdt data", func(t *testing.T) {
		t.Parallel()

		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(_ core.EnableEpochFlag) bool {
				return false
			},
		}
		currentEsdtData := &esdt.ESDigitalToken{
			Reserved: []byte{1},
			Value:    big.NewInt(100),
		}
		transferEsdtData := &esdt.ESDigitalToken{
			Reserved: []byte{2},
			Value:    big.NewInt(200),
		}

		latestEsdtData, err := getLatestMetaData(currentEsdtData, transferEsdtData, enableEpochsHandler, &mock.MarshalizerMock{})
		assert.Nil(t, err)
		assert.Equal(t, transferEsdtData, latestEsdtData)
	})
	t.Run("flag enabled and transfer esdt data version is not set should merge", func(t *testing.T) {
		t.Parallel()
		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(_ core.EnableEpochFlag) bool {
				return true
			},
		}
		name := []byte("name")
		creator := []byte("creator")
		newCreator := []byte("newCreator")
		royalties := uint32(25)
		newRoyalties := uint32(11)
		hash := []byte("hash")
		uris := [][]byte{[]byte("uri1"), []byte("uri2")}
		attributes := []byte("attributes")
		newAttributes := []byte("newAttributes")
		transferEsdtData := &esdt.ESDigitalToken{
			Reserved: []byte{1},
			TokenMetaData: &esdt.MetaData{
				Nonce:      0,
				Name:       name,
				Creator:    creator,
				Royalties:  royalties,
				Hash:       hash,
				URIs:       uris,
				Attributes: attributes,
			},
		}
		currentEsdtVersion := &esdt.MetaDataVersion{
			Creator:    2,
			Royalties:  2,
			Attributes: 2,
		}
		versionBytes, _ := (&mock.MarshalizerMock{}).Marshal(currentEsdtVersion)
		currentEsdtData := &esdt.ESDigitalToken{
			Reserved: versionBytes,
			TokenMetaData: &esdt.MetaData{
				Creator:    newCreator,
				Royalties:  newRoyalties,
				Attributes: newAttributes,
			},
		}

		latestEsdtData, err := getLatestMetaData(currentEsdtData, transferEsdtData, enableEpochsHandler, &mock.MarshalizerMock{})
		assert.Nil(t, err)
		assert.Equal(t, versionBytes, latestEsdtData.Reserved)
		assert.Equal(t, newCreator, latestEsdtData.TokenMetaData.Creator)
		assert.Equal(t, newRoyalties, latestEsdtData.TokenMetaData.Royalties)
		assert.Equal(t, newAttributes, latestEsdtData.TokenMetaData.Attributes)

		assert.Equal(t, name, latestEsdtData.TokenMetaData.Name)
		assert.Equal(t, hash, latestEsdtData.TokenMetaData.Hash)
		assert.Equal(t, uris, latestEsdtData.TokenMetaData.URIs)
	})
	t.Run("different versions for different fields should merge", func(t *testing.T) {
		t.Parallel()
		enableEpochsHandler := &mock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(_ core.EnableEpochFlag) bool {
				return true
			},
		}
		name := []byte("name")
		creator := []byte("creator")
		newCreator := []byte("newCreator")
		royalties := uint32(25)
		newRoyalties := uint32(11)
		hash := []byte("hash")
		uris := [][]byte{[]byte("uri1"), []byte("uri2")}
		attributes := []byte("attributes")
		newAttributes := []byte("newAttributes")
		transferEsdtVersion := &esdt.MetaDataVersion{
			Name:       3,
			Creator:    0,
			Royalties:  0,
			Hash:       3,
			URIs:       3,
			Attributes: 3,
		}
		versionBytes, _ := (&mock.MarshalizerMock{}).Marshal(transferEsdtVersion)
		transferEsdtData := &esdt.ESDigitalToken{
			Reserved: versionBytes,
			TokenMetaData: &esdt.MetaData{
				Nonce:      0,
				Name:       name,
				Creator:    creator,
				Royalties:  royalties,
				Hash:       hash,
				URIs:       uris,
				Attributes: attributes,
			},
		}
		currentEsdtVersion := &esdt.MetaDataVersion{
			Name:       0,
			Creator:    2,
			Royalties:  2,
			Hash:       0,
			URIs:       0,
			Attributes: 2,
		}
		versionBytes, _ = (&mock.MarshalizerMock{}).Marshal(currentEsdtVersion)
		currentEsdtData := &esdt.ESDigitalToken{
			Reserved: versionBytes,
			TokenMetaData: &esdt.MetaData{
				Creator:    newCreator,
				Royalties:  newRoyalties,
				Attributes: newAttributes,
			},
		}

		latestEsdtData, err := getLatestMetaData(currentEsdtData, transferEsdtData, enableEpochsHandler, &mock.MarshalizerMock{})
		assert.Nil(t, err)
		expectedVersion := &esdt.MetaDataVersion{
			Name:       3,
			Creator:    2,
			Royalties:  2,
			Hash:       3,
			URIs:       3,
			Attributes: 3,
		}
		expectedVersionBytes, _ := (&mock.MarshalizerMock{}).Marshal(expectedVersion)
		assert.Equal(t, expectedVersionBytes, latestEsdtData.Reserved)
		assert.Equal(t, newCreator, latestEsdtData.TokenMetaData.Creator)
		assert.Equal(t, newRoyalties, latestEsdtData.TokenMetaData.Royalties)
		assert.Equal(t, name, latestEsdtData.TokenMetaData.Name)
		assert.Equal(t, hash, latestEsdtData.TokenMetaData.Hash)
		assert.Equal(t, uris, latestEsdtData.TokenMetaData.URIs)
		assert.Equal(t, attributes, latestEsdtData.TokenMetaData.Attributes)
	})
}
