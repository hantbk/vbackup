package policy

import (
	"fmt"
	"time"

	"github.com/asdine/storm/q"
	"github.com/hantbk/vbackup/internal/entity/v1/repository"
	"github.com/hantbk/vbackup/internal/service/v1/common"
)

type Service interface {
	common.DBService
	Create(policy *repository.ForgetPolicy, options common.DBOptions) error
	List(options common.DBOptions) ([]repository.ForgetPolicy, error)
	Search(repoId int, path string, options common.DBOptions) ([]repository.ForgetPolicy, error)
	Get(id int, options common.DBOptions) (*repository.ForgetPolicy, error)
	Delete(id int, options common.DBOptions) error
	DeleteByRepo(repoId int, options common.DBOptions) error
	Update(policy *repository.ForgetPolicy, options common.DBOptions) error
	UpdateField(id int, fieldName string, value interface{}, options common.DBOptions) error
}

func GetService() Service {
	return &Policy{
		DefaultDBService: common.DefaultDBService{},
	}
}

type Policy struct {
	common.DefaultDBService
}

func (p Policy) DeleteByRepo(repoId int, options common.DBOptions) error {
	db := p.GetDB(options)
	policies, err := p.Search(repoId, "", options)
	if err != nil {
		return err
	}
	for policy := range policies {
		err = db.DeleteStruct(&policy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Policy) List(options common.DBOptions) (policies []repository.ForgetPolicy, err error) {
	db := p.GetDB(options)
	policies = make([]repository.ForgetPolicy, 0)
	var ms []q.Matcher
	query := db.Select(q.And(ms...)).OrderBy("CreatedAt").Reverse()
	if err = query.Find(&policies); err != nil {
		return
	}
	return
}

func (p Policy) Create(policy *repository.ForgetPolicy, options common.DBOptions) error {
	db := p.GetDB(options)
	policies, err := p.Search(policy.RepositoryId, policy.Path, options)
	if err != nil && err.Error() != "not found" {
		return err
	}
	if len(policies) > 0 {
		return fmt.Errorf("数据 %d,%s 已存在", policy.RepositoryId, policy.Path)
	}
	policy.CreatedAt = time.Now()
	return db.Save(policy)
}

func (p Policy) Search(repoId int, path string, options common.DBOptions) (policies []repository.ForgetPolicy, err error) {
	db := p.GetDB(options)
	policies = make([]repository.ForgetPolicy, 0)
	var ms []q.Matcher
	if repoId > 0 {
		ms = append(ms, q.Eq("RepositoryId", repoId))
	}
	if path != "" {
		ms = append(ms, q.Eq("Path", path))
	}
	query := db.Select(q.And(ms...)).OrderBy("CreatedAt").Reverse()
	if err = query.Find(&policies); err != nil {
		return
	}
	return
}

func (p Policy) Get(id int, options common.DBOptions) (*repository.ForgetPolicy, error) {
	db := p.GetDB(options)
	var policy repository.ForgetPolicy
	err := db.One("Id", id, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (p Policy) Delete(id int, options common.DBOptions) error {
	db := p.GetDB(options)
	policy, err := p.Get(id, options)
	if err != nil {
		return err
	}
	return db.DeleteStruct(policy)
}

func (p Policy) Update(policy *repository.ForgetPolicy, options common.DBOptions) error {
	db := p.GetDB(options)
	policy.UpdatedAt = time.Now()
	return db.Update(policy)
}

func (p Policy) UpdateField(id int, fieldName string, value interface{}, options common.DBOptions) error {
	db := p.GetDB(options)
	th := &repository.ForgetPolicy{}
	th.Id = id
	th.UpdatedAt = time.Now()
	return db.UpdateField(th, fieldName, value)
}
