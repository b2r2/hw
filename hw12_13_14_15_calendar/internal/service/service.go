package service

import (
	"context"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	pb "github.com/b2r2/hw/hw12_13_14_15_calendar/pkg/service"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedCalendarServer
	app app.App
	log logger.Logger
}

func NewService(log logger.Logger, app app.App) pb.CalendarServer {
	return &Service{log: log, app: app}
}

func (s *Service) Create(ctx context.Context, event *pb.Event) (*pb.EventID, error) {
	id, err := s.app.CreateEvent(ctx, convertPb2Event(event))
	if err != nil {
		return nil, err
	}

	s.log.Info("create new event", event.GetId())

	return &pb.EventID{Id: id}, nil
}

func (s *Service) Update(ctx context.Context, event *pb.Event) (*emptypb.Empty, error) {
	if err := s.app.UpdateEvent(ctx, event.GetId(), convertPb2Event(event)); err != nil {
		return &emptypb.Empty{}, err
	}

	s.log.Info("update event", event.GetId())

	return &emptypb.Empty{}, nil
}

func (s *Service) Delete(ctx context.Context, id *pb.EventID) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, id.GetId()); err != nil {
		return &emptypb.Empty{}, err
	}

	s.log.Info("delete event", id.GetId())

	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteAll(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.app.DeleteAllEvent(ctx); err != nil {
		return &emptypb.Empty{}, err
	}

	s.log.Info("remove all events")

	return &emptypb.Empty{}, nil
}

func (s *Service) Get(ctx context.Context, id *pb.EventID) (*pb.Event, error) {
	event, err := s.app.GetEvent(ctx, id.GetId())
	if err != nil {
		return nil, err
	}

	s.log.Info("get event", id.GetId())

	return convertEvent2Pb(event), nil
}

func (s *Service) ListAll(ctx context.Context, _ *emptypb.Empty) (*pb.AllEvents, error) {
	events, err := s.app.ListAllEvents(ctx)
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		eventsProto = append(eventsProto, convertEvent2Pb(e))
	}

	s.log.Info("get all events")

	return &pb.AllEvents{AllEvents: eventsProto}, nil
}

func (s *Service) ListDay(ctx context.Context, date *timestamppb.Timestamp) (*pb.AllEvents, error) {
	events, err := s.app.ListDayEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		eventsProto = append(eventsProto, convertEvent2Pb(e))
	}

	s.log.Info("get events day", date.AsTime())

	return &pb.AllEvents{AllEvents: eventsProto}, nil
}

func (s *Service) ListWeek(ctx context.Context, date *timestamppb.Timestamp) (*pb.AllEvents, error) {
	events, err := s.app.ListWeekEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		eventsProto = append(eventsProto, convertEvent2Pb(e))
	}

	s.log.Info("get events week", date.AsTime())

	return &pb.AllEvents{AllEvents: eventsProto}, nil
}

func (s *Service) ListMonth(ctx context.Context, date *timestamppb.Timestamp) (*pb.AllEvents, error) {
	events, err := s.app.ListMonthEvents(ctx, date.AsTime())
	if err != nil {
		return nil, err
	}
	eventsProto := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		eventsProto = append(eventsProto, convertEvent2Pb(e))
	}

	s.log.Info("get events month", date.AsTime())

	return &pb.AllEvents{AllEvents: eventsProto}, nil
}

func convertPb2Event(event *pb.Event) *storage.Event {
	duration := event.GetNotification().AsDuration()
	return &storage.Event{
		ID:           event.GetId(),
		Title:        event.GetTitle(),
		Start:        event.GetStart().AsTime(),
		Stop:         event.GetStop().AsTime(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserId(),
		Notification: &duration,
	}
}

func convertEvent2Pb(event *storage.Event) *pb.Event {
	return &pb.Event{
		Id:           event.ID,
		Title:        event.Title,
		Start:        timestamppb.New(event.Start),
		Stop:         timestamppb.New(event.Stop),
		Description:  event.Description,
		UserId:       event.UserID,
		Notification: durationpb.New(*event.Notification),
	}
}
