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

type Action string

const (
	ActionAdd    Action = "Add"
	ActionRemove Action = "Remove"
	ActionEdit   Action = "Edit"
	ActionGet    Action = "Get"

	actionGetMessage     Action = "GetMessage"
	actionMessageConfirm Action = "MessageConfirm"
)

var AllAction = []Action{
	ActionAdd,
	ActionRemove,
	ActionEdit,
	ActionGet,
	actionGetMessage,
	actionMessageConfirm,
}

func (e Action) IsValid() bool {
	switch e {
	case ActionAdd,
		ActionRemove,
		ActionEdit,
		ActionGet,
		actionGetMessage,
		actionMessageConfirm:
		return true
	}
	return false
}

func (e Action) String() string {
	return string(e)
}

type Datatype string

const (
	DatatypeSubdivisionOrg                       Datatype = "subdivision_org"
	DatatypeCampaign                             Datatype = "campaign"
	DatatypeAchievements                         Datatype = "achievements"
	DatatypeAdmissionVolume                      Datatype = "admission_volume"
	DatatypeDistributedAdmissionVolume           Datatype = "distributed_admission_volume"
	DatatypeCompetitiveGroups                    Datatype = "competitive_groups"
	DatatypeCompetitiveGroupPrograms             Datatype = "competitive_group_programs"
	DatatypeCompetitiveBenefits                  Datatype = "competitive_benefits"
	DatatypeEntranceTests                        Datatype = "entrance_tests"
	DatatypeEntranceTestBenefits                 Datatype = "entrance_test_benefits"
	DatatypeEntrants                             Datatype = "entrants"
	DatatypeCompatriot                           Datatype = "compatriot"
	DatatypeComposition                          Datatype = "composition"
	DatatypeDisability                           Datatype = "disability"
	DatatypeEducations                           Datatype = "educations"
	DatatypeEge                                  Datatype = "ege"
	DatatypeIdentification                       Datatype = "identification"
	DatatypeMilitaries                           Datatype = "militaries"
	DatatypeOlympics                             Datatype = "olympics"
	DatatypeOrphans                              Datatype = "orphans"
	DatatypeOther                                Datatype = "other"
	DatatypeParentsLost                          Datatype = "parents_lost"
	DatatypeRadiationWork                        Datatype = "radiation_work"
	DatatypeVeteran                              Datatype = "veteran"
	DatatypeApplications                         Datatype = "applications"
	DatatypeEditApplicationStatus                Datatype = "edit_application_status"
	DatatypeEntranceTestAgreed                   Datatype = "entrance_test_agreed"
	DatatypeEntranceTestResult                   Datatype = "entrance_test_result"
	DatatypeOrderAdmission                       Datatype = "order_admission"
	DatatypeCompletitiveGroupsApplicationsRating Datatype = "completitive_groups_applications_rating"
	DatatypeAppAchievements                      Datatype = "app_achievements"
	DatatypeApplicationsRating                   Datatype = "applications_rating"
	DatatypeCompetitiveGroupsApplicationsRating  Datatype = "competitive_groups_applications_rating"
	DatatypeEntrantPhotoFiles                    Datatype = "entrant_photo_files"
)

var AllDataType = []Datatype{
	DatatypeSubdivisionOrg,
	DatatypeCampaign,
	DatatypeAchievements,
	DatatypeAdmissionVolume,
	DatatypeDistributedAdmissionVolume,
	DatatypeCompetitiveGroups,
	DatatypeCompetitiveGroupPrograms,
	DatatypeCompetitiveBenefits,
	DatatypeEntranceTests,
	DatatypeEntranceTestBenefits,
	DatatypeEntrants,
	DatatypeCompatriot,
	DatatypeComposition,
	DatatypeDisability,
	DatatypeEducations,
	DatatypeEge,
	DatatypeIdentification,
	DatatypeMilitaries,
	DatatypeOlympics,
	DatatypeOrphans,
	DatatypeOther,
	DatatypeParentsLost,
	DatatypeRadiationWork,
	DatatypeVeteran,
	DatatypeApplications,
	DatatypeEditApplicationStatus,
	DatatypeEntranceTestAgreed,
	DatatypeEntranceTestResult,
	DatatypeOrderAdmission,
	DatatypeCompletitiveGroupsApplicationsRating,
	DatatypeAppAchievements,
	DatatypeApplicationsRating,
	DatatypeCompetitiveGroupsApplicationsRating,
	DatatypeEntrantPhotoFiles,
}

func (e Datatype) IsValid() bool {
	switch e {
	case DatatypeSubdivisionOrg,
		DatatypeCampaign,
		DatatypeAchievements,
		DatatypeAdmissionVolume,
		DatatypeDistributedAdmissionVolume,
		DatatypeCompetitiveGroups,
		DatatypeCompetitiveGroupPrograms,
		DatatypeCompetitiveBenefits,
		DatatypeEntranceTests,
		DatatypeEntranceTestBenefits,
		DatatypeEntrants,
		DatatypeCompatriot,
		DatatypeComposition,
		DatatypeDisability,
		DatatypeEducations,
		DatatypeEge,
		DatatypeIdentification,
		DatatypeMilitaries,
		DatatypeOlympics,
		DatatypeOrphans,
		DatatypeOther,
		DatatypeParentsLost,
		DatatypeRadiationWork,
		DatatypeVeteran,
		DatatypeApplications,
		DatatypeEditApplicationStatus,
		DatatypeEntranceTestAgreed,
		DatatypeEntranceTestResult,
		DatatypeOrderAdmission,
		DatatypeCompletitiveGroupsApplicationsRating,
		DatatypeAppAchievements,
		DatatypeApplicationsRating,
		DatatypeCompetitiveGroupsApplicationsRating,
		DatatypeEntrantPhotoFiles:
		return true
	}
	return false
}

func (e Datatype) String() string {
	return string(e)
}

type ActionMessage struct {
	SignMessage
}

func NewActionMessage(crypto sspvo.Crypto, action Action, datatype Datatype, data []byte) *ActionMessage {
	msg := &ActionMessage{}
	msg.Init(crypto, data)
	msg.UpdateJWTFields(setAction(action), setDatatype(datatype))

	return msg
}

func (m *ActionMessage) PathMethod() string {
	return pathMethodAction
}

func (m *ActionMessage) Response() sspvo.Response {
	resp := response.NewSignResponse()
	resp.SetCryptoHandler(m.crypto.GetVerifyCrypto)
	return resp
}
