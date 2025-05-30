// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameQuiz = "quizzes"

// Quiz mapped from table <quizzes>
type Quiz struct {
	QuizID      string      `gorm:"column:quiz_id;type:uuid;primaryKey" json:"quiz_id"`
	CreatedAt   int64       `gorm:"column:created_at;type:bigint;not null" json:"created_at"`
	ModifiedAt  *int64      `gorm:"column:modified_at;type:bigint" json:"modified_at"`
	DeletedAt   *int64      `gorm:"column:deleted_at;type:bigint" json:"deleted_at"`
	CreatedBy   string      `gorm:"column:created_by;type:character varying;not null" json:"created_by"`
	ModifiedBy  *string     `gorm:"column:modified_by;type:character varying" json:"modified_by"`
	DeletedBy   *string     `gorm:"column:deleted_by;type:character varying" json:"deleted_by"`
	Title       string      `gorm:"column:title;type:character varying;not null" json:"title"`
	Description *string     `gorm:"column:description;type:character varying" json:"description"`
	Questions   []*Question `gorm:"foreignKey:quiz_id;references:quiz_id" json:"questions"`
}

// TableName Quiz's table name
func (*Quiz) TableName() string {
	return TableNameQuiz
}
