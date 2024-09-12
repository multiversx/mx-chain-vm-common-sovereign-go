package vmcommon

import "errors"

// ErrSubtractionOverflow signals that uint64 subtraction overflowed
var ErrSubtractionOverflow = errors.New("uint64 subtraction overflowed")

// ErrAsyncParams signals that there was an error with the async parameters
var ErrAsyncParams = errors.New("async parameters error")

// ErrInvalidVMType signals that invalid vm type was provided
var ErrInvalidVMType = errors.New("invalid VM type")

// ErrTransfersNotIndexed signals that transfers were found unindexed
var ErrTransfersNotIndexed = errors.New("unindexed transfers found")

// ErrNilTransferIndexer signals that the provided transfer indexer is nil
var ErrNilTransferIndexer = errors.New("nil NextOutputTransferIndexProvider")

// ErrNilRequest signals that the provided request is nil
var ErrNilRequest = errors.New("nil request")

// ErrInvalidSourceIdentifier signals that the source identifier is invalid
var ErrInvalidSourceIdentifier = errors.New("invalid source identifier")

// ErrInvalidRequestedIdentifier signals that the requested identifier is invalid
var ErrInvalidRequestedIdentifier = errors.New("invalid requested identifier")

// ErrSourceIdentifierMatchesRequestedIdentifier signals that the source identifier matches the requested identifier
var ErrSourceIdentifierMatchesRequestedIdentifier = errors.New("source identifier matches requested identifier")
