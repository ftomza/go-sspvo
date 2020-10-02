/*
 * Copyright © 2020-present Artem V. Zaborskiy <ftomza@yandex.ru>. All rights reserved.
 *
 * This source code is licensed under the Apache 2.0 license found
 * in the LICENSE file in the root directory of this source tree.
 */

package message

import (
	"github.com/ftomza/go-sspvo"
	"github.com/ftomza/go-sspvo/response"
)

type CLS string

const (
	CLSDirections                CLS = "Directions"
	CLSCampaignType              CLS = "CampaignType"
	CLSCampaignStatus            CLS = "CampaignStatus"
	CLSBenefit                   CLS = "Benefit"
	CLSEducationForm             CLS = "EducationForm"
	CLSEducationLevel            CLS = "EducationLevel"
	CLSEducationSource           CLS = "EducationSource"
	CLSEntranceTestType          CLS = "EntranceTestType"
	CLSLevelBudget               CLS = "LevelBudget"
	CLSOlympicDiplomaType        CLS = "OlympicDiplomaType"
	CLSOlympicLevel              CLS = "OlympicLevel"
	CLSSubject                   CLS = "Subject"
	CLSEduLevelsCampaignTypes    CLS = "EduLevelsCampaignTypes"
	CLSAchievementCategory       CLS = "AchievementCategory"
	CLSApplicationStatuses       CLS = "ApplicationStatuses"
	CLSCompatriotCategories      CLS = "CompatriotCategories"
	CLSCompositionThemes         CLS = "CompositionThemes"
	CLSDisabilityTypes           CLS = "DisabilityTypes"
	CLSDocumentCategories        CLS = "DocumentCategories"
	CLSDocumentTypes             CLS = "DocumentTypes"
	CLSEntranceTestDocumentTypes CLS = "EntranceTestDocumentTypes"
	CLSEntranceTestResultSources CLS = "EntranceTestResultSources"
	CLSGenders                   CLS = "Genders"
	CLSMinScoreSubjects          CLS = "MinScoreSubjects"
	CLSOkcms                     CLS = "Okcms"
	CLSOktmos                    CLS = "Oktmos"
	CLSOlympicMinEge             CLS = "OlympicMinEge"
	CLSOrderAdmissionStatuses    CLS = "OrderAdmissionStatuses"
	CLSOrderAdmissionTypes       CLS = "OrderAdmissionTypes"
	CLSOrphanCategories          CLS = "OrphanCategories"
	CLSParentsLostCategories     CLS = "ParentsLostCategories"
	CLSRadiationWorkCategories   CLS = "RadiationWorkCategories"
	CLSRegions                   CLS = "Regions"
	CLSReturnTypes               CLS = "ReturnTypes"
	CLSVeteranCategories         CLS = "VeteranCategories"
	CLSViolationTypes            CLS = "ViolationTypes"
	CLSOlympicsProfiles          CLS = "OlympicsProfiles"
	CLSOlyProfiles               CLS = "OlyProfiles"
	CLSOlympics                  CLS = "Olympics"
	CLSAppealStatuses            CLS = "AppealStatuses"
	CLSMilitaryCategories        CLS = "MilitaryCategories"
)

var AllCLS = []CLS{
	CLSDirections,
	CLSCampaignType,
	CLSCampaignStatus,
	CLSBenefit,
	CLSEducationForm,
	CLSEducationLevel,
	CLSEducationSource,
	CLSEntranceTestType,
	CLSLevelBudget,
	CLSOlympicDiplomaType,
	CLSOlympicLevel,
	CLSSubject,
	CLSEduLevelsCampaignTypes,
	CLSAchievementCategory,
	CLSApplicationStatuses,
	CLSCompatriotCategories,
	CLSCompositionThemes,
	CLSDisabilityTypes,
	CLSDocumentCategories,
	CLSDocumentTypes,
	CLSEntranceTestDocumentTypes,
	CLSEntranceTestResultSources,
	CLSGenders,
	CLSMinScoreSubjects,
	CLSOkcms,
	CLSOktmos,
	CLSOlympicMinEge,
	CLSOrderAdmissionStatuses,
	CLSOrderAdmissionTypes,
	CLSOrphanCategories,
	CLSParentsLostCategories,
	CLSRadiationWorkCategories,
	CLSRegions,
	CLSReturnTypes,
	CLSVeteranCategories,
	CLSViolationTypes,
	CLSOlympicsProfiles,
	CLSOlyProfiles,
	CLSOlympics,
	CLSAppealStatuses,
	CLSMilitaryCategories,
}

func (e CLS) IsValid() bool {
	switch e {
	case CLSDirections,
		CLSCampaignType,
		CLSCampaignStatus,
		CLSBenefit,
		CLSEducationForm,
		CLSEducationLevel,
		CLSEducationSource,
		CLSEntranceTestType,
		CLSLevelBudget,
		CLSOlympicDiplomaType,
		CLSOlympicLevel,
		CLSSubject,
		CLSEduLevelsCampaignTypes,
		CLSAchievementCategory,
		CLSApplicationStatuses,
		CLSCompatriotCategories,
		CLSCompositionThemes,
		CLSDisabilityTypes,
		CLSDocumentCategories,
		CLSDocumentTypes,
		CLSEntranceTestDocumentTypes,
		CLSEntranceTestResultSources,
		CLSGenders,
		CLSMinScoreSubjects,
		CLSOkcms,
		CLSOktmos,
		CLSOlympicMinEge,
		CLSOrderAdmissionStatuses,
		CLSOrderAdmissionTypes,
		CLSOrphanCategories,
		CLSParentsLostCategories,
		CLSRadiationWorkCategories,
		CLSRegions,
		CLSReturnTypes,
		CLSVeteranCategories,
		CLSViolationTypes,
		CLSOlympicsProfiles,
		CLSOlyProfiles,
		CLSOlympics,
		CLSAppealStatuses,
		CLSMilitaryCategories:
		return true
	}
	return false
}

func (e CLS) String() string {
	return string(e)
}

type CLSMessage struct {
	Message
}

func NewCLSMessage(cls CLS) *CLSMessage {
	msg := &CLSMessage{}
	msg.Init()
	msg.UpdateJWTFields(setCLS(cls))

	return msg
}

func (m *CLSMessage) PathMethod() string {
	return pathMethodCLS
}

func (m *CLSMessage) Response() sspvo.Response {
	return response.NewResponse()
}
