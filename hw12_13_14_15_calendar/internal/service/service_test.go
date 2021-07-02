package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/require"

	pb "github.com/b2r2/hw/hw12_13_14_15_calendar/pkg/service"

	grpcserver "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/server/grpc"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/sirupsen/logrus"
)

func TestService(t *testing.T) {
	ctx, stop := context.WithTimeout(context.Background(), time.Second*5)
	defer stop()

	log := logrus.New()
	db := memorystorage.New(log)
	g := grpcserver.NewGRPCServer(log)
	s := NewService(logrus.New(), app.New(log, db))

	pb.RegisterCalendarServer(g, s)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		require.NoError(t, g.Start(":8081"))
	}()

	client, cc, err := func() (pb.CalendarClient, *grpc.ClientConn, error) {
		cc, err := grpc.Dial(":8081", grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		client := pb.NewCalendarClient(cc)
		return client, cc, nil
	}()
	defer func() {
		require.NoError(t, cc.Close())
	}()

	require.NoError(t, err)

	var id int32 = 1
	res, err := client.Create(ctx, newEvent(id, time.Now().Add(time.Minute)))
	require.NoError(t, err)
	require.Equal(t, res.Id, int32(1))

	id++
	e := newEvent(id, time.Now().Add(time.Hour))
	e.Notification = nil

	res, err = client.Create(ctx, e)
	require.NoError(t, err)
	require.Equal(t, res.Id, int32(2))

	getEvent, err := client.Get(ctx, &pb.EventID{Id: e.GetId()})
	require.NoError(t, err)
	require.Equal(t, id, getEvent.GetId())

	e.Description = "updated"
	_, err = client.Update(ctx, e)
	require.NoError(t, err)
	getEvent, err = client.Get(ctx, &pb.EventID{Id: e.GetId()})
	require.NoError(t, err)
	require.Equal(t, "updated", getEvent.GetDescription())

	_, err = client.DeleteAll(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	g.Stop()
	wg.Wait()
}

func newEvent(id int32, t time.Time) *pb.Event {
	return &pb.Event{
		Id:           id,
		Title:        "event",
		Start:        timestamppb.New(t),
		Stop:         timestamppb.New(t),
		Description:  "some desc",
		UserId:       id,
		Notification: durationpb.New(time.Hour),
	}
}
