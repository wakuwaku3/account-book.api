package models

import (
	"time"
)

type (
	// Account は アカウントです
	Account struct {
		Email          string `db:"Email" firestore:"-"`
		UserID         string `db:"UserID" firestore:"user-id"`
		HashedPassword string `db:"UserID" firestore:"hashed-password"`
	}
	// User はユーザーです
	User struct {
		UserID   string `db:"UserID" firestore:"-"`
		UserName string `db:"UserName" firestore:"user-name"`
		Email    string `db:"Email" firestore:"email"`
	}
	// Plan は計画です
	Plan struct {
		PlanID     string     `db:"PlanID" firestore:"-"`
		PlanName   string     `db:"PlanName" firestore:"plan-name"`
		Interval   int        `db:"Interval" firestore:"interval"`
		PlanAmount int        `db:"PlanAmount" firestore:"plan-amount"`
		IsIncome   bool       `db:"IsIncome" firestore:"is-income"`
		Start      *time.Time `db:"Start" firestore:"start"`
		End        *time.Time `db:"End" firestore:"end"`
		IsDeleted  bool       `db:"IsDeleted" firestore:"is-deleted"`
	}
	// Transaction は取引です
	Transaction struct {
		TransactionID string    `db:"TransactionID" firestore:"-"`
		Amount        int       `db:"Amount" firestore:"amount"`
		Category      int       `db:"Category" firestore:"category"`
		Date          time.Time `db:"Date" firestore:"date"`
		Notes         *string   `db:"Notes" firestore:"notes"`
		DailyID       *string   `db:"DailyID" firestore:"daily-id"`
	}
	// Dashboard はダッシュボードです
	Dashboard struct {
		DashboardID         string    `db:"DashboardID" firestore:"-"`
		Date                time.Time `db:"Date" firestore:"date"`
		Income              *int      `db:"Income" firestore:"income"`
		Expense             *int      `db:"Expense" firestore:"expense"`
		CurrentBalance      *int      `db:"CurrentBalance" firestore:"current-balance"`
		Balance             *int      `db:"Balance" firestore:"balance"`
		PreviousDashboardID *string   `db:"PreviousDashboardID" firestore:"previous-dashboard-id"`
		PreviousBalance     *int      `db:"PreviousBalance" firestore:"previous-balance"`
		Daily               []Daily   `firestore:"-"`
		Actual              []Actual  `firestore:"-"`
	}
	// Daily は日毎のデータです
	Daily struct {
		DailyID string    `db:"DailyID" firestore:"-"`
		Date    time.Time `db:"Date" firestore:"date"`
		Income  int       `db:"Income" firestore:"income"`
		Expense int       `db:"Expense" firestore:"expense"`
	}
	// Actual は実費のデータです
	Actual struct {
		ActualID     string    `db:"ActualID" firestore:"-"`
		Date         time.Time `db:"Date" firestore:"-"`
		ActualAmount int       `db:"ActualAmount" firestore:"actual-amount"`
		PlanID       string    `db:"PlanID" firestore:"plan-id"`
		PlanName     string    `db:"PlanName" firestore:"plan-name"`
		PlanAmount   int       `db:"PlanAmount" firestore:"plan-amount"`
		IsIncome     bool      `db:"IsIncome" firestore:"is-income"`
	}
)
