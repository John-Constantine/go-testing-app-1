package service

import (
	"encoding/json"
	"flow/cache"
	"flow/db/repository"
	"flow/logger"
	"flow/model"
	"flow/model/response_dto"
	"flow/utility"
	"strings"
	"time"
)

type JourneyService struct {
	JourneyServiceUtil *JourneyServiceUtil
	RequestValidator   *utility.RequestValidator
}

func NewJourneyService(journeyServiceUtil *JourneyServiceUtil, validator *utility.RequestValidator) *JourneyService {
	service := &JourneyService{
		JourneyServiceUtil: journeyServiceUtil,
		RequestValidator:   validator,
	}
	return service
}

const ttl int = 24 * 60 * 60 * 5

func (u JourneyService) GetJourneys(merchantId string, tenantId string, channelId string) response_dto.JourneyResponsesDto {
	methodName := "GetJourneys"
	logger.SugarLogger.Info(methodName, "Recieved request to get all the flows associated with merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisClient := cache.GetRedisClient()
	var journeyResponsesDto response_dto.JourneyResponsesDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponsesDto
	}
	logger.SugarLogger.Info(methodName, "Fetching flows from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
	redisKey := u.RequestValidator.GenerateRedisKey(merchantId, tenantId, channelId)
	cachedFlow, err := redisClient.Get(redisKey).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId, " with error: ", err)
	}
	if cachedFlow == "" {
		logger.SugarLogger.Info(methodName, "No flows exist in redis cache for merchant: ", merchantId, " tenantId: ", tenantId, " channelId: ", channelId)
		flowContext := model.FlowContext{
			MerchantId: merchantId,
			TenantId:   tenantId,
			ChannelId:  channelId}
		flows := u.JourneyServiceUtil.FetchAllJourneysFromDB(flowContext)
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap := u.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flows...)
		journeyResponse := u.JourneyServiceUtil.ConstructJourneysResponse(flows, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)

		//Do not set redis key when there is no entry for given flowContext.
		if len(journeyResponse.JourneyResponses) == 0 {
			return journeyResponse
		}

		response, err := json.Marshal(journeyResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " couldn't update redis as failed to marshal response with err: ", err)
			return journeyResponse
		}

		logger.SugarLogger.Info(methodName, " Adding redis key: ", redisKey)
		setStatus := redisClient.Set(redisKey, response, time.Duration(ttl))
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", redisKey)
		return journeyResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponsesDto)
	return journeyResponsesDto
}

func (f JourneyService) GetModuleByModuleID(moduleID string) response_dto.ModuleVersionResponseDto{
	moduleVersion := repository.NewModuleRepository().FetchModuleVersion(moduleID)
	sectionVersionMap, fieldVersionMap := f.JourneyServiceUtil.GetSectionAndFieldVersionNumberList(moduleID)
	return f.JourneyServiceUtil.GetModuleVersionResponseDto(moduleVersion, sectionVersionMap, fieldVersionMap)
}

func (f JourneyService) GetJourneyById(journeyExternalId string) response_dto.JourneyResponseDto {
	methodName := "GetJourneyById"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", journeyExternalId)
	redisClient := cache.GetRedisClient()
	var journeyResponseDto response_dto.JourneyResponseDto
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponseDto
	}
	logger.SugarLogger.Info("Fetching the flow data from redis for journeyExternalId ", journeyExternalId)
	key := "journeyId:" + journeyExternalId + ":nested"
	cachedFlow, err := redisClient.Get(key).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for journeyExternalId: ", journeyExternalId, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.JourneyServiceUtil.FetchJourneyByIdFromDB(journeyExternalId)
		if len(flow.Name) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", journeyExternalId)
			return journeyResponseDto
		}
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap := f.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flow)
		flowsResponse := f.JourneyServiceUtil.ConstructFlowResponseWithModuleFieldSection(flow, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)

		response, err := json.Marshal(flowsResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return flowsResponse
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", journeyExternalId)
		setStatus := redisClient.Set(key, response, time.Duration(ttl))
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", journeyExternalId)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponseDto)
	return journeyResponseDto
}

