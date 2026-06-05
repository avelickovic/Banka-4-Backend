package grpc

import (
	"context"
	"fmt"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/interbank-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/interbank-service/internal/service"
)

type outboxBankingClient struct{}

// NewOutboxBankingClient returns a stub BankingServiceClient. The full
// implementation requires a proto RPC (CommitInterbankPaymentByTxID) that
// does not yet exist in banking-service. OTC transactions never call this —
// they go through MessageProcessor directly.
func NewOutboxBankingClient(_ *client.BankingServiceConn) service.BankingServiceClient {
	return &outboxBankingClient{}
}

func (c *outboxBankingClient) CommitInterbankPayment(_ context.Context, bankingTxID uint) error {
	return fmt.Errorf("CommitInterbankPayment(txID=%d): not implemented — banking proto RPC missing", bankingTxID)
}

func (c *outboxBankingClient) RollbackInterbankPayment(_ context.Context, bankingTxID uint) error {
	return fmt.Errorf("RollbackInterbankPayment(txID=%d): not implemented — banking proto RPC missing", bankingTxID)
}
