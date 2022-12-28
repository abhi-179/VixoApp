package config

import (
	"net/url"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
)

func GeneratePaginationFromRequest(query url.Values) models.Pagination {
	// Initializing default
	//	var mode string
	limit := 20
	page := 0
	sort := "created_at asc"

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break

		}
	}
	return models.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}
