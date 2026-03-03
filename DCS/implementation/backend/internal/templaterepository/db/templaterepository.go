package db

import (
	approvaltask2 "digital-contracting-service/internal/templaterepository/datatype/approvaltask"
	"digital-contracting-service/internal/templaterepository/datatype/approvaltaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtask"
	"digital-contracting-service/internal/templaterepository/datatype/reviewtaskstate"
	"digital-contracting-service/internal/templaterepository/datatype/templaterepository"
	"digital-contracting-service/internal/templaterepository/datatype/templatestate"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReviewTaskRepo interface {
	Create(tx *sqlx.Tx, data reviewtask.TaskData) (*time.Time, error)
	IsValidReviewer(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string) (bool, error)
	ReopenTasks(tx *sqlx.Tx, did string, documentNumber int, version int) error
	ReadAll(tx *sqlx.Tx, did string) ([]reviewtask.TaskData, error)
	ReadAllByID(tx *sqlx.Tx, did string, documentNumber int, version int) ([]reviewtask.TaskData, error)
	ReadAllByReviewer(tx *sqlx.Tx, reviewer string) ([]reviewtask.TaskData, error)
	Update(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) error
	AnyTasksInState(tx *sqlx.Tx, did string, documentNumber int, version int, states ...reviewtaskstate.ReviewTaskState) (bool, error)
	TaskExistsInState(tx *sqlx.Tx, did string, documentNumber int, version int, reviewer string, state reviewtaskstate.ReviewTaskState) (bool, error)
	TaskExist(tx *sqlx.Tx, did string, documentNumber int, version int) (bool, error)
	Delete(tx *sqlx.Tx, did string, documentNumber int, version int) error
}

type ApprovalTaskRepo interface {
	Create(tx *sqlx.Tx, data approvaltask2.TaskData) (*time.Time, error)
	ReopenTasks(tx *sqlx.Tx, did string, documentNumber int, version int) error
	ReadAll(dtx *sqlx.Tx, id string) ([]approvaltask2.TaskData, error)
	ReadAllByApprover(tx *sqlx.Tx, approver string) ([]approvaltask2.TaskData, error)
	Update(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state approvaltaskstate.ApprovalTaskState) error
	IsValidApprover(tx *sqlx.Tx, did string, documentNumber int, version int, approver string) (bool, error)
	TaskExistsInState(tx *sqlx.Tx, did string, documentNumber int, version int, approver string, state approvaltaskstate.ApprovalTaskState) (bool, error)
	TaskExists(tx *sqlx.Tx, did string, documentNumber int, version int) (bool, error)
	Delete(tx *sqlx.Tx, did string, documentNumber int, version int) error
}

type TemplateRepository interface {
	Create(tx *sqlx.Tx, data templaterepository.ContractTemplate) (*time.Time, error)
	ReadDataByID(tx *sqlx.Tx, did string, documentNumber int, version int) (*templaterepository.ContractTemplate, error)
	ReadAllMetaData(tx *sqlx.Tx) ([]templaterepository.MetaData, error)
	ReadAllMetaDataByFilter(tx *sqlx.Tx, values templaterepository.SearchValues) ([]templaterepository.MetaData, error)
	ReadProcessData(tx *sqlx.Tx, did string, documentNumber int, version int) (*templaterepository.ProcessData, error)
	UpdateState(tx *sqlx.Tx, did string, documentNumber int, version int, state templatestate.TemplateState) error
	Update(tx *sqlx.Tx, data templaterepository.UpdateData) error
}
