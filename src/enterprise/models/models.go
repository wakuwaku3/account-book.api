package models

import (
	"time"
)

type (
	// Account は アカウントです
	Account struct {
		Email          string `firestore:"-"`
		UserID         string `firestore:"userId"`
		HashedPassword string `firestore:"hashedPassword"`
		AccountToken   string `firestore:"accountToken"`
	}
	// PasswordResetToken は パスワードリセットのトークンです
	PasswordResetToken struct {
		PasswordResetToken string    `firestore:"-"`
		Expires            time.Time `firestore:"expires"`
		Email              string    `firestore:"email"`
	}
	// SignUpToken は ユーザー作成のトークンです
	SignUpToken struct {
		SignUpToken string    `firestore:"-"`
		Expires     time.Time `firestore:"expires"`
		Email       string    `firestore:"email"`
	}
	// User はユーザーです
	User struct {
		UserID       string    `firestore:"-"`
		UserName     string    `firestore:"userName"`
		Email        string    `firestore:"email"`
		Culture      string    `firestore:"culture"`
		UseStartDate time.Time `firestore:"useStartDate"`
	}
	// Plan は計画です
	Plan struct {
		PlanID     string     `firestore:"-"`
		PlanName   string     `firestore:"planName"`
		Interval   int        `firestore:"interval"`
		PlanAmount int        `firestore:"planAmount"`
		IsIncome   bool       `firestore:"isIncome"`
		Start      *time.Time `firestore:"start"`
		End        *time.Time `firestore:"end"`
		IsDeleted  bool       `firestore:"isDeleted"`
		CreatedAt  time.Time  `firestore:"createdAt"`
	}
	// Transaction は取引です
	Transaction struct {
		TransactionID string    `firestore:"-"`
		Amount        int       `firestore:"amount"`
		Category      int       `firestore:"category"`
		Date          time.Time `firestore:"date"`
		Notes         *string   `firestore:"notes"`
		DailyID       *string   `firestore:"dailyId"`
	}
	// Dashboard はダッシュボードです
	Dashboard struct {
		DashboardID         string    `firestore:"-"`
		Date                time.Time `firestore:"date"`
		Income              *int      `firestore:"income"`
		Expense             *int      `firestore:"expense"`
		CurrentBalance      *int      `firestore:"currentBalance"`
		Balance             *int      `firestore:"balance"`
		PreviousDashboardID *string   `firestore:"previousDashboardId"`
		PreviousBalance     *int      `firestore:"previousBalance"`
		State               string    `firestore:"state"`
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
		ActualID      string    `firestore:"-"`
		ActualAmount  int       `firestore:"actualAmount"`
		PlanID        string    `firestore:"planId"`
		PlanName      string    `firestore:"planName"`
		PlanAmount    int       `firestore:"planAmount"`
		IsIncome      bool      `firestore:"isIncome"`
		PlanCreatedAt time.Time `firestore:"planCreatedAt"`
	}
	// ActualKey はActualを特定するための要素です
	ActualKey struct {
		PlanID        string
		ActualID      *string
		DashboardID   *string
		SelectedMonth *time.Time
	}
)
