module github.com/multiversx/mx-chain-vm-common-go

replace github.com/multiversx/mx-chain-core-go => github.com/multiversx/mx-chain-core-sovereign-go v1.0.0-sov

go 1.20

require (
	github.com/mitchellh/mapstructure v1.4.1
	github.com/multiversx/mx-chain-core-go v1.2.24-0.20241119082458-e2451e147ab1
	github.com/multiversx/mx-chain-logger-go v1.0.15-0.20240508072523-3f00a726af57
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/gogo/protobuf => github.com/multiversx/protobuf v1.3.2
