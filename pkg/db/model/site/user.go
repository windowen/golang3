package site

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"

	"serverApi/pkg/tools/cast"

	"gorm.io/gorm"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/passwordhelper"
	"serverApi/pkg/tools/strhelper"
)

// Register 用户注册
func (s *Site) Register(ctx context.Context, req *site.RegisterReq, parentId int, countryCode, avatar string) (int, error) {
	var (
		err     error
		timeNow = time.Now()
		userId  int
	)

	invite := strhelper.NewGenerator(constant.UserInviteCodeLen)

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := &table.User{
			Avatar:      avatar,
			CountryCode: countryCode,
			Category:    constant.UserCategoryGeneral,
			ParentId:    parentId,
			Status:      constant.UserStatusNormal,
			LevelId:     1, // 默认1
			CreatedAt:   timeNow,
		}

		// 创建用户
		if err = tx.Create(user).Error; err != nil {
			return err
		}

		// 更新昵称、邀请码
		if err = tx.Model(user).Select("invite_code", "nickname").
			Updates(table.User{
				InviteCode: invite.Encode(cast.ToUint64(user.Id)),
				Nickname:   fmt.Sprintf("user_%d", user.Id),
			}).
			Error; err != nil {
			return err
		}

		// 创建用户认证信息
		if err = tx.Create(&table.UserAuth{
			UserId:   user.Id,
			AreaCode: req.GetAreaCode(),
			Mobile:   req.GetMobile(),
			Email:    req.GetEmail(),
			Password: passwordhelper.GenPassword(req.GetPassword()),
		}).Error; err != nil {
			return err
		}

		if address, privateKey, err := s.generateTronAddressAndPrivateKey(); err != nil {
			// 创建钱包
			tx.Create(&table.UserWallet{
				UserId:    user.Id,
				CreatedAt: timeNow,
			})
		} else {
			// 创建钱包
			tx.Create(&table.UserWallet{
				UserId:          user.Id,
				CreatedAt:       timeNow,
				Trc20Address:    address,
				Trc20PrivateKey: privateKey,
			})
		}

		userId = user.Id

		return nil
	})

	return userId, err
}

func (s *Site) FindUserInfo(ctx context.Context, opts ...table.UserWhereOption) (*table.User, error) {
	var (
		user  table.User
		query = s.DB.WithContext(ctx).Scopes(user.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&user).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Site) UpdateUserInfo(ctx context.Context, userID int, opts ...table.UserUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	err := s.DB.WithContext(ctx).
		Model(&table.User{}).
		Where(table.FieldUserId.Eq(userID)).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}

// generateTronAddressAndPrivateKey 生成一个 Tron 地址和私钥
func (s *Site) generateTronAddressAndPrivateKey() (string, string, error) {
	// 生成 ECDSA 私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("无法生成私钥: %v", err)
	}

	// 提取公钥
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// 执行 sha256 哈希
	sha256Hash := sha256.Sum256(pubKey)

	ripeHash := ripemd160.New()
	_, err = ripeHash.Write(sha256Hash[:])
	if err != nil {
		return "", "", fmt.Errorf("ripemd160 hash failed: %v", err)
	}
	pubKeyHash := ripeHash.Sum(nil)

	// 添加 Tron 地址的前缀 (0x41)
	address := append([]byte{0x41}, pubKeyHash...)

	// 计算校验和
	checksum := sha256.Sum256(address)
	checksum = sha256.Sum256(checksum[:])

	// 将校验和的前4字节添加到地址末尾
	address = append(address, checksum[:4]...)

	// 使用 Base58 编码
	tronAddress := base58.Encode(address)

	// 获取私钥的十六进制表示
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	return tronAddress, privateKeyHex, nil
}
