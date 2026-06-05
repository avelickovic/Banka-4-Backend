package service

import "context"

// BankingServiceClient wraps the banking-service gRPC operations needed by
// the outbox worker to finalize the banking-side of a regular inter-bank
// payment after receiving a YES/NO vote.
//
// The full implementation (committing/rolling back a banking transaction by
// its banking TX ID) requires a proto RPC that does not yet exist. Until
// that RPC is added and the banking-service implements it, the stub in
// client/grpc/outbox_banking_client_impl.go returns "not implemented".
// OTC transactions never call this interface — they go through MessageProcessor.
type BankingServiceClient interface {
	CommitInterbankPayment(ctx context.Context, bankingTxID uint) error
	RollbackInterbankPayment(ctx context.Context, bankingTxID uint) error
}