func (f JourneyService) GetJourneyDetailsAsList(journeyExternalId string) response_dto.JourneyResponseDtoList {
	methodName := "GetJourneyById"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", journeyExternalId)
	redisClient := cache.GetRedisClient()
	var journeyResponseDto response_dto.JourneyResponseDtoList
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponseDto
	}
	logger.SugarLogger.Info("Fetching the flow data from redis for journeyExternalId ", journeyExternalId)
	//redisClient.FlushAll()
	key := "journeyId:" + journeyExternalId
	cachedFlow, err := redisClient.Get(key).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for journeyExternalId: ", journeyExternalId, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.JourneyServiceUtil.FetchJourneyByIdFromDB(journeyExternalId)
		if len(flow.Name) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", journeyExternalId)
			return journeyResponseDto
		}
		moduleVersionsMap, sectionVersionsMap, fieldVersionsMap := f.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(flow)
		flowsResponse := f.JourneyServiceUtil.ConstructFlowResponseAsList(flow, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)

		response, err := json.Marshal(flowsResponse)
		if err != nil {
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return flowsResponse
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", journeyExternalId)
		setStatus := redisClient.Set(key,response,time.Duration(ttl))
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", journeyExternalId)
		return flowsResponse
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponseDto)
	return journeyResponseDto
}

func (f JourneyService) GetJourneyDetailsListForJourneyIds(journeyExternalIds []string) []response_dto.JourneyResponseDtoList {
	methodName := "GetJourneyDetailsListForJourneyIds"
	logger.SugarLogger.Info(methodName, "Recieved request to get flow id ", journeyExternalIds)
	redisClient := cache.GetRedisClient()
	var journeyResponseDtoList []response_dto.JourneyResponseDtoList
	if redisClient == nil {
		logger.SugarLogger.Info(methodName, "Failed to connect with redis client. ")
		return journeyResponseDtoList
	}
	journeyIds := strings.Join(journeyExternalIds,":")
	logger.SugarLogger.Info("Fetching the flow data from redis for journeyExternalId ", journeyIds)
	key := "journeyId:" + journeyIds
	cachedFlow, err := redisClient.Get(key).Result()
	if err != nil {
		logger.SugarLogger.Info(methodName, "Failed to fetch flow from redis cache for journeyExternalId: ", journeyIds, " with error: ", err)
	}
	if len(cachedFlow) == 0 {
		flow := f.JourneyServiceUtil.FetchJourneyByJourneyIdListFromDB(journeyExternalIds)
		if len(flow) <= 0 {
			logger.SugarLogger.Error(methodName, " Invalid flow id passed : ", journeyIds)
			return journeyResponseDtoList
		}
		for _, journey := range flow {
			moduleVersionsMap, sectionVersionsMap, fieldVersionsMap := f.JourneyServiceUtil.GetModuleSectionAndFieldVersionsAndActiveVersionNumberList(journey)
			flowsResponse := f.JourneyServiceUtil.ConstructFlowResponseAsList(journey, moduleVersionsMap, sectionVersionsMap, fieldVersionsMap)
			journeyResponseDtoList = append(journeyResponseDtoList, flowsResponse)
		}


		response, err := json.Marshal(journeyResponseDtoList)
		if err != nil {
			logger.SugarLogger.Error(methodName, " failed to marshal response with err: will not be able to update redis", err)
			return journeyResponseDtoList
		}
		logger.SugarLogger.Info(methodName, " Adding redis key: ", journeyIds)
		setStatus := redisClient.Set(key, response, time.Duration(ttl))
		logger.SugarLogger.Info(methodName, " Set redis key status: ", setStatus.Val(), " for key: ", journeyIds)
		return journeyResponseDtoList
	}
	logger.SugarLogger.Info(methodName, " UnMarshlling the cached flow response")
	json.Unmarshal([]byte(cachedFlow), &journeyResponseDtoList)
	return journeyResponseDtoList
}
