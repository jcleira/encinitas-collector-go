package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

// transactionsRetriever  defines the methods needed to retrievettransactions.
type transactionsRetriever interface {
	GetTransactionDetailAggregated(context.Context) ([]aggregates.TransactionDetailAggregated, error)
}

// TransactionsRetrieverHandler defines the dependencies to retrieve transactions.
type TransactionsRetrieverHandler struct {
	transactionsRetriever
}

// NewTransactionsRetriever initializes a new TransactionsRetrieverHandler.
func NewTransactionsRetriever(
	transactionsRetriever transactionsRetriever) *TransactionsRetrieverHandler {
	return &TransactionsRetrieverHandler{
		transactionsRetriever: transactionsRetriever,
	}
}

// Handle is the handler function to retrieve transactions
func (ech *TransactionsRetrieverHandler) Handle(c *gin.Context) {
	transactions, err := ech.transactionsRetriever.GetTransactionDetailAggregated(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpTransactionsResponse := struct {
		Transactions []struct {
			ProgramAddress string  `json:"program_address"`
			Percentage     float64 `json:"percentage"`
		} `json:"transactions"`
	}{}

	for _, transaction := range transactions {
		httpTransactionsResponse.Transactions = append(
			httpTransactionsResponse.Transactions,
			struct {
				ProgramAddress string  `json:"program_address"`
				Percentage     float64 `json:"percentage"`
			}{
				ProgramAddress: transaction.ProgramAddress,
				Percentage:     transaction.Percentage,
			},
		)
	}

	c.JSON(http.StatusOK, httpTransactionsResponse)
}
