package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var transactionCompleted = "completed"
var transactionReassign = "reassign"

func (m Module) SlackCallback(ctx context.Context, r *events.APIGatewayProxyRequest) (Response, error) {
	var input = struct {
		Challenge string `json:"challenge"`
		Type      string `json:"type"`
		Text      string `json:"text"`
		User      string `json:"user"`
	}{}

	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Text == "challenge" {
		return SendJSON(input)
	}

	if input.Type == "message" {
		func() {
			payloads := strings.Split(input.Text, " ")
			if len(payloads) != 2 {
				m.sendSlackMessage(ctx, fmt.Sprintf("Invalid command sent by <@%s>", input.User))
				return
			}

			command := payloads[0]
			transactionID := payloads[1]

			transaction, err := m.db.Transaction(ctx, transactionID)
			if err != nil {
				m.sendSlackMessage(ctx, fmt.Sprintf("Invalid transaction ID in command from <@%s>", input.User))
				return
			}

			switch command {
			case transactionCompleted:
				if err := m.db.UpdateTransactionStatus(ctx, transaction.ID, TransactionStatuses.Completed); err != nil {
					m.sendSlackMessage(ctx, fmt.Sprintf("error in marking transaction as completed, ID: %s. Manual check required. Completed by <@%s>",
					 transactionID, input.User))
					 return
				}
				m.sendSlackMessage(ctx, fmt.Sprintf("Transaction with ID: %s marked as completed by <@%s>", transactionID, input.User))
			case transactionReassign:
				if err = m.assignTransactionToAgent(ctx, transaction, true); err != nil {
					m.sendSlackMessage(ctx, fmt.Sprintf("Cannot reassign transaction with ID: %s as requested by <@%s>", transactionID, input.User))
					return
				}
				m.sendSlackMessage(ctx, fmt.Sprintf("Transaction with ID %s reassigned by <@%s>", transactionID, input.User))
			}
		}()
	}

	return SendJSON("ok")
}
