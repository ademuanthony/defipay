package app

import (
	"context"
	"deficonnect/defipayapi/web"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

const (
	NOTIFICATION_STATUS_NEW  = 0
	NOTIFICATION_STATUS_READ = 1

	NOTIFICATION_TYPE_TOPBAR    = 0
	NOTIFICATION_TYPE_DASHBOARD = 1
)

type sendNotificationInput struct {
	Titile     string `json:"title"`
	Content    string `json:"content"`
	ActionLink string `json:"action_link"`
	ActionText string `json:"action_text"`
	Type       int    `josn:"type"`
}

func (m Module) getUnReadNotificationCount(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	notificationType, _ := strconv.ParseInt(r.QueryStringParameters["type"], 10, 64)
	count, err := m.db.UnReadNotificationCount(ctx, m.server.GetUserIDTokenCtxSls(r), int(notificationType))
	if err != nil {
		log.Error("getNotificationCount", "UnReadNotificationCount", err)
		return SendErrorfJSON("Something went wrong. Please try again later")
	}
	return SendJSON(count)
}

func (m Module) getNewNotifications(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	pagedReq := web.GetPaginationInfoSls(r)
	notificationType, _ := strconv.ParseInt(r.QueryStringParameters["type"], 10, 64)
	notification, count, err := m.db.GetNewNotifications(ctx, m.server.GetUserIDTokenCtxSls(r), int(notificationType), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		return m.sendSomethingWentWrong("GetNewNotifications", err)
	}
	return SendPagedJSON(notification, count)
}

func (m Module) getNotifications(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	pagedReq := web.GetPaginationInfoSls(r)
	notificationType, _ := strconv.ParseInt(r.QueryStringParameters["type"], 10, 64)
	notification, count, err := m.db.GetNotifications(ctx, m.server.GetUserIDTokenCtxSls(r), int(notificationType), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		return m.sendSomethingWentWrong("GetNotifications", err)
	}
	return SendPagedJSON(notification, count)
}

func (m Module) getNotification(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	notification, err := m.db.GetNotification(ctx, r.QueryStringParameters["id"])
	if err != nil {
		return m.sendSomethingWentWrong("GetNotification", err)
	}
	if notification.AccountID != m.server.GetUserIDTokenCtxSls(r) {
		return SendErrorfJSON("Access denied")
	}
	return SendJSON(notification)
}

func (m Module) sendSomethingWentWrong(fn string, err error) (Response, error) {
	log.Error(fn, err)
	return m.handleError(err)
}
