package service

import (
	"banking-service/internal/model"
	"banking-service/internal/repository"
	"common/pkg/errors"
	"context"
	
	"gorm.io/gorm"
)

var BankAccounts = map[model.CurrencyCode]string{
	model.RSD: "999-0000000000000-00",
	model.EUR: "999-0000000000001-00",
	model.USD: "999-0000000000002-00",
	model.CHF: "999-0000000000003-00",
	model.GBP: "999-0000000000004-00",
	model.JPY: "999-0000000000005-00",
	model.CAD: "999-0000000000006-00",
	model.AUD: "999-0000000000007-00",
}

type TransactionProcessor struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	exchangeService ExchangeService
	db              *gorm.DB
}

func NewTransactionProcessor(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, exchangeService ExchangeService, db *gorm.DB) *TransactionProcessor {
	return &TransactionProcessor{accountRepo: accountRepo, transactionRepo: transactionRepo, exchangeService: exchangeService, db: db}
}

func (tp *TransactionProcessor) Process(ctx context.Context, transactionID uint) error {
	return tp.db.Transaction(func(tx *gorm.DB) error {
		transaction, err := tp.transactionRepo.GetByIDWithTx(ctx, tx, transactionID)
		if err != nil {
				return errors.InternalErr(err)
		}

		payer, err := tp.accountRepo.GetByAccountNumberWithTx(ctx, tx, transaction.PayerAccountNumber)
		if err != nil {
				return errors.InternalErr(err)
		}

		recipient, err := tp.accountRepo.GetByAccountNumberWithTx(ctx, tx, transaction.RecipientAccountNumber)
		if err != nil {
				return errors.InternalErr(err)
		}

		banksAccountTo, err := tp.accountRepo.GetByAccountNumberWithTx(ctx, tx, BankAccounts[transaction.StartCurrencyCode])
		if err != nil {
				return errors.InternalErr(err)
		}

		banksAccountFrom, err := tp.accountRepo.GetByAccountNumberWithTx(ctx, tx, BankAccounts[transaction.EndCurrencyCode])
		if err != nil {
				return errors.InternalErr(err)
		}

		// Check funds
		if payer.AvailableBalance < transaction.StartAmount {
				return errors.BadRequestErr("insufficient payer funds")
		}
		if banksAccountFrom.AvailableBalance < transaction.EndAmount {
				return errors.BadRequestErr("insufficient banks funds")
		}

		model.UpdateBalances(payer, -transaction.StartAmount)
		model.UpdateBalances(banksAccountTo, transaction.StartAmount)
		model.UpdateBalances(banksAccountFrom, -transaction.EndAmount)
		model.UpdateBalances(recipient, transaction.EndAmount)

		if err := tp.accountRepo.UpdateWithTx(ctx, tx, payer); err != nil {
				return errors.InternalErr(err)
		}
		if err := tp.accountRepo.UpdateWithTx(ctx, tx, recipient); err != nil {
				return errors.InternalErr(err)
		}
		if err := tp.accountRepo.UpdateWithTx(ctx, tx, banksAccountTo); err != nil {
				return errors.InternalErr(err)
		}
		if err := tp.accountRepo.UpdateWithTx(ctx, tx, banksAccountFrom); err != nil {
				return errors.InternalErr(err)
		}

		transaction.Status = model.TransactionCompleted
		return tp.transactionRepo.UpdateWithTx(ctx, tx, transaction)
	})
}
