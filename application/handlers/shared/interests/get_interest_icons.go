package shared_interests

import (
	"chat_app_backend/application/models/interests/get"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/mapper"
	"chat_app_backend/internal/s3"
	"chat_app_backend/internal/sqlc/db_queries"
	"context"
)

func GetInterestIcons(rawInterests []db_queries.Interest, client s3.IClient, ctx context.Context) ([]get.GetInterestResponseDto, exceptions.ITrackableException) {
	mappedInterests := make([]get.GetInterestResponseDto, len(rawInterests))

	for idx, rawInterest := range rawInterests {
		link, s3Error := client.GetDownloadUrl(ctx, rawInterest.IconFileName, s3.InterestsIconBucket)
		if s3Error != nil {
			return nil, exceptions.WrapErrorWithTrackableException(s3Error)
		}

		mappingErr := mapper.Mapper{}.Map(
			&mappedInterests[idx],
			rawInterest,
			struct {
				IconDownloadLink string
			}{
				IconDownloadLink: link,
			},
		)

		if mappingErr != nil {
			return nil, exceptions.WrapErrorWithTrackableException(mappingErr)
		}
	}

	return mappedInterests, nil
}
