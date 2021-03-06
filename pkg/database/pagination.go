// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"math"

	"github.com/google/exposure-notifications-verification-server/pkg/pagination"
	"github.com/jinzhu/gorm"
)

// Paginate is a helper that paginates a gorm query into the given result. In
// addition to reflecting into the provided result, it returns a pagination
// struct.
//
// If page is 0, it defaults to 1. If limit is 0, it defaults to the global
// pagination default limit.
func Paginate(query *gorm.DB, result interface{}, page, limit uint64) (*pagination.Paginator, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = pagination.DefaultLimit
	}

	var total uint64
	if err := query.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// Calculate offset.
	offset := uint64((page - 1) * limit)

	// Get the list of records, limiting and offsetting as needed.
	if err := query.
		Limit(limit).
		Offset(offset).
		Find(result).Error; err != nil {
		return nil, err
	}

	// Calculate how many future pages exist.
	var nextPages uint64
	if current := offset + limit; current < total {
		remaining := total - current
		nextPages = uint64(math.Ceil(float64(remaining) / float64(limit)))
		if nextPages > 10 {
			nextPages = 10
		}
	}

	prevPages := page - 1
	if prevPages > 10 {
		prevPages = 10
	}

	var paginator pagination.Paginator
	paginator.CurrPage = &pagination.Page{
		Number: page,
	}

	for i := prevPages; i > 0; i-- {
		paginator.PrevPages = append(paginator.PrevPages, &pagination.Page{
			Number: page - i,
		})
	}
	if prevPages > 0 {
		paginator.PrevPage = &pagination.Page{
			Number: page - 1,
		}
	}

	for i := uint64(1); i <= nextPages; i++ {
		paginator.NextPages = append(paginator.NextPages, &pagination.Page{
			Number: page + i,
		})
	}
	if nextPages > 0 {
		paginator.NextPage = &pagination.Page{
			Number: page + 1,
		}
	}

	return &paginator, nil
}
