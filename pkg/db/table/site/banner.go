package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const TableNameBanner = "site_banner"

var (
	FieldBannerId          = field.NewInt(TableNameBanner, "id")
	FieldBannerSort        = field.NewInt(TableNameBanner, "sort")
	FieldBannerCategory    = field.NewInt(TableNameBanner, "category")
	FieldBannerCountryJson = field.NewString(TableNameBanner, "country_json")
	FieldBannerName        = field.NewString(TableNameBanner, "name")
	FieldBannerUri         = field.NewString(TableNameBanner, "uri")
	FieldBannerShowType    = field.NewInt(TableNameBanner, "show_type")
	FieldBannerExtInfo     = field.NewString(TableNameBanner, "ext_info")
	FieldBannerStatus      = field.NewInt(TableNameBanner, "status")
	FieldBannerCreatedAt   = field.NewTime(TableNameBanner, "created_at")
	FieldBannerUpdatedAt   = field.NewTime(TableNameBanner, "updated_at")
)

type Banner struct {
	Id          int       `gorm:"column:id" json:"id"`
	Sort        int       `gorm:"column:sort" json:"sort"`                 // 排序 默认0,越小越靠前
	Category    int       `gorm:"column:category" json:"category"`         // 类别(固定 1-热门 2-游戏 3-附近 4-推荐)
	CountryJson string    `gorm:"column:country_json" json:"country_json"` // 可用国家 如： ["IN", "VN", "BR", "ID"]
	Name        string    `gorm:"column:name" json:"name"`                 // 轮播名称
	Uri         string    `gorm:"column:uri" json:"uri"`                   // 图片地址
	ShowType    int       `gorm:"column:show_type" json:"show_type"`       // 展示类型 1-纯展示 2-外链 3-直播间
	ExtInfo     string    `gorm:"column:ext_info" json:"ext_info"`         // 扩展信息 外链地址/直播间id
	Status      int       `gorm:"column:status" json:"status"`             // 1-启用2-禁用 3-删除
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

func (*Banner) TableName() string {
	return TableNameBanner
}

func (a *Banner) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// BannerWhereOption 条件项
type BannerWhereOption func(*gorm.DB) *gorm.DB

// WhereBannerId 查询
func WhereBannerId(value int) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerId.Eq(value))
	}
}

// WhereBannerSort 查询排序 默认0,越小越靠前
func WhereBannerSort(value int) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerSort.Eq(value))
	}
}

// WhereBannerCategory 查询类别(固定 1-热门 2-游戏 3-附近 4-推荐)
func WhereBannerCategory(value int) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerCategory.Eq(value))
	}
}

// WhereBannerName 查询轮播名称
func WhereBannerName(value string) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerName.Eq(value))
	}
}

// WhereBannerUri 查询图片地址
func WhereBannerUri(value string) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerUri.Eq(value))
	}
}

// WhereBannerShowType 查询展示类型 1-纯展示 2-外链 3-直播间
func WhereBannerShowType(value int) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerShowType.Eq(value))
	}
}

// WhereBannerExtInfo 查询扩展信息 外链地址/直播间id
func WhereBannerExtInfo(value string) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerExtInfo.Eq(value))
	}
}

// WhereBannerStatus 查询1-启用2-禁用 3-删除
func WhereBannerStatus(value int) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerStatus.Eq(value))
	}
}

// WhereBannerCreatedAt 查询创建时间
func WhereBannerCreatedAt(value time.Time) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerCreatedAt.Eq(value))
	}
}

// WhereBannerUpdatedAt 查询更新时间
func WhereBannerUpdatedAt(value time.Time) BannerWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldBannerUpdatedAt.Eq(value))
	}
}

// BannerUpdateOption 修改项
type BannerUpdateOption func(map[string]interface{})

// SetBannerId 设置
func SetBannerId(value int) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerId.ColumnName().String()] = value
	}
}

// SetBannerSort 设置排序 默认0,越小越靠前
func SetBannerSort(value int) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerSort.ColumnName().String()] = value
	}
}

// SetBannerCategory 设置类别(固定 1-热门 2-游戏 3-附近 4-推荐)
func SetBannerCategory(value int) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerCategory.ColumnName().String()] = value
	}
}

// SetBannerName 设置轮播名称
func SetBannerName(value string) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerName.ColumnName().String()] = value
	}
}

// SetBannerUri 设置图片地址
func SetBannerUri(value string) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerUri.ColumnName().String()] = value
	}
}

// SetBannerShowType 设置展示类型 1-纯展示 2-外链 3-直播间
func SetBannerShowType(value int) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerShowType.ColumnName().String()] = value
	}
}

// SetBannerExtInfo 设置扩展信息 外链地址/直播间id
func SetBannerExtInfo(value string) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerExtInfo.ColumnName().String()] = value
	}
}

// SetBannerStatus 设置1-启用2-禁用 3-删除
func SetBannerStatus(value int) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerStatus.ColumnName().String()] = value
	}
}

// SetBannerCreatedAt 设置创建时间
func SetBannerCreatedAt(value time.Time) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerCreatedAt.ColumnName().String()] = value
	}
}

// SetBannerUpdatedAt 设置更新时间
func SetBannerUpdatedAt(value time.Time) BannerUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldBannerUpdatedAt.ColumnName().String()] = value
	}
}
