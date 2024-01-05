package servers

import (
	"context"
	"sync/atomic"

	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	"go.uber.org/zap"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

// =============== server ====================== //

func NewRemoteServer(cl *itisadb.Client, number int32, logger *zap.Logger) *RemoteServer {
	return &RemoteServer{
		sdk:    cl,
		number: number,
		tries:  atomic.Uint32{},
		ram:    gost.NewRwLock(models.RAM{}),
		logger: logger,
	}
}

type RemoteServer struct {
	tries  atomic.Uint32
	ram    gost.RwLock[models.RAM]
	number int32
	logger *zap.Logger

	sdk *itisadb.Client
}

func (s *RemoteServer) Number() int32 {
	return s.number
}

func (s *RemoteServer) IsOffline() bool {
	return s.tries.Load() >= constants.MaxServerTries
}

type resulterr interface {
	IsErr() bool
	Error() *gost.Error
}

func after[Re resulterr, RePtr *Re](s *RemoteServer, res RePtr) {
	if res == nil { // TODO: may cause problems
		s.logger.Warn("res in server handler is nil", zap.Int32("server number", s.number))
		return
	}

	resUnwrapped := *res

	if !resUnwrapped.IsErr() {
		s.resetTries()
		return
	}

	if resUnwrapped.Error().BaseCode() != 0 {
		return
	}

	s.incTries()
}

func (s *RemoteServer) GetOne(ctx context.Context, _ int, key string, opt models.GetOptions) (res gost.Result[models.Value]) {
	defer after(s, &res)

	r := s.sdk.GetOne(ctx, key, opt.ToSDK())
	if r.IsErr() {
		return res.Err(r.Error())
	}

	val := r.Unwrap()
	return res.Ok(models.Value{
		ReadOnly: val.ReadOnly,
		Level:    models.Level(val.Level),
		Value:    val.Value,
	})
}

func (s *RemoteServer) DelOne(ctx context.Context, _ int, key string, opt models.DeleteOptions) (res gost.Result[gost.Nothing]) {
	defer after(s, &res)
	return s.sdk.DelOne(ctx, key, opt.ToSDK())
}

func (s *RemoteServer) SetOne(ctx context.Context, _ int, key string, val string, opts models.SetOptions) (res gost.Result[int32]) {
	defer after(s, &res)
	return s.sdk.SetOne(ctx, key, val, opts.ToSDK())
}

func (s *RemoteServer) RAM() models.RAM {
	defer s.ram.Release()
	return s.ram.RBorrow().Read()
}

func (s *RemoteServer) RefreshRAM(ctx context.Context) (res gost.Result[gost.Nothing]) {
	defer after(s, &res)

	r := itisadb.Internal.GetRAM(ctx, s.sdk)
	if r.IsErr() {
		return res.Err(r.Error())
	}

	s.ram.SetWithLock(models.RAM(r.Unwrap()))
	return gost.Ok(gost.Nothing{})
}

func (s *RemoteServer) incTries() uint32 {
	return s.tries.Add(1)
}

func (s *RemoteServer) resetTries() {
	s.tries.Store(0)
}

func (s *RemoteServer) NewObject(ctx context.Context, _ int, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing]) {
	defer after(s, &res)

	r := s.sdk.Object(ctx, name, opts.ToSDK())
	if r.IsOk() {
		return res.Ok(gost.Nothing{})
	}

	return res.Err(r.Error())
}

func (s *RemoteServer) GetFromObject(ctx context.Context, _ int, object string, key string, opts models.GetFromObjectOptions) (res gost.Result[string]) {
	defer after(s, &res)

	r := s.sdk.Object(ctx, object)
	if r.IsErr() {
		return res.Err(r.Error().WrapfNotNilMsg("error while GetFromObject [%s]", object))
	}

	gerRes := r.Unwrap().Get(ctx, key, opts.ToSDK())
	if gerRes.IsErr() {
		return res.Err(gerRes.Error().WrapfNotNilMsg("error while GetFromObject [%s.%s]", object, key))
	}

	return res.Ok(gerRes.Unwrap())
}

func (s *RemoteServer) SetToObject(ctx context.Context, userID int, object string, key string, value string, opts models.SetToObjectOptions) (res gost.Result[gost.Nothing]) {
	defer after(s, &res)

	r := s.sdk.Object(ctx, object)
	if r.IsErr() {
		return res.Err(r.Error())
	}

	setResult := r.Unwrap().Set(ctx, key, value, opts.ToSDK())
	if setResult.IsErr() {
		return res.Err(setResult.Error())
	}

	return res.Ok(gost.Nothing{})
}

func (s *RemoteServer) ObjectToJSON(ctx context.Context, userID int, name string, opts models.ObjectToJSONOptions) (res gost.Result[string]) {
	defer after(s, &res)

	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) ObjectSize(ctx context.Context, userID int, object string, opts models.SizeOptions) (res gost.Result[uint64]) {
	defer after(s, &res)

	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) DeleteObject(ctx context.Context, userID int, object string, opts models.DeleteObjectOptions) (res gost.ResultN) {
	defer after(s, &res)

	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) AttachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) (res gost.ResultN) {
	defer after(s, &res)

	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) ObjectDeleteKey(ctx context.Context, userID int, object, key string, opts models.DeleteAttrOptions) (res gost.ResultN) {
	defer after(s, &res)

	//TODO implement me
	panic("implement me")
}
