package models

import (
	"time"
)

type (
	// Account は アカウントです
	Account struct {
		Email          string `firestore:"-"`
		UserID         string `firestore:"user-id"`
		HashedPassword string `firestore:"hashed-password"`
	}
	// PasswordResetToken は パスワードリセットのトークンです
	PasswordResetToken struct {
		PasswordResetToken string    `firestore:"-"`
		Expires            time.Time `firestore:"expires"`
	}
	// User はユーザーです
	User struct {
		UserID   string `firestore:"-"`
		UserName string `firestore:"user-name"`
		Email    string `firestore:"email"`
	}
	// Plan は計画です
	Plan struct {
		PlanID     string     `firestore:"-"`
		PlanName   string     `firestore:"plan-name"`
		Interval   int        `firestore:"interval"`
		PlanAmount int        `firestore:"plan-amount"`
		IsIncome   bool       `firestore:"is-income"`
		Start      *time.Time `firestore:"start"`
		End        *time.Time `firestore:"end"`
		IsDeleted  bool       `firestore:"is-deleted"`
	}
	// Transaction は取引です
	Transaction struct {
		TransactionID string    `firestore:"-"`
		Amount        int       `firestore:"amount"`
		Category      int       `firestore:"category"`
		Date          time.Time `firestore:"date"`
		Notes         *string   `firestore:"notes"`
		DailyID       *string   `firestore:"daily-id"`
	}
	// Dashboard はダッシュボードです
	Dashboard struct {
		DashboardID         string    `firestore:"-"`
		Date                time.Time `firestore:"date"`
		Income              *int      `firestore:"income"`
		Expense             *int      `firestore:"expense"`
		CurrentBalance      *int      `firestore:"current-balance"`
		Balance             *int      `firestore:"balance"`
		PreviousDashboardID *string   `firestore:"previous-dashboard-id"`
		PreviousBalance     *int      `firestore:"previous-balance"`
		Daily               []Daily   `firestore:"-"`
		Actual              []Actual  `firestore:"-"`
	}
	// Daily は日毎のデータです
	Daily struct {
		DailyID string    `firestore:"-"`
		Date    time.Time `firestore:"date"`
		Income  int       `firestore:"income"`
		Expense int       `firestore:"expense"`
	}
	// Actual は実費のデータです
	Actual struct {
		ActualID     string    `firestore:"-"`
		Date         time.Time `firestore:"-"`
		ActualAmount int       `firestore:"actual-amount"`
		PlanID       string    `firestore:"plan-id"`
		PlanName     string    `firestore:"plan-name"`
		PlanAmount   int       `firestore:"plan-amount"`
		IsIncome     bool      `firestore:"is-income"`
	}
)
