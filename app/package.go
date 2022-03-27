package app

import (
	"database/sql"
	"encoding/json"
	"merryworld/metatradas/postgres/models"
	"merryworld/metatradas/web"
	"net/http"
)

type UpdatePackageInput struct {
	ID                string `json:"id"`
	Name              string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Price             int64  `boil:"price" json:"price" toml:"price" yaml:"price"`
	MinReturnPerMonth int    `boil:"min_return_per_month" json:"min_return_per_month" toml:"min_return_per_month" yaml:"min_return_per_month"`
	MaxReturnPerMonth int    `boil:"max_return_per_month" json:"max_return_per_month" toml:"max_return_per_month" yaml:"max_return_per_month"`
	TradesPerDay      int    `boil:"trades_per_day" json:"trades_per_day" toml:"trades_per_day" yaml:"trades_per_day"`
	Accuracy          int    `boil:"accuracy" json:"accuracy" toml:"accuracy" yaml:"accuracy"`
}

type SubscriptionResponse struct {
	*models.Subscription
	Package models.Package `json:"package"`
}

type InvestInput struct {
	Amount int64 `json:"amount"`
}

func (m module) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		log.Critical("CreatePackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	if pkg.Name == "" || pkg.Accuracy == 0 || pkg.MinReturnPerMonth == 0 || pkg.MaxReturnPerMonth == 0 {
		web.SendErrorfJSON(w, "Required filedss not sent")
		return
	}

	if _, err := m.db.GetPackageByName(r.Context(), pkg.Name); err == nil {
		web.SendErrorfJSON(w, "Name not available")
		return
	}

	if err := m.db.CreatePackage(r.Context(), pkg); err != nil {
		log.Error("CreatePackage", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg UpdatePackageInput
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		log.Critical("UpdatePackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	oldPkg, err := m.db.GetPackage(r.Context(), pkg.ID)
	if err != nil {
		web.SendErrorfJSON(w, "Invalid package ID")
		return
	}

	if pkg.Name != "" {
		if p, err := m.db.GetPackageByName(r.Context(), pkg.Name); err == nil && p.ID != oldPkg.ID {
			web.SendErrorfJSON(w, "Name not available")
			return
		}
	}

	if err := m.db.PatchPackage(r.Context(), oldPkg.ID, pkg); err != nil {
		log.Error("UpdatePackage", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) GetPackage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	pkg, err := m.db.GetPackage(r.Context(), id)
	if err == sql.ErrNoRows {
		web.SendErrorfJSON(w, "Invalid package ID")
		return
	}

	if err != nil {
		log.Error("CreatePackage", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, pkg)
}

func (m module) GetPackages(w http.ResponseWriter, r *http.Request) {
	packages, err := m.db.GetPackages(r.Context())

	if err != nil {
		log.Error("CreatePackage", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, packages)
}

type buyPackageInput struct {
	ID string `json:"id"`
}

func (m module) BuyPackage(w http.ResponseWriter, r *http.Request) {
	var input buyPackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("BuyPackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}
	acc, err := m.currentAccount(r)
	if err != nil {
		log.Info("BuyPackage", "currentAccount", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	if sub, _ := m.db.ActiveSubscription(r.Context(), acc.ID); sub != nil {
		web.SendErrorfJSON(w, "You have an active subscription. Please use the upgrade function")
		return
	}

	pkg, err := m.db.GetPackage(r.Context(), input.ID)
	if err != nil {
		web.SendErrorfJSON(w, "Invalid package ID")
		return
	}

	if acc.Balance < pkg.Price {
		web.SendErrorfJSON(w, "Insufficient fund, please topup your account to continue")
		return
	}

	if err := m.db.CreateSubscription(r.Context(), acc.ID, pkg.ID, false); err != nil {
		log.Error("BuyPackage", "CreateSubscription", err)
		web.SendErrorfJSON(w, "Error in creating subscription")
		return
	}

	web.SendJSON(w, true)
}

type createSubscriptionC250 struct {
	PackageID string `json:"package_id"`
	Username  string `json:"username"`
}

func (m module) createSubscriptionC250(w http.ResponseWriter, r *http.Request) {
	var input createSubscriptionC250
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("BuyPackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	acc, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		web.SendErrorfJSON(w, "Account not found. Please activate your METATRADAS account")
	}

	if err := m.db.CreateSubscription(r.Context(), input.PackageID, acc.ID, true); err != nil {
		log.Error("createSubscriptionC250", "CreateSubscription", err)
		web.SendErrorfJSON(w, "Error in creating subscription")
		return
	}

	web.SendJSON(w, true)
}

func (m module) GetActiveSubscription(w http.ResponseWriter, r *http.Request) {
	sub, err := m.db.ActiveSubscription(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err == sql.ErrNoRows {
		web.SendErrorfJSON(w, "You do not have an active subscription")
		return
	}

	if err != nil {
		log.Critical("GetActiveSubscription", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	resp := SubscriptionResponse{
		Subscription: sub,
		Package:      *sub.R.Package,
	}

	web.SendJSON(w, resp)
}
