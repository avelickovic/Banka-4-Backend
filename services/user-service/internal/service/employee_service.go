package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/permission"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/user-service/internal/config"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/user-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/user-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/user-service/internal/repository"
)

type EmployeeService struct {
	employeeRepo        repository.EmployeeRepository
	identityRepo        repository.IdentityRepository
	activationTokenRepo repository.ActivationTokenRepository
	positionRepo        repository.PositionRepository
	emailService        Mailer
	cfg                 *config.Configuration
	txManager           repository.TransactionManager
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	identityRepo repository.IdentityRepository,
	activationTokenRepo repository.ActivationTokenRepository,
	positionRepo repository.PositionRepository,
	emailService Mailer,
	cfg *config.Configuration,
	txManager repository.TransactionManager,
) *EmployeeService {
	return &EmployeeService{
		employeeRepo:        employeeRepo,
		identityRepo:        identityRepo,
		activationTokenRepo: activationTokenRepo,
		positionRepo:        positionRepo,
		emailService:        emailService,
		cfg:                 cfg,
		txManager:           txManager,
	}
}

func (s *EmployeeService) Register(ctx context.Context, req *dto.CreateEmployeeRequest) (*model.Employee, error) {
	emailExists, err := s.identityRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	if emailExists {
		return nil, errors.ConflictErr("email already in use")
	}

	usernameExists, err := s.identityRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	if usernameExists {
		return nil, errors.ConflictErr("username already in use")
	}

	positionValid, err := s.positionRepo.Exists(ctx, req.PositionID)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	if !positionValid {
		return nil, errors.BadRequestErr("invalid position id")
	}

	identity := &model.Identity{
		Email:    req.Email,
		Username: req.Username,
		Type:     auth.IdentityEmployee,
		Active:   req.Active,
	}

	employee := &model.Employee{
		IdentityID:  0,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Department:  req.Department,
		PositionID:  req.PositionID,
		Permissions: mapPermissions(0, req.Permissions),
	}

	tokenStr, err := generateSecureToken(16)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	activationToken := &model.ActivationToken{
		IdentityID: 0,
		Token:      tokenStr,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.identityRepo.Create(txCtx, identity); err != nil {
			return errors.InternalErr(err)
		}

		employee.IdentityID = identity.ID
		if err := s.employeeRepo.Create(txCtx, employee); err != nil {
			return errors.InternalErr(err)
		}

		activationToken.IdentityID = identity.ID
		if err := s.activationTokenRepo.Create(txCtx, activationToken); err != nil {
			return errors.InternalErr(err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	activationBase := strings.TrimRight(s.cfg.URLs.FrontendBaseURL, "/")
	link := fmt.Sprintf("%s/activate?token=%s", activationBase, url.QueryEscape(tokenStr))

	if err := s.emailService.Send(
		identity.Email,
		"Welcome!",
		fmt.Sprintf("Kliknite ovde da postavite lozinku: %s", link),
	); err != nil {
		return nil, errors.ServiceUnavailableErr(err)
	}

	employee.Identity = *identity
	return employee, nil
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, id uint, req *dto.UpdateEmployeeRequest) (*model.Employee, error) {
	employee, err := s.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	if employee == nil {
		return nil, errors.NotFoundErr("employee not found")
	}

	identity, err := s.identityRepo.FindByID(ctx, employee.IdentityID)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	if identity == nil {
		return nil, errors.NotFoundErr("identity not found")
	}

	// TODO: Ability for other admins, or the admin themselves, to update their details?
	if employee.IsAdmin() {
		return nil, errors.ForbiddenErr("cannot modify admin")
	}

	identityChanged := false

	if req.Email != nil && *req.Email != identity.Email {
		emailExists, err := s.identityRepo.EmailExists(ctx, *req.Email)
		if err != nil {
			return nil, errors.InternalErr(err)
		}

		if emailExists {
			return nil, errors.ConflictErr("email already in use")
		}

		identity.Email = *req.Email
		identityChanged = true
	}

	if req.Username != nil && *req.Username != identity.Username {
		usernameExists, err := s.identityRepo.UsernameExists(ctx, *req.Username)
		if err != nil {
			return nil, errors.InternalErr(err)
		}

		if usernameExists {
			return nil, errors.ConflictErr("username already in use")
		}

		identity.Username = *req.Username
		identityChanged = true
	}

	if req.Active != nil && *req.Active != identity.Active {
		identity.Active = *req.Active
		identityChanged = true
	}

	if req.PositionID != nil && *req.PositionID != employee.PositionID {
		exists, err := s.positionRepo.Exists(ctx, *req.PositionID)
		if err != nil {
			return nil, errors.InternalErr(err)
		}

		if !exists {
			return nil, errors.BadRequestErr("invalid position_id")
		}
	}

	setIfNotNil(&employee.FirstName, req.FirstName)
	setIfNotNil(&employee.LastName, req.LastName)
	setIfNotNil(&employee.Gender, req.Gender)
	setIfNotNil(&employee.DateOfBirth, req.DateOfBirth)
	setIfNotNil(&employee.PhoneNumber, req.PhoneNumber)
	setIfNotNil(&employee.Address, req.Address)
	setIfNotNil(&employee.Department, req.Department)
	setIfNotNil(&employee.PositionID, req.PositionID)

	if req.Permissions != nil {
		employee.Permissions = mapPermissions(employee.EmployeeID, *req.Permissions)
	}

	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if identityChanged {
			if err := s.identityRepo.Update(txCtx, identity); err != nil {
				return errors.InternalErr(err)
			}
		}

		if err := s.employeeRepo.Update(txCtx, employee); err != nil {
			return errors.InternalErr(err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	employee.Identity = *identity
	return employee, nil
}

func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id uint) (*dto.EmployeeResponse, error) {
	employee, err := s.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.InternalErr(err)
	}
	if employee == nil {
		return nil, errors.NotFoundErr("employee not found")
	}

	return dto.ToEmployeeResponse(employee), nil
}

func (s *EmployeeService) GetAllEmployees(ctx context.Context, query *dto.ListEmployeesQuery) (*dto.ListEmployeesResponse, error) {
	employees, total, err := s.employeeRepo.GetAll(ctx, query.Email, query.FirstName, query.LastName, query.Position, query.Page, query.PageSize)
	if err != nil {
		return nil, errors.InternalErr(err)
	}

	return dto.ToEmployeeResponseList(employees, total, query.Page, query.PageSize), nil
}

func setIfNotNil[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

func mapPermissions(employeeID uint, permissions []permission.Permission) []model.EmployeePermission {
	result := make([]model.EmployeePermission, len(permissions))
	for i, p := range permissions {
		result[i] = model.EmployeePermission{
			EmployeeID: employeeID,
			Permission: p,
		}
	}
	return result
}

func (s *EmployeeService) DeactivateEmployee(ctx context.Context, id uint) error {
	employee, err := s.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return errors.InternalErr(err)
	}
	if employee == nil {
		return errors.NotFoundErr("employee not found")
	}

	if employee.IsAdmin() {
		return errors.ForbiddenErr("cannot deactivate admin")
	}

	identity, err := s.identityRepo.FindByID(ctx, employee.IdentityID)
	if err != nil {
		return errors.InternalErr(err)
	}
	if identity == nil {
		return errors.NotFoundErr("identity not found")
	}

	identity.Active = false

	if err := s.identityRepo.Update(ctx, identity); err != nil {
		return errors.InternalErr(err)
	}

	return nil
}
