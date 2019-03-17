package services

import (
	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	dashboard struct {
		repos domains.DashboardRepository
	}
	// Dashboard is DashboardService
	Dashboard interface {
	}
)

// NewDashboard is create instance
func NewDashboard(repos domains.DashboardRepository) Dashboard {
	return &dashboard{repos}
}
