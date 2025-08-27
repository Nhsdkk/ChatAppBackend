package interests

import (
	"chat_app_backend/application/models/interests/update"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/request_env"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/service_wrapper"
	"chat_app_backend/internal/sqlc/db_queries"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UpdateInterestsHandler struct{}

func (u UpdateInterestsHandler) Handle(
	request *update.UpdateInterestRequestDto,
	services service_wrapper.IServiceWrapper,
	ctx *gin.Context,
	requestEnvironment *request_env.RequestEnv,
) (*update.UpdateInterestResponseDto, exceptions.ITrackableException) {
	if requestEnvironment.User.Role != db_queries.RoleTypeADMIN {
		return nil, &common_exceptions.ForbiddenException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.CreateTrackableExceptionFromStringF("interest update forbidden"),
				Message:             "only admins can access this route",
			},
		}
	}

	interest, getInterestDbError := services.GetDbConnection().
		GetQueries().
		GetInterestById(ctx, request.ID)

	switch {
	case errors.Is(getInterestDbError, pgx.ErrNoRows):
		return nil, &common_exceptions.ResourceNotFoundException{
			BaseRestException: exceptions.BaseRestException{
				ITrackableException: exceptions.WrapErrorWithTrackableException(getInterestDbError),
				Message:             "interest not found",
			},
		}
	case getInterestDbError != nil:
		return nil, exceptions.WrapErrorWithTrackableException(getInterestDbError)
	}

	var downloadLink string

	if request.Icon != nil {
		newDownloadLink, uploadError := services.GetS3Client().
			ModifyFileContents(ctx, request.Icon, interest.IconFileName, s3.InterestsIconBucket)

		if uploadError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(uploadError)
		}

		downloadLink = newDownloadLink
	} else {
		storedDownloadLink, uploadError := services.GetS3Client().
			GetDownloadUrl(ctx, interest.IconFileName, s3.InterestsIconBucket)

		if uploadError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(uploadError)
		}

		downloadLink = storedDownloadLink
	}

	if request.Description != nil {
		var updateParams db_queries.UpdateInterestDescriptionParams

		updateParamsMappingError := mapper.Mapper{}.Map(
			&updateParams,
			*request,
			struct {
				Description string
			}{
				Description: *request.Description,
			},
		)

		if updateParamsMappingError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(updateParamsMappingError)
		}

		newInterest, descriptionUpdateError := services.GetDbConnection().
			GetQueries().
			UpdateInterestDescription(ctx, updateParams)

		if descriptionUpdateError != nil {
			return nil, exceptions.WrapErrorWithTrackableException(descriptionUpdateError)
		}

		interest = newInterest
	}

	var result update.UpdateInterestResponseDto

	resultMappingError := mapper.Mapper{}.Map(
		&result,
		interest,
		struct {
			IconDownloadLink string
		}{
			IconDownloadLink: downloadLink,
		},
	)

	if resultMappingError != nil {
		return nil, exceptions.WrapErrorWithTrackableException(resultMappingError)
	}

	return &result, nil
}
