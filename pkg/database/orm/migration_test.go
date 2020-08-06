package orm_test

import (
	"github.com/cryptogarageinc/server-common-go/pkg/database/orm"
	"github.com/cryptogarageinc/server-common-go/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMigrationModel struct {
	Name string
}

func TestMigrationInitialize_IsInitialized(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormConfig := &orm.Config{}
	test.InitializeConfig(ormConfig)
	ormInstance := orm.NewORM(ormConfig, test.NewLogger())
	ormInstance.Initialize()
	migrator := orm.NewMigrator(ormInstance, &TestMigrationModel{})

	// Act
	migrator.Initialize()

	// Assert
	assert.True(migrator.IsInitialized())
}

func TestMigrationInitialize_HasCorrectTable(t *testing.T) {
	// Arrange
	assert := assert.New(t)
	ormConfig := &orm.Config{}
	test.InitializeConfig(ormConfig)
	ormInstance := orm.NewORM(ormConfig, test.NewLogger())
	ormInstance.Initialize()
	migrator := orm.NewMigrator(ormInstance, &TestMigrationModel{})
	migrator.Initialize()

	// Act
	var result []TestMigrationModel
	err := ormInstance.GetDB().Find(&result).Error

	// Assert
	assert.NoError(err)
}
