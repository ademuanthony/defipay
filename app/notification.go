package app

import (
	"merryworld/metatradas/web"
	"net/http"
)

const (
	NOTIFICATION_STATUS_NEW  = 0
	NOTIFICATION_STATUS_READ = 1
)

func (m module) getUnReadNotificationCount(w http.ResponseWriter, r *http.Request) {
	count, err := m.db.UnReadNotificationCount(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Error("getNotificationCount", "UnReadNotificationCount", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
	}
	web.SendJSON(w, count)
}

func (m module) getNotifications(w http.ResponseWriter, r *http.Request) {
	pagedReq := web.GetPanitionInfo(r)
	notification, count, err := m.db.GetNotifications(r.Context(), m.server.GetUserIDTokenCtx(r), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		m.sendSomethingWentWrong(w, err)
		return
	}
	web.SendPagedJSON(w, notification, count)
}

func (m module) getNotification(w http.ResponseWriter, r *http.Request) {
	notification, err := m.db.GetNotification(r.Context(), r.FormValue("id"))
	if err != nil {
		m.sendSomethingWentWrong(w, err)
		return
	}
	if notification.AccountID != m.server.GetUserIDTokenCtx(r) {
		web.SendErrorfJSON(w, "Access denied")
		return
	}
	web.SendJSON(w, notification)
}

func (m module) sendSomethingWentWrong(w http.ResponseWriter, err error) {
	log.Error("getNotificationCount", "UnReadNotificationCount", err)
	web.SendErrorfJSON(w, "Something went wrong. Please try again later")
}
