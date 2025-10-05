package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
)

type InstanceService struct {
	tokenService     *TokenService
	instanceRepo     instance.InstanceRepository
	instanceRegistry instance.InstanceRegistry
	eventbus         events.EventBus
}

func NewInstanceService(tokenService *TokenService, instRepo instance.InstanceRepository, instRegistry instance.InstanceRegistry, eventbus events.EventBus) *InstanceService {
	return &InstanceService{
		tokenService:     tokenService,
		instanceRepo:     instRepo,
		instanceRegistry: instRegistry,
		eventbus:         eventbus,
	}
}

func (s *InstanceService) List(ctx context.Context) ([]*instance.Instance, *app.AppError) {
	instances, err := s.instanceRepo.List()
	if err != nil {
		return nil, app.NewDatabaseError("instance service", err)
	}

	return instances, nil
}

func (s *InstanceService) Create(ctx context.Context, inp input.CreateInstance) (*instance.Instance, *token.Token, *app.AppError) {
	l := app.GetInstanceServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, nil, app.TranslateError("instance service", err)
	}

	inst := instance.NewInstance(inp.Name)

	if err := s.instanceRepo.Insert(inst); err != nil {
		return nil, nil, app.NewDatabaseError("instance service", err)
	}

	token, appErr := s.tokenService.CreateToken(inst.ID)
	if appErr != nil {
		return nil, nil, app.WrapLoc("instance service", appErr)
	}

	go s.eventbus.Publish(inst.EventCreated())
	l.Info("Instance created", "instance", inst.ID)

	return inst, token, nil
}

func (s *InstanceService) Get(ctx context.Context, inp input.GetInstance) (*instance.Instance, *app.AppError) {
	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("instance service", err)
	}

	inst, ok := s.instanceRegistry.Get(inp.ID)
	if ok {
		return inst, nil
	}

	instance, err := s.instanceRepo.Get(instance.WhereID(inp.ID))
	if err != nil {
		return nil, app.NewDatabaseError("instance service", err)
	}

	return instance, nil
}

func (s *InstanceService) RenewToken(ctx context.Context, inp input.RenewInstanceToken) (*token.Token, *app.AppError) {
	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("instance service", err)
	}

	inst, appErr := s.Get(ctx, input.GetInstance(inp))

	if appErr != nil {
		return nil, app.WrapLoc("instance service", appErr)
	}

	if inst == nil {
		return nil, app.TranslateError("instance service", instance.ErrInstanceNotFound)
	}

	toks, err := s.tokenService.GetByInstanceID(inst.ID)
	if err != nil {
		return nil, app.NewDatabaseError("instance service", err)
	}

	if len(toks) <= 0 {
		tok, err := s.tokenService.CreateToken(inst.ID)
		if err != nil {
			return nil, app.WrapLoc("instance service", err)
		}
		return tok, nil
	}

	tok, appErr := s.tokenService.RenewToken(ctx, toks[0].ID, true)
	if appErr != nil {
		return nil, app.WrapLoc("instance service", appErr)
	}

	return tok, nil
}
