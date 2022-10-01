package postgres

import (
	"context"
	"database/sql"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres/models"
	"errors"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreateAgent(ctx context.Context, input app.CreateAgentInput) error {
	item := models.Agent{
		SlackUsername: input.SlackUsername,
		Name:          input.Name,
		Balance:       input.Balance,
		Status:        input.Status,
	}
	return item.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) UpdateAgent(ctx context.Context, agentID int, input app.UpdateAgentInput) (*app.AgentOutput, error) {
	var col = models.M{}
	if input.Balance != nil {
		col[models.AgentColumns.Balance] = input.Balance
	}

	if input.Status != nil {
		col[models.AgentColumns.Status] = input.Status
	}

	_, err := models.Agents(models.AgentWhere.ID.EQ(agentID)).UpdateAll(ctx, pg.Db, col)
	if err != nil {
		return nil, err
	}

	agent, err := models.Agents(models.AgentWhere.ID.EQ(agentID)).One(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	return agentFromModel(agent), err
}

func agentFromModel(model *models.Agent) *app.AgentOutput {
	return &app.AgentOutput{
		ID:            model.ID,
		SlackUsername: model.SlackUsername,
		Name:          model.Name,
		Balance:       model.Balance,
		Status:        model.Status,
	}
}

func (pg PgDb) AgentExists(ctx context.Context, slackUsername string) (bool, error) {
	return models.Agents(models.AgentWhere.SlackUsername.EQ(slackUsername)).Exists(ctx, pg.Db)
}

func (pg PgDb) GetAgent(ctx context.Context, slackUsername string) (*app.AgentOutput, error) {
	agent, err := models.Agents(models.AgentWhere.SlackUsername.EQ(slackUsername)).One(ctx, pg.Db)
	if err != nil {
		return nil, err
	}
	return agentFromModel(agent), nil
}

func (pg PgDb) GetAgents(ctx context.Context, input app.GetAgentsInput) ([]*app.AgentOutput, int64, error) {
	var query []qm.QueryMod
	if input.Query != "" {
		query = append(query, qm.Where(`slack_username like %1 or slack_username like $2 or slack_username like $3 or
		name like %4 or name like $5 or name like $6`,
			("%"+input.Query),
			("%"+input.Query+"%"),
			(input.Query+"%"),
			("%"+input.Query),
			("%"+input.Query+"%"),
			(input.Query+"%"),
		))
	}

	count, err := models.Agents(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Offset(input.Offset), qm.Limit(input.Limit))

	agents, err := models.Agents(query...).All(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	var result []*app.AgentOutput
	for _, agent := range agents {
		result = append(result, agentFromModel(agent))
	}

	return result, count, err
}

func (pg PgDb) NextAvailableAgent(ctx context.Context, transactionAmount int64) (*app.AgentOutput, error) {
	lastAssignment, err := models.TransactionAssignments(
		qm.Select(models.TransactionAssignmentColumns.AgentID),
		qm.OrderBy(models.TransactionAssignmentColumns.Date+" desc"),
	).One(ctx, pg.Db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	lastAgentID := 0
	if err == nil {
		lastAgentID = lastAssignment.AgentID
	}

	maxAgent, _ := models.Agents(
		qm.Select(models.AgentColumns.ID),
		qm.OrderBy(models.AgentColumns.ID+" desc"),
	).One(ctx, pg.Db)

	if maxAgent == nil || lastAgentID == maxAgent.ID {
		lastAgentID = 0
	}

	for id := lastAgentID + 1; id <= maxAgent.ID; id++ {
		agent, err := models.FindAgent(ctx, pg.Db, id)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, err
		}
		if agent.Balance < transactionAmount {
			continue
		}
		if agent.Status == int(app.AgentStatuses.InActive) {
			continue
		}

		return agentFromModel(agent), nil
	}

	return nil, errors.New("agent not found")
}

func (pg PgDb) AssignAgent(ctx context.Context, agentID int, transactionID string, amount int64) error {
	assignment := models.TransactionAssignment{
		ID:            1,
		AgentID:       agentID,
		TransactionID: transactionID,
		Amount:        amount,
		Date:          time.Now().Unix(),
		Status:        int(app.AssignmentStatuses.Pending),
	}

	agent, err := models.FindAgent(ctx, pg.Db, agentID)
	if err != nil {
		return errors.New("invalid agent ID")
	}

	err = assignment.Insert(ctx, pg.Db, boil.Infer())
	if err != nil {
		return nil
	}

	balance := agent.Balance - amount
	_, err = pg.UpdateAgent(ctx, agentID, app.UpdateAgentInput{
		Balance: &balance,
	})
	return err
}

func (pg PgDb) GetAssignedAgent(ctx context.Context, transactionID string) (*app.AgentOutput, error) {
	assignment, err := models.TransactionAssignments(
		qm.Load(models.TransactionAssignmentRels.Agent),
		models.TransactionAssignmentWhere.TransactionID.EQ(transactionID),
	).One(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	return agentFromModel(assignment.R.Agent), nil
}

func (pg PgDb) GetAgentAssignments(ctx context.Context, input app.GetAgentAssignmentsInput) ([]app.AssignmentOutput, error) {
	assignments, err := models.TransactionAssignments(
		qm.Load(models.TransactionAssignmentRels.Agent),
		models.TransactionAssignmentWhere.AgentID.EQ(input.AgentID),
	).All(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	var result []app.AssignmentOutput
	for _, item := range assignments {
		result = append(result, app.AssignmentOutput{
			Date:          item.Date,
			Amount:        item.Amount,
			TransactionID: item.TransactionID,
			Status:        app.AssignmentStatus(item.Status),
		})
	}

	return result, nil
}
