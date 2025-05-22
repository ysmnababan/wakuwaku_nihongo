package main

import (
	"wakuwaku_nihongo/internals/factory"
	"wakuwaku_nihongo/internals/pkg/database"
	"wakuwaku_nihongo/internals/utils/env"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

func init() {
	env := env.NewEnv()
	env.Load(`../../.env`)
}

func main() {
	database.Init("std")
	f := factory.NewFactory()
	db := f.Db

	g := gen.NewGenerator(gen.Config{
		OutPath: "../../internals/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		// WithUnitTest:     true,
		FieldNullable:    true,
		FieldCoverable:   true,
		FieldSignable:    true,
		FieldWithTypeTag: true,
	})
	g.UseDB(db)

	customers := g.GenerateModel("customers",
		gen.FieldNewTag("password", field.Tag{
			"json": "-",
		}))
	jlpt_books := g.GenerateModel("jlpt_books")
	questions := g.GenerateModel("questions")
	answers := g.GenerateModel("answers")

	quizzes := g.GenerateModel("quizzes",
		gen.FieldRelate(
			field.HasMany,
			"Questions",
			questions,
			&field.RelateConfig{
				RelateSlicePointer: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"quiz_id"},
					"references": []string{"quiz_id"},
				},
			},
		),
	)

	questions = g.GenerateModel("questions",
		gen.FieldRelate(
			field.BelongsTo,
			"Quiz",
			quizzes,
			&field.RelateConfig{
				RelatePointer: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"quiz_id"},
					"references": []string{"quiz_id"},
				},
			},
		),
		gen.FieldRelate(
			field.HasMany,
			"Answers",
			answers,
			&field.RelateConfig{
				RelateSlicePointer: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"question_id"},
					"references": []string{"question_id"},
				},
			},
		),
	)

	answers = g.GenerateModel("answers",
		gen.FieldRelate(
			field.BelongsTo,
			"Question",
			questions,
			&field.RelateConfig{
				RelatePointer: true,
				GORMTag: field.GormTag{
					"foreignKey": []string{"question_id"},
					"references": []string{"question_id"},
				},
			},
		))

	g.ApplyBasic(
		customers,
		jlpt_books,
		quizzes,
		questions,
		answers,
	)
	g.Execute()
}
