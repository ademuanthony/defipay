package app

import (
	"context"
	"deficonnect/defipayapi/web"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type CreateAgentInput struct {
	SlackUsername string `json:"slack_username"`
	Name          string `json:"name"`
	Balance       int64  `json:"balance"`
	Status        int    `json:"status"`
}

type AgentOutput struct {
	SlackUsername string `json:"slack_username"`
	Name          string `json:"name"`
	Balance       int64  `json:"balance"`
	Status        int    `json:"status"`
}

type UpdateAgentInput struct {
	Balance *int64 `json:"balance"`
	Status  *int64 `json:"status"`
}

type GetAgentsInput struct {
	Query  string `json:"slackUsername"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type GetAgentAssignmentsInput struct {
	SlackUsername string `json:"slackUsername"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
}

type AssignmentOutput struct {
	TransactionID string           `json:"transactionId"`
	Date          int64            `json:"date"`
	Amount        int64            `json:"amount"`
	Status        AssignmentStatus `json:"status"`
}

type AgentStatus int

var AgentStatuses = struct {
	InActive AgentStatus
	Active   AgentStatus
}{
	InActive: 0,
	Active:   1,
}

type AssignmentStatus int

var AssignmentStatuses = struct {
	Pending    AssignmentStatus
	Processing AssignmentStatus
	Completed  AssignmentStatus
}{
	Pending:    0,
	Processing: 1,
	Completed:  2,
}

func (m Module) CreateAgent(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input CreateAgentInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Status != int(AgentStatuses.Active) && input.Status != int(AgentStatuses.InActive) {
		return SendErrorfJSON("Invalid agent status")
	}

	err := m.db.CreateAgent(ctx, input)
	if err != nil {
		return m.handleError(err)
	}

	return SendJSON(input)
}

func (m Module) GetAgents(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	pagedReq := web.GetPaginationInfoSls(r)
	var input = GetAgentsInput{
		Offset: pagedReq.Offset,
		Limit:  pagedReq.Limit,
		Query:  r.QueryStringParameters["q"],
	}

	agents, totalCount, err := m.db.GetAgents(ctx, input)
	if err != nil {
		return m.handleError(err)
	}

	return SendPagedJSON(agents, totalCount)
}

func (m Module) GetAgent(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	slackUsername := r.PathParameters["slackUsername"]
	if slackUsername == "" {
		SendErrorfJSON("Username is required")
	}

	agent, err := m.db.GetAgent(ctx, slackUsername)
	if err != nil {
		return m.handleError(err)
	}

	return SendJSON(agent)
}

func (m Module) UpdateAgent(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input UpdateAgentInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Status != nil && (AgentStatuses.Active != AgentStatus(*input.Status) &&
		AgentStatuses.InActive != AgentStatus(*input.Status)) {
		return SendErrorfJSON("Invalid agent status")
	}

	slackUsername := r.PathParameters["slackUsername"]
	if slackUsername == "" {
		SendErrorfJSON("Username is required")
	}

	agent, err := m.db.UpdateAgent(ctx, slackUsername, input)
	if err != nil {
		return m.handleError(err)
	}

	return SendJSON(agent)
}
