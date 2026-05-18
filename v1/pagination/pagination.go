package pagination

import (
	utils_v1 "github.com/makesalekz/utils/api/utils/v1"
)

// defaul limit for all services
const PAGINATION_LIMIT = 100

// InitializePagination initializes pagination with default limit if not provided
func InitializePagination(pagination *utils_v1.PaginateRequest) *utils_v1.PaginateRequest {
	if pagination == nil {
		return &utils_v1.PaginateRequest{
			Limit: PAGINATION_LIMIT,
		}
	}

	if pagination.Limit == 0 {
		pagination.Limit = PAGINATION_LIMIT
	}

	return pagination
}

// UpdatePaginationForAroundId resets pagination for around id, returns actual limit, because pagination.limit changes for second request
func UpdatePaginationForAroundId(pagination *utils_v1.PaginateRequest, lenList int32) int32 {
	// save actual limit, because pagination.limit changes for second request
	actualLimit := pagination.Limit

	// set from id to around id and reset around id for second request
	pagination.FromId = pagination.AroundId
	pagination.AroundId = 0

	// update limit for second request, half of limit or not enough part of list by first request
	if pagination.Limit != lenList {
		pagination.Limit = max(pagination.Limit-lenList, pagination.Limit/2)
	} else {
		pagination.Limit /= 2
	}

	return actualLimit
}

func GetListForAroundId[S ~[]E, E any](list1, list2 S, actualLimit int32, isDescDB bool, pagination *utils_v1.PaginateRequest) S {
	// assign length of lists to variables
	lenList1, lenList2 := int32(len(list1)), int32(len(list2))

	// if sum of lengths of lists is less than or equal to actual limit, return concatenated list
	if lenList1+lenList2 <= actualLimit {
		if pagination.Descending == isDescDB {
			return append(list1, list2...)
		}
		return append(list2, list1...)
	}

	// create variables for cutting lists
	cutList1, cutList2 := int32(0), int32(0)

	// calculate how many elements to cut from each list
	if len(list1) >= len(list2) {
		cutList2 = min(lenList2, actualLimit/2)
		cutList1 = actualLimit - cutList2
	} else {
		cutList1 = min(lenList1, actualLimit/2)
		cutList2 = actualLimit - cutList1
	}

	// return concatenated list with cut elements
	if pagination.Descending == isDescDB {
		return append(list1[lenList1-cutList1:], list2[:cutList2]...)
	}
	return append(list2[lenList2-cutList2:], list1[:cutList1]...)
}
