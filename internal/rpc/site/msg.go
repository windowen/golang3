package site

import (
	"context"

	"serverApi/pkg/common/mctx"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"
)

// RedPoint 是否红点以及未读数
func (s *siteSrv) RedPoint(ctx context.Context, _ *commonPb.Empty) (*pb.RedPointResp, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("RedPoint CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 取用户未读的交易信息
	count, err := s.siteDB.CountTransactionMsg(ctx, userId)
	if err != nil {
		zlogger.Errorf("RedPoint CountTransactionMsg err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	return &pb.RedPointResp{
		IsShowRedPoint: count > 0,
		UnreadCount:    int32(count),
	}, nil
}

// ReadAll 一键已读所有
func (s *siteSrv) ReadAll(ctx context.Context, _ *commonPb.Empty) (*commonPb.Empty, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("ReadAll CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("user_not_exist")
	}

	// 更新用户所有的交易消息为已读
	err = s.siteDB.UpdateTransactionMsg(ctx, userId)
	if err != nil {
		zlogger.Errorf("ReadAll UpdateTransactionMsg err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return &commonPb.Empty{}, nil
}

// MsgList 消息列表(平台消息、活动消息、交易消息)
func (s *siteSrv) MsgList(ctx context.Context, req *pb.MsgListReq) (*pb.MsgListResp, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("MsgList CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("user_not_exist")
	}

	langCode := mctx.GetLanguage(ctx)
	if len(langCode) == 0 {
		return nil, errs.ErrArgs.Wrap("parameter_error")
	}

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 || req.PageSize >= 50 {
		req.PageSize = 15
	}

	var count int64
	list := make([]*pb.MsgListItem, 0)

	// 1-平台消息2-交易消息3-活动消息
	switch req.MsgType {
	case 1:
		platformList, platformCount, _ := s.siteDB.PlatformMsgList(ctx, req.PageNumber, req.PageSize, langCode)
		for _, platformMsg := range platformList {
			formattedTime := platformMsg.CreatedAt.Format(utils.TimeBarFormat)
			list = append(list, &pb.MsgListItem{
				Id:        int32(platformMsg.Id),
				Content:   platformMsg.Content,
				CreatedAt: formattedTime,
			})
		}
		count = platformCount
	case 2:
		transactionList, transactionCount, _ := s.siteDB.TransactionMsgList(ctx, req.PageNumber, req.PageSize, userId)
		for _, transactionMsg := range transactionList {
			formattedTime := transactionMsg.CreatedAt.Format(utils.TimeBarFormat)
			list = append(list, &pb.MsgListItem{
				Id:         int32(transactionMsg.Id),
				Content:    transactionMsg.Content,
				CreatedAt:  formattedTime,
				ReadStatus: int32(transactionMsg.ReadStatus),
			})
		}
		count = transactionCount
	case 3:
		activityMsgList, activityMsgCount, _ := s.siteDB.ActivityMsgList(ctx, req.PageNumber, req.PageSize, langCode)
		for _, activityMsg := range activityMsgList {
			formattedTime := activityMsg.CreatedAt.Format(utils.TimeBarFormat)
			list = append(list, &pb.MsgListItem{
				Id:        int32(activityMsg.Id),
				Content:   activityMsg.Content,
				CreatedAt: formattedTime,
			})
		}
		count = activityMsgCount
	default:
		platformList, platformCount, _ := s.siteDB.PlatformMsgList(ctx, req.PageNumber, req.PageSize, langCode)
		for _, platformMsg := range platformList {
			formattedTime := platformMsg.CreatedAt.Format(utils.TimeBarFormat)
			list = append(list, &pb.MsgListItem{
				Id:        int32(platformMsg.Id),
				Content:   platformMsg.Content,
				CreatedAt: formattedTime,
			})
		}
		count = platformCount
	}

	return &pb.MsgListResp{
		List:        list,
		TotalRecord: int32(count),
		PageSize:    req.PageSize,
		PageNum:     req.PageNumber,
	}, nil
}

// MarkRead 标记交易消息已读
func (s *siteSrv) MarkRead(ctx context.Context, req *pb.MarkReadReq) (*commonPb.Empty, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("MarkRead CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("user_not_exist")
	}

	if len(req.GetIds()) == 0 {
		return nil, errs.ErrArgs.Wrap("parameter_error")
	}

	// 更新用户部分交易消息为已读
	err = s.siteDB.UpdatePartTransactionMsg(ctx, userId, req.GetIds())
	if err != nil {
		zlogger.Errorf("MarkRead UpdatePartTransactionMsg err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return &commonPb.Empty{}, nil
}

// Rotating 跑马灯
func (s *siteSrv) Rotating(ctx context.Context, req *pb.RotatingReq) (*pb.RotatingResp, error) {
	countryCode := mctx.GetCountryCode(ctx)
	if len(countryCode) == 0 {
		return nil, errs.ErrArgs.Wrap("parameter_error")
	}

	languageCode := mctx.GetLanguage(ctx)
	if len(languageCode) == 0 {
		return nil, errs.ErrArgs.Wrap("parameter_error")
	}

	list, err := s.siteDB.RotatingList(ctx, req.RotatingType, countryCode, languageCode)
	if err != nil {
		zlogger.Errorf("Rotating RotatingList err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	listItem := make([]*pb.RotatingItem, 0)
	for _, item := range list {
		listItem = append(listItem, &pb.RotatingItem{
			Id:   int32(item.Id),
			Name: item.Name,
		})
	}

	return &pb.RotatingResp{
		List: listItem,
	}, nil
}

// ClearTransactionMsgReadAll 清除交易已读消息
func (s *siteSrv) ClearTransactionMsgReadAll(ctx context.Context, _ *commonPb.Empty) (*commonPb.Empty, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("ClearTransactionMsgReadAll CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("user_not_exist")
	}

	// 更新用户所有的交易消息为已读
	err = s.siteDB.ClearTransactionMsgReadAll(ctx, userId)
	if err != nil {
		zlogger.Errorf("ClearTransactionMsgReadAll UpdateTransactionMsg err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return &commonPb.Empty{}, nil
}

// MsgSummary 消息汇总(列表) 同时返回交易、平台、活动的最后一条消息，用户未读的交易信息条数
func (s *siteSrv) MsgSummary(ctx context.Context, _ *commonPb.Empty) (*pb.MsgSummaryResp, error) {
	// 获取用户信息
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("MsgSummary CheckUser err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}
	resp := &pb.MsgSummaryResp{}

	// 取用户未读的交易信息
	transactionCount, err := s.siteDB.CountTransactionMsg(ctx, userId)
	resp.TransactionUnreadCount = int32(transactionCount)

	if err != nil {
		zlogger.Errorf("MsgSummary CountTransactionMsg err=%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	transaction, err := s.siteDB.LastTransactionMsg(ctx, userId)
	if err == nil && transaction != nil {
		resp.TransactionCreatedAt = transaction.CreatedAt.Format(utils.TimeBarFormat)
		resp.TransactionContent = transaction.Content
	}

	activityMsg, err := s.siteDB.LastActivityMsg(ctx, userId)
	if err == nil && activityMsg != nil {
		resp.ActivityCreatedAt = activityMsg.CreatedAt.Format(utils.TimeBarFormat)
		resp.ActivityContent = activityMsg.Content
	}

	platformMsg, err := s.siteDB.LastPlatformMsg(ctx, userId)
	if err == nil && platformMsg != nil {
		resp.PlatformCreatedAt = platformMsg.CreatedAt.Format(utils.TimeBarFormat)
		resp.PlatformContent = platformMsg.Content
	}
	return resp, nil
}
