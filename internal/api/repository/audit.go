package repository

import (
	"github.com/kirananto/review-system/internal/db"
	models "github.com/kirananto/review-system/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	CreateAuditLog(auditLog *models.AuditLog) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(dataSource *db.DataSource) AuditLogRepository {
	return &auditLogRepository{
		db: dataSource.Db,
	}
}

func (r *auditLogRepository) CreateAuditLog(auditLog *models.AuditLog) error {
	return r.db.Create(auditLog).Error
}
