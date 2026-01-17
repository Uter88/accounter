package task

import (
	"accounter/internal/domain/common"
)

var (
	ErrTaskNotFound = &common.DomainError{Code: "err_task_not_found", Message: "task not found", StatusCode: 404}
	ErrTaskDelete   = &common.DomainError{Code: "err_delete_task", Message: "error delete task", StatusCode: 400}
	ErrTaskCreate   = &common.DomainError{Code: "err_create_task", Message: "error create task", StatusCode: 400}
	ErrTaskUpdate   = &common.DomainError{Code: "err_update_task", Message: "error update task", StatusCode: 400}
	ErrLogTask      = &common.DomainError{Code: "err_log_task", Message: "error log task", StatusCode: 400}
)
