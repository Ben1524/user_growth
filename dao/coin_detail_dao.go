package dao

import (
	"context"
	comm "user_growth/comm"
	"user_growth/db"
	"user_growth/models"
	"xorm.io/xorm"
)

type CoinDetailDao struct {
	db  *xorm.Engine
	ctx context.Context
}

func NewCoinDetailDaoWithCtx(ctx context.Context) *CoinDetailDao {
	return &CoinDetailDao{
		db:  db.DbEngine,
		ctx: ctx,
	}
}

func NewCoinDetailDao(ctx context.Context, db *xorm.Engine) *CoinDetailDao {
	return &CoinDetailDao{
		db:  db,
		ctx: ctx,
	}
}
func (dao *CoinDetailDao) Get(id int) (*models.TbCoinDetail, error) {
	data := &models.TbCoinDetail{}
	if _, err := dao.db.ID(id).Get(data); err != nil {
		return nil, err
	} else if data == nil || data.Id == 0 {
		return nil, nil
	} else {
		return data, nil
	}
}

// FindByUid 查询用户的积分明细列表，会有多条
func (dao *CoinDetailDao) FindByUid(uid, page, size int) ([]models.TbCoinDetail, int64, error) {
	dataList := make([]models.TbCoinDetail, 0)
	sess := dao.db.Where("`uid`=?", uid)
	start := (page - 1) * size
	total, err := sess.Desc("id").Limit(size, start).FindAndCount(&dataList)
	return dataList, total, err
}

// FindAllPager get all models
func (dao *CoinDetailDao) FindAllPager(page, size int) ([]models.TbCoinDetail, int64, error) {
	datalist := make([]models.TbCoinDetail, 0)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 100
	}
	start := (page - 1) * size
	total, err := dao.db.Desc("id").Limit(size, start).FindAndCount(&datalist)
	return datalist, total, err
}

// Insert one row
func (dao *CoinDetailDao) Insert(data *models.TbCoinDetail) error {
	data.SysCreated = comm.Now()
	data.SysUpdated = comm.Now()
	_, err := dao.db.Insert(data)
	return err
}

// Update one row  ,mustColumns指定必须更新的列
// Xorm 的默认更新行为：“非零值才更新”
func (dao *CoinDetailDao) Update(data *models.TbCoinDetail, mustColumns ...string) error {
	sess := dao.db.ID(data.Id)
	if len(mustColumns) > 0 { // 如果指定了必须更新的列，则使用MustCols
		sess.MustCols(mustColumns...)
	}
	_, err := sess.Update(data) // 更新数据
	return err
}

// Save with Insert and Update
func (dao *CoinDetailDao) Save(data *models.TbCoinDetail, mustColumns ...string) error {
	if data.Id > 0 {
		return dao.Update(data, mustColumns...)
	} else {
		return dao.Insert(data)
	}
}
