package servers

import (
	"context"
	"sync/atomic"

	"github.com/egorgasay/gost"
	"github.com/egorgasay/itisadb-go-sdk"
	"itisadb/internal/constants"
	"itisadb/internal/models"
)

// =============== server ====================== //

func NewRemoteServer(cl *itisadb.Client, number int32, control interface{ Disconnect(int32) }) *RemoteServer {
	return &RemoteServer{
		sdk:     cl,
		number:  number,
		tries:   atomic.Uint32{},
		ram:     gost.NewRwLock(models.RAM{}),
		control: control,
	}
}

type RemoteServer struct {
	tries   atomic.Uint32
	ram     gost.RwLock[models.RAM]
	number  int32
	control interface{ Disconnect(int32) }

	sdk *itisadb.Client
}

func (s *RemoteServer) Number() int32 {
	return s.number
}

func (s *RemoteServer) Tries() uint32 {
	return s.tries.Load()
}

func (rs *RemoteServer) OnServerError(err *gost.Error) {
	if err.BaseCode() != 0 {
		return
	}

	if rs.IncTries() > constants.MaxServerTries {
		rs.control.Disconnect(rs.Number())
	}
}

func (s *RemoteServer) GetOne(ctx context.Context, userID int, key string, opt models.GetOptions) (res gost.Result[models.Value]) {
	r := s.sdk.GetOne(ctx, key, opt.ToSDK())
	if r.IsErr() {
		s.OnServerError(r.Error())
		return res.Err(r.Error())
	}

	val := r.Unwrap()
	return res.Ok(models.Value{
		ReadOnly: val.ReadOnly,
		Level:    models.Level(val.Level),
		Value:    val.Value,
	})
}

func (s *RemoteServer) DelOne(ctx context.Context, userID int, key string, opt models.DeleteOptions) gost.Result[gost.Nothing] {
	r := s.sdk.DelOne(ctx, key, opt.ToSDK())
	if r.IsErr() {
		s.OnServerError(r.Error())
		return r
	}

	return r
}

func (s *RemoteServer) SetOne(ctx context.Context, userID int, key string, val string, opts models.SetOptions) (res gost.Result[int32]) {
	r := s.sdk.SetOne(ctx, key, val, opts.ToSDK())
	if r.IsErr() {
		s.OnServerError(r.Error())
		return res.Err(r.Error())
	}

	return r
}

func (s *RemoteServer) RAM() models.RAM {
	defer s.ram.Release()
	return s.ram.RBorrow().Read()
}

func (s *RemoteServer) RefreshRAM(ctx context.Context) (res gost.Result[gost.Nothing]) {
	r := itisadb.Internal.GetRAM(ctx, s.sdk)
	if r.IsErr() {
		s.OnServerError(r.Error())
		return res.Err(r.Error())
	}

	s.ram.SetWithLock(models.RAM(r.Unwrap()))
	return gost.Ok(gost.Nothing{})
}

func (s *RemoteServer) IncTries() uint32 {
	return s.tries.Add(1)
}

func (s *RemoteServer) ResetTries() {
	s.tries.Store(0)
}

func (s *RemoteServer) NewObject(ctx context.Context, userID int, name string, opts models.ObjectOptions) (res gost.Result[gost.Nothing]) {
	r := s.sdk.Object(ctx, name, opts.ToSDK())
	if r.IsOk() {
		return res.Ok(gost.Nothing{})
	}

	s.OnServerError(r.Error())
	return res.Err(r.Error())
}

func (s *RemoteServer) GetFromObject(ctx context.Context, userID int, object string, key string, opts models.GetFromObjectOptions) (res gost.Result[string]) {
	defer s.OnServerError(res.Error())

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
	defer s.OnServerError(res.Error())

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
	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) ObjectSize(ctx context.Context, userID int, object string, opts models.SizeOptions) (res gost.Result[uint64]) {
	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) DeleteObject(ctx context.Context, userID int, object string, opts models.DeleteObjectOptions) gost.ResultN {
	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) AttachToObject(ctx context.Context, userID int, dst, src string, opts models.AttachToObjectOptions) gost.ResultN {
	//TODO implement me
	panic("implement me")
}

func (s *RemoteServer) ObjectDeleteKey(ctx context.Context, userID int, object, key string, opts models.DeleteAttrOptions) gost.ResultN {
	//TODO implement me
	panic("implement me")
}
