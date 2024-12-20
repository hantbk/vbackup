package oplog

import (
	"time"

	"github.com/asdine/storm/q"
	"github.com/hantbk/vbackup/internal/entity/v1/oplog"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	"github.com/hantbk/vbackup/pkg/storm"
)

type Service interface {
	common.DBService
	Create(log *oplog.OperationLog, options common.DBOptions) error
	Search(num, size int, operator, operation, url, data string, options common.DBOptions) (int, []oplog.OperationLog, error)
}

func GetService() Service {
	return &OperationLog{
		DefaultDBService: common.DefaultDBService{},
	}
}

type OperationLog struct {
	common.DefaultDBService
}

func (o OperationLog) Create(log *oplog.OperationLog, options common.DBOptions) error {
	db := o.GetDB(options)
	log.CreatedAt = time.Now()
	return db.Save(log)
}

func (o OperationLog) Search(num, size int, operator, operation, url, data string, options common.DBOptions) (total int, res []oplog.OperationLog, err error) {
	total = 0
	res = make([]oplog.OperationLog, 0)
	db := o.GetDB(options)
	var ms []q.Matcher
	if operator != "" {
		ms = append(ms, q.Eq("Operator", operator))
	}
	if operation != "" {
		ms = append(ms, q.Eq("Operation", operation))
	}
	if url != "" {
		ms = append(ms, storm.Like("Url", url))
	}
	if data != "" {
		ms = append(ms, storm.Like("Data", data))
	}
	query := db.Select(q.And(ms...)).OrderBy("CreatedAt").Reverse()
	total, err = query.Count(&oplog.OperationLog{})
	if err != nil {
		return
	}
	if size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	if err = query.Find(&res); err != nil {
		return
	}
	return
}
