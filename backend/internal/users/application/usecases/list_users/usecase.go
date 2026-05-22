package listusers

import (
	"context"
	"math"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/pkg/paginacion"
)

type ListarUsuariosCasoDeUso struct {
	userRepo domain.UserRepository
	authSvc  domain.AuthorizationService
}

func NewListarUsuariosCasoDeUso(userRepo domain.UserRepository, authSvc domain.AuthorizationService) *ListarUsuariosCasoDeUso {
	return &ListarUsuariosCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

func (uc *ListarUsuariosCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarUsuarios) (*RespuestaListarUsuarios, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.ViewUsers)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// Obtener usuarios sin paginación para contar total
	allUsers, err := uc.userRepo.Listar(ctx, cmd.Especificacion, paginacion.Paginacion{})
	if err != nil {
		return nil, err
	}
	total := int64(len(allUsers))

	// Obtener usuarios con paginación
	users, err := uc.userRepo.Listar(ctx, cmd.Especificacion, cmd.Paginacion)
	if err != nil {
		return nil, err
	}

	// Convertir a DTOs
	dtos := make([]UsuarioDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, UsuarioDTO{
			ID:        u.ID(),
			Email:     u.Email(),
			FullName:  u.FullName(),
			Active:    u.IsActive(),
			CreatedAt: u.CreatedAt(),
		})
	}

	// Calcular total de páginas
	pageSize := cmd.Paginacion.Limit()
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &RespuestaListarUsuarios{
		Usuarios:     dtos,
		Pagina:       cmd.Paginacion.Pagina,
		TamanoPagina: pageSize,
		Total:        total,
		TotalPages:   totalPages,
	}, nil
}
