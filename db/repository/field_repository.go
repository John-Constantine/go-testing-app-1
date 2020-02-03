package repository

import (
	"flow/db"
	"flow/logger"
	"flow/model"
	"flow/utility"
)

type FieldRepository struct {
	MapUtil utility.MapUtil
	DBService db.DBService
}

func (f FieldRepository) FetchFieldFromFieldVersion(completeFieldVersionNumberList map[int]bool) []model.FieldVersion {
	methodName := "FetchFieldFromFieldVersion"
	logger.SugarLogger.Info(methodName, " Fetching the field data with join on field versions")
	var fieldVersions []model.FieldVersion
	dbConnection := f.DBService.GetDB()
	if dbConnection == nil {
		return fieldVersions
	}
	dbConnection.Joins("JOIN field ON field.id = field_version.field_id ").Where("field_version.id in (?) ", f.MapUtil.GetKeyListFromKeyValueMap(completeFieldVersionNumberList)).Find(&fieldVersions)

	return fieldVersions
}

func (f FieldRepository) FindByExternalId(flowExternalId string) model.Flow {
	var flow model.Flow
	dbConnection := f.DBService.GetDB()
	dbConnection.Where(" external_id = ? ", flowExternalId).Find(&flow)
	return flow
}
