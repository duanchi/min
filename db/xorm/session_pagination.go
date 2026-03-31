package xorm

import (
	"strconv"

	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/types"
)

func (session *Session) ListPage(records any, sess *Session, pageAndSize ...any) (result types.PaginationData, err error) {
	page := 1
	size := 20

	if len(pageAndSize) == 1 {
		switch pageAndSize[0].(type) {
		case *context.Context:
			page, _ = strconv.Atoi(pageAndSize[0].(*context.Context).Query("page", "1"))
			size, _ = strconv.Atoi(pageAndSize[0].(*context.Context).Query("size", "20"))
		case map[string]string:
			page, _ = strconv.Atoi(pageAndSize[0].(map[string]string)["page"])
			size, _ = strconv.Atoi(pageAndSize[0].(map[string]string)["size"])
		}
	} else if len(pageAndSize) == 2 {
		page = pageAndSize[0].(int)
		size = pageAndSize[1].(int)
	}

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	total, err := sess.Limit(size, (page-1)*size).FindAndCount(records)
	if err != nil {
		return
	}
	pages := int(total) / size
	if int(total)%size > 0 {
		pages++
	}
	result.Pagination = types.Pagination{
		Total:   int(total),
		Size:    size,
		Pages:   pages,
		Current: page,
	}

	result.Records = records
	return
}
