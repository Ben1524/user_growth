package ugserver

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
	"user_growth/models"
	"user_growth/pb"
	"user_growth/service"
)

// UgCoinServer 用户积分服务
type UgCoinServer struct {
	pb.UnimplementedUserCoinServer // 这是一个默认的实现，如果父类有未实现的方法，则使用这个默认实现在调用时会返回未实现错误
}

// ListTasks 获取所有的积分任务列表
func (s *UgCoinServer) ListTasks(ctx context.Context, in *pb.ListTasksRequest) (*pb.ListTasksReply, error) {
	log.Printf("UgCoinServer.ListTasksRequest=%+v\n", *in)
	//return nil, status.Errorf(codes.Unimplemented, "方法待实现")
	coinTaskSvc := service.NewCoinTaskService(ctx)
	datalist, err := coinTaskSvc.FindAll()
	if err != nil {
		return nil, err
	}
	dlist := make([]*pb.TbCoinTask, len(datalist))
	for i := range datalist {
		dlist[i] = models.CoinTaskToMessage(&datalist[i])
	}
	out := &pb.ListTasksReply{
		Datalist: dlist,
	}
	return out, nil
}

// UserCoinInfo 获取用户的积分信息
func (s *UgCoinServer) UserCoinInfo(ctx context.Context, in *pb.UserCoinInfoRequest) (*pb.UserCoinInfoReply, error) {
	log.Printf("UgCoinServer.UserCoinInfoRequest=%+v\n", *in)
	//return nil, status.Errorf(codes.Unimplemented, "方法待实现")
	coinUserSvc := service.NewCoinUserService(ctx)
	uid := int(in.Uid)
	data, err := coinUserSvc.GetByUid(uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取用户积分信息失败: %v", err)
	}
	d := models.CoinUserToMessage(data)
	out := &pb.UserCoinInfoReply{
		Data: d,
	}
	return out, nil
}

// UserDetails 获取用户的积分明细列表
func (s *UgCoinServer) UserDetails(ctx context.Context, in *pb.UserDetailsRequest) (*pb.UserDetailsReply, error) {
	log.Printf("UgCoinServer.UserDetailsRequest=%+v\n", *in)
	//return nil, status.Errorf(codes.Unimplemented, "方法待实现")
	uid := int(in.Uid)
	page := int(in.Page)
	size := int(in.Size)
	coinDetailSvc := service.NewCoinDetailService(ctx)
	datalist, total, err := coinDetailSvc.FindByUid(uid, page, size)
	if err != nil {
		return nil, err
	}
	dlist := make([]*pb.TbCoinDetail, len(datalist))
	for i := range datalist {
		dlist[i] = models.CoinDetailToMessage(&datalist[i])
	}
	out := &pb.UserDetailsReply{
		Datalist: dlist,
		Total:    int32(total),
	}
	return out, nil
}

// UserCoinChange 调整用户积分-奖励和惩罚都是这个接口
func (s *UgCoinServer) UserCoinChange(ctx context.Context, in *pb.UserCoinChangeRequest) (*pb.UserCoinChangeReply, error) {
	log.Printf("UgCoinServer.UserCoinChangeRequest=%+v\n", *in)
	//return nil, status.Errorf(codes.Unimplemented, "方法待实现")
	uid := int(in.Uid)
	task := in.Task
	coin := int(in.Coin)
	taskInfo, err := service.NewCoinTaskService(ctx).GetByTask(task)
	if err != nil {
		return nil, err
	}
	if taskInfo == nil {
		return nil, errors.New("任务不存在")
	}
	// 插入详情
	coinDetail := models.TbCoinDetail{
		Uid:    uid,
		TaskId: taskInfo.Id,
		Coin:   coin,
	}
	err = service.NewCoinDetailService(ctx).Save(&coinDetail)
	if err != nil {
		return nil, err
	}
	// 更新用户信息
	coinUserSvc := service.NewCoinUserService(ctx)
	coinUser, err := coinUserSvc.GetByUid(uid)
	if err != nil {
		return nil, err
	}
	if coinUser == nil {
		coinUser = &models.TbCoinUser{
			Uid:   uid,
			Coins: coin,
		}
	} else {
		coinUser.Coins += coin
		coinUser.SysCreated = time.Time{}
		coinUser.SysUpdated = time.Now()
	}

	err = coinUserSvc.Save(coinUser)
	if err != nil {
		return nil, err
	}
	out := &pb.UserCoinChangeReply{
		User: models.CoinUserToMessage(coinUser),
	}
	return out, nil
}
