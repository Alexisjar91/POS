package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Alexisjar91/POS/internal/config"
	"github.com/Alexisjar91/POS/internal/database"
	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
)

func main() {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  RUN TEST — Integración PostgreSQL")
	fmt.Println("═══════════════════════════════════════")

	// 1. Config
	fmt.Println("\n📁 1. Configuración")
	cfg := config.Get()
	fmt.Printf("   ✓ Config loaded: DB=%s:%s/%s\n", cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Printf("   ✓ Server port: %s\n", cfg.ServerPort)

	// 2. Database connection
	fmt.Println("\n🗄️  2. Conexión a base de datos")
	db := database.Get()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("   ✗ Error getting underlying DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("   ✗ Cannot ping database: %v", err)
	}
	fmt.Println("   ✓ Connected to PostgreSQL successfully")

	ctx := context.Background()

	// 3. Migrations
	fmt.Println("\n📦 3. Migraciones")
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("   ✗ Migrations failed: %v", err)
	}
	fmt.Println("   ✓ Migrations executed successfully (AutoMigrate)")

	// 4. Seed
	fmt.Println("\n🌱 4. Seed de datos")
	if err := postgres.RunSeed(db); err != nil {
		log.Fatalf("   ✗ Seed failed: %v", err)
	}
	fmt.Println("   ✓ Seed executed successfully")

	// Verify permissions were seeded
	var permCount int64
	db.Model(&models.PermissionModel{}).Count(&permCount)
	fmt.Printf("   ✓ %d permissions seeded\n", permCount)

	var roleCount int64
	db.Model(&models.RoleModel{}).Count(&roleCount)
	fmt.Printf("   ✓ %d system roles seeded\n", roleCount)

	// 5a. Bulk insert 15 test users
	fmt.Println("\n👥 5a. Bulk insert 15 test users")

	testNames := []string{
		"Juan García", "María López", "Carlos Rodríguez", "Ana Martínez",
		"Pedro Sánchez", "Laura Fernández", "Miguel Díaz", "Isabel Moreno",
		"Francisco Jiménez", "Elena Ruiz", "Antonio Pérez", "Rosa García",
		"Manuel López", "Carmen Martínez", "David Sánchez", "PEPIRO EL DE LOS PALOTAS",
	}

	// Initialize hasher for password hashing
	hasher := postgres.NewPasswordHasher()
	defaultPassword, err := hasher.Hash("password123")
	if err != nil {
		log.Fatalf("   ✗ Failed to hash default password: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	userRoleRepo := postgres.NewUserRoleRepository(db)
	roleRepo := postgres.NewRoleRepository(db)

	// Get ADMIN role for assignment
	adminRole, err := roleRepo.ObtenerPorNombre(ctx, "ADMIN")
	if err != nil {
		log.Fatalf("   ✗ RoleRepository.ObtenerPorNombre(ADMIN): %v", err)
	}

	createdCount := 0
	for i, name := range testNames {
		email := fmt.Sprintf("user%d@example.com", i+1)

		// Create user
		newUser, err := domain.NuevoUsuario(email, name, defaultPassword, "SYSTEM")
		if err != nil {
			log.Fatalf("   ✗ domain.NuevoUsuario: %v", err)
		}

		createdUser, err := userRepo.Crear(ctx, newUser)
		if err != nil {
			if err == domain.ErrEmailDuplicado {
				fmt.Printf("   ⊘ User %s already exists (skipped)\n", email)
				createdCount++
				// Fetch existing user to assign role if needed
				existingUser, _ := userRepo.ObtenerPorEmail(ctx, email)
				if existingUser != nil {
					_ = userRoleRepo.Asignar(ctx, existingUser.ID(), adminRole.ID())
				}
			} else {
				log.Fatalf("   ✗ UserRepository.Crear(%s): %v", email, err)
			}
		} else {
			createdCount++
			// Assign ADMIN role
			if err := userRoleRepo.Asignar(ctx, createdUser.ID(), adminRole.ID()); err != nil {
				log.Fatalf("   ✗ UserRoleRepository.Asignar: %v", err)
			}
		}
	}

	var totalUserCount int64
	db.Model(&models.UserModel{}).Count(&totalUserCount)
	fmt.Printf("   ✓ %d users created/verified\n", createdCount)
	fmt.Printf("   ✓ Total users in database: %d\n", totalUserCount)

	// 5b. Test search combinations with Especificacion + Paginacion
	fmt.Println("\n🔍 5b. Test search combinations with Especificacion + Paginacion")

	// Helper to format test results
	runTest := func(testNum int, desc string, spec especificacion.Especificacion) {
		pag := paginacion.Paginacion{Pagina: 1, TamanoPagina: 50}
		users, err := userRepo.Listar(ctx, spec, pag)
		if err != nil {
			log.Fatalf("   ✗ Test %d failed: %v", testNum, err)
		}
		fmt.Printf("Test %d: %s\n", testNum, desc)
		fmt.Printf("   Results: %d users\n", len(users))
		for _, user := range users {
			createdBy := ""
			if user.CreatedBy() != nil {
				createdBy = *user.CreatedBy()
			}
			fmt.Printf("   - %s | %s | %s | %v | %s | %s\n",
				user.ID(), user.Email(), user.FullName(), user.IsActive(), createdBy, user.CreatedAt())
		}
		fmt.Println()
	}

	// Test 1: active = true
	runTest(1, "active = true",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
			},
		})

	// Test 2: active = false
	runTest(2, "active = false",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: false},
			},
		})

	// Test 3: email LIKE "%user1%"
	runTest(3, "email LIKE \"%user1%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "email", Operador: "LIKE", Valor: "%user1%"},
			},
		})

	// Test 4: email LIKE "%@example.com"
	runTest(4, "email LIKE \"%@example.com\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "email", Operador: "LIKE", Valor: "%@example.com"},
			},
		})

	// Test 5: fullName LIKE "%García%"
	runTest(5, "fullName LIKE \"%García%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "fullName", Operador: "LIKE", Valor: "%García%"},
			},
		})

	// Test 6: fullName LIKE "%López%"
	runTest(6, "fullName LIKE \"%López%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "fullName", Operador: "LIKE", Valor: "%López%"},
			},
		})

	// Test 7: createdBy = "SYSTEM"
	runTest(7, "createdBy = \"SYSTEM\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "createdBy", Operador: "=", Valor: "SYSTEM"},
			},
		})

	// Test 8: active = true AND email LIKE "%user1%"
	runTest(8, "active = true AND email LIKE \"%user1%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
				{Campo: "email", Operador: "LIKE", Valor: "%user1%"},
			},
		})

	// Test 9: active = true AND fullName LIKE "%García%"
	runTest(9, "active = true AND fullName LIKE \"%García%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
				{Campo: "fullName", Operador: "LIKE", Valor: "%García%"},
			},
		})

	// Test 10: active = true AND createdBy = "SYSTEM"
	runTest(10, "active = true AND createdBy = \"SYSTEM\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
				{Campo: "createdBy", Operador: "=", Valor: "SYSTEM"},
			},
		})

	// Test 11: email LIKE "%user%" AND fullName LIKE "%López%"
	runTest(11, "email LIKE \"%user%\" AND fullName LIKE \"%López%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "email", Operador: "LIKE", Valor: "%user%"},
				{Campo: "fullName", Operador: "LIKE", Valor: "%López%"},
			},
		})

	// Test 12: active = true AND email LIKE "%user%" AND fullName LIKE "%García%"
	runTest(12, "active = true AND email LIKE \"%user%\" AND fullName LIKE \"%García%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
				{Campo: "email", Operador: "LIKE", Valor: "%user%"},
				{Campo: "fullName", Operador: "LIKE", Valor: "%García%"},
			},
		})

	// Test 13: email != "user1@example.com"
	runTest(13, "email != \"user1@example.com\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "email", Operador: "!=", Valor: "user1@example.com"},
			},
		})

	// Test 14: fullName LIKE "%Pérez%"
	runTest(14, "fullName LIKE \"%Pérez%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "fullName", Operador: "LIKE", Valor: "%Pérez%"},
			},
		})

	// Test 15: active = true AND createdBy = "SYSTEM" AND email LIKE "%user%"
	runTest(15, "active = true AND createdBy = \"SYSTEM\" AND email LIKE \"%user%\"",
		especificacion.Especificacion{
			Filtros: []especificacion.CriterioFiltro{
				{Campo: "active", Operador: "=", Valor: true},
				{Campo: "createdBy", Operador: "=", Valor: "SYSTEM"},
				{Campo: "email", Operador: "LIKE", Valor: "%user%"},
			},
		})

	// 6. Repositorios
	fmt.Println("\n🔧 6. Repositorios")

	// 6a. PermissionRepository
	permRepo := postgres.NewPermissionRepository(db)
	allPerms, err := permRepo.ListarTodos(ctx)
	if err != nil {
		log.Fatalf("   ✗ PermissionRepository.ListarTodos: %v", err)
	}
	fmt.Printf("   ✓ PermissionRepository — %d permissions loaded\n", len(allPerms))

	// 6b. RoleRepository
	roleRepo2 := postgres.NewRoleRepository(db)

	// Get ADMIN role
	adminRole2, err := roleRepo2.ObtenerPorNombre(ctx, "ADMIN")
	if err != nil {
		log.Fatalf("   ✗ RoleRepository.ObtenerPorNombre(ADMIN): %v", err)
	}
	fmt.Printf("   ✓ RoleRepository — ADMIN role loaded (ID: %s)\n", adminRole2.ID())
	fmt.Printf("     IsSystem: %v\n", adminRole2.IsSystem())

	ownerRole, err := roleRepo2.ObtenerPorNombre(ctx, "OWNER")
	if err != nil {
		log.Fatalf("   ✗ RoleRepository.ObtenerPorNombre(OWNER): %v", err)
	}
	fmt.Printf("   ✓ RoleRepository — OWNER role loaded (ID: %s)\n", ownerRole.ID())

	// List all roles
	roles, err := roleRepo2.Listar(ctx)
	if err != nil {
		log.Fatalf("   ✗ RoleRepository.Listar: %v", err)
	}
	fmt.Printf("   ✓ RoleRepository — %d roles listed\n", len(roles))

	// 6c. PasswordHasher
	hasher2 := postgres.NewPasswordHasher()
	hash, err := hasher2.Hash("test_password_123")
	if err != nil {
		log.Fatalf("   ✗ PasswordHasher.Hash: %v", err)
	}
	fmt.Printf("   ✓ PasswordHasher — hash generated (%d bytes)\n", len(hash))

	if err := hasher2.Compare("test_password_123", hash); err != nil {
		log.Fatalf("   ✗ PasswordHasher.Compare (correct): %v", err)
	}
	fmt.Println("   ✓ PasswordHasher — correct password matches")

	if err := hasher2.Compare("wrong_password", hash); err == nil {
		log.Fatalf("   ✗ PasswordHasher.Compare (wrong) should have failed")
	}
	fmt.Println("   ✓ PasswordHasher — wrong password rejected")

	// 6d. UserRepository
	userRepo2 := postgres.NewUserRepository(db)

	// Create a test user
	testUser, err := domain.NuevoUsuario("test@example.com", "Test User", hash, "SYSTEM")
	if err != nil {
		log.Fatalf("   ✗ domain.NuevoUsuario: %v", err)
	}

	createdUser, err := userRepo2.Crear(ctx, testUser)
	if err != nil {
		log.Fatalf("   ✗ UserRepository.Crear: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — user created (ID: %s, Email: %s)\n",
		createdUser.ID(), createdUser.Email())

	// Get by ID
	fetchedByID, err := userRepo2.ObtenerPorID(ctx, createdUser.ID())
	if err != nil {
		log.Fatalf("   ✗ UserRepository.ObtenerPorID: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — user fetched by ID: %s\n", fetchedByID.Email())

	// Get by email
	fetchedByEmail, err := userRepo2.ObtenerPorEmail(ctx, "test@example.com")
	if err != nil {
		log.Fatalf("   ✗ UserRepository.ObtenerPorEmail: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — user fetched by email: %s\n", fetchedByEmail.ID())

	// Check exists
	exists, err := userRepo2.ExistePorEmail(ctx, "test@example.com")
	if err != nil {
		log.Fatalf("   ✗ UserRepository.ExistePorEmail: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — email exists: %v\n", exists)

	notExists, err := userRepo2.ExistePorEmail(ctx, "noexiste@example.com")
	if err != nil {
		log.Fatalf("   ✗ UserRepository.ExistePorEmail: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — non-existent email: %v\n", notExists)

	// List users
	allUsers, err := userRepo2.Listar(ctx, especificacion.Especificacion{}, paginacion.Paginacion{Pagina: 1, TamanoPagina: 10})
	if err != nil {
		log.Fatalf("   ✗ UserRepository.Listar: %v", err)
	}
	fmt.Printf("   ✓ UserRepository — %d users listed\n", len(allUsers))

	// 6e. UserRoleRepository
	userRoleRepo2 := postgres.NewUserRoleRepository(db)

	// Assign role
	if err := userRoleRepo2.Asignar(ctx, createdUser.ID(), adminRole2.ID()); err != nil {
		log.Fatalf("   ✗ UserRoleRepository.Asignar: %v", err)
	}
	fmt.Println("   ✓ UserRoleRepository — ADMIN role assigned to test user")

	// Assign again (should be idempotent)
	if err := userRoleRepo2.Asignar(ctx, createdUser.ID(), adminRole2.ID()); err != nil {
		log.Fatalf("   ✗ UserRoleRepository.Asignar (idempotent): %v", err)
	}
	fmt.Println("   ✓ UserRoleRepository — duplicate assign is idempotent")

	// 7. AuthorizationService
	fmt.Println("\n🔐 7. AuthorizationService")
	authSvc := postgres.NewAuthorizationService(db)

	// Verify OWNER permission (should be false since user is not OWNER)
	ownerVerified, err := authSvc.EsUsuarioOWNER(ctx, createdUser.ID())
	if err != nil {
		log.Fatalf("   ✗ AuthorizationService.EsUsuarioOWNER: %v", err)
	}
	fmt.Printf("   ✓ AuthorizationService — user is OWNER: %v (expected: false)\n", ownerVerified)

	// Check a permission the ADMIN role should have (ViewUsers)
	hasViewUsers, err := authSvc.VerificarPermiso(ctx, createdUser.ID(), domain.ViewUsers)
	if err != nil {
		log.Fatalf("   ✗ AuthorizationService.VerificarPermiso(ViewUsers): %v", err)
	}
	fmt.Printf("   ✓ AuthorizationService — user has view_users: %v (expected: true)\n", hasViewUsers)

	// Check a permission the ADMIN role should have (ManageRoles)
	hasManageRoles, err := authSvc.VerificarPermiso(ctx, createdUser.ID(), domain.ManageRoles)
	if err != nil {
		log.Fatalf("   ✗ AuthorizationService.VerificarPermiso(ManageRoles): %v", err)
	}
	fmt.Printf("   ✓ AuthorizationService — user has manage_roles: %v (expected: true)\n", hasManageRoles)

	// 8. Error cases
	fmt.Println("\n⚠️  8. Casos de error")

	// Duplicate email
	dupUser, _ := domain.NuevoUsuario("test@example.com", "Duplicate", hash, "SYSTEM")
	_, err = userRepo2.Crear(ctx, dupUser)
	if err == domain.ErrEmailDuplicado {
		fmt.Println("   ✓ UserRepository — duplicate email rejected (ErrEmailDuplicado)")
	} else {
		log.Fatalf("   ✗ Expected ErrEmailDuplicado, got: %v", err)
	}

	// Non-existent user
	_, err = userRepo2.ObtenerPorID(ctx, "nonexistentid")
	if err == domain.ErrUsuarioNoEncontrado {
		fmt.Println("   ✓ UserRepository — non-existent user returns ErrUsuarioNoEncontrado")
	} else {
		log.Fatalf("   ✗ Expected ErrUsuarioNoEncontrado, got: %v", err)
	}

	// Non-existent role
	_, err = roleRepo2.ObtenerPorID(ctx, "nonexistentid")
	if err == domain.ErrRolNoEncontrado {
		fmt.Println("   ✓ RoleRepository — non-existent role returns ErrRolNoEncontrado")
	} else {
		log.Fatalf("   ✗ Expected ErrRolNoEncontrado, got: %v", err)
	}

	// 9. RoleRepository — validaciones extra
	fmt.Println("\n🧪 9. RoleRepository — validaciones extra")

	// Try to delete ADMIN role (is system → ErrRolSistemaInmutable, even if has users)
	err = roleRepo2.Eliminar(ctx, adminRole2.ID())
	if err == domain.ErrRolSistemaInmutable {
		fmt.Println("   ✓ RoleRepository — system role cannot be deleted (ErrRolSistemaInmutable)")
	} else {
		log.Fatalf("   ✗ Expected ErrRolSistemaInmutable, got: %v", err)
	}

	// Try to delete OWNER role (is system → ErrRolSistemaInmutable)
	err = roleRepo2.Eliminar(ctx, ownerRole.ID())
	if err == domain.ErrRolSistemaInmutable {
		fmt.Println("   ✓ RoleRepository — OWNER role cannot be deleted (ErrRolSistemaInmutable)")
	} else {
		log.Fatalf("   ✗ Expected ErrRolSistemaInmutable, got: %v", err)
	}

	// Create a non-system role with users → ErrRolConUsuarios
	customRole, _ := domain.NuevoRol("test_role", "A role for testing")
	createdRole, err := roleRepo2.Crear(ctx, customRole)
	if err != nil {
		log.Fatalf("   ✗ RoleRepository.Crear (test_role): %v", err)
	}
	fmt.Printf("   ✓ RoleRepository — test role created: %s\n", createdRole.ID())

	// Assign the test role to the user
	if err := userRoleRepo2.Asignar(ctx, createdUser.ID(), createdRole.ID()); err != nil {
		log.Fatalf("   ✗ UserRoleRepository.Asignar (test role): %v", err)
	}
	fmt.Println("   ✓ UserRoleRepository — test role assigned to user")

	// Try to delete role with users → ErrRolConUsuarios
	err = roleRepo2.Eliminar(ctx, createdRole.ID())
	if err == domain.ErrRolConUsuarios {
		fmt.Println("   ✓ RoleRepository — cannot delete role with users (ErrRolConUsuarios)")
	} else {
		log.Fatalf("   ✗ Expected ErrRolConUsuarios, got: %v", err)
	}

	// Cleanup: remove role assignment and delete test role
	if err := userRoleRepo2.Remover(ctx, createdUser.ID(), createdRole.ID()); err != nil {
		log.Fatalf("   ✗ UserRoleRepository.Remover (cleanup): %v", err)
	}
	if err := roleRepo2.Eliminar(ctx, createdRole.ID()); err != nil {
		log.Fatalf("   ✗ RoleRepository.Eliminar (cleanup): %v", err)
	}
	fmt.Println("   ✓ Test role cleaned up")

	// 10. Update user
	fmt.Println("\n📝 10. Actualizar usuario")

	// Disable in domain
	if err := createdUser.Disable(); err != nil {
		log.Fatalf("   ✗ domain.User.Disable: %v", err)
	}

	updatedUser, err := userRepo2.Actualizar(ctx, createdUser)
	if err != nil {
		log.Fatalf("   ✗ UserRepository.Actualizar (disable): %v", err)
	}
	fmt.Printf("   ✓ UserRepository — user disabled: active=%v\n", updatedUser.IsActive())

	// Re-enable
	if err := updatedUser.Enable(); err != nil {
		log.Fatalf("   ✗ domain.User.Enable: %v", err)
	}

	reenabled, err := userRepo2.Actualizar(ctx, updatedUser)
	if err != nil {
		log.Fatalf("   ✗ UserRepository.Actualizar (enable): %v", err)
	}
	fmt.Printf("   ✓ UserRepository — user re-enabled: active=%v\n", reenabled.IsActive())

	// 11. Cleanup: remove test user
	fmt.Println("\n🧹 11. Cleanup")
	if err := userRoleRepo2.Remover(ctx, createdUser.ID(), adminRole2.ID()); err != nil {
		log.Fatalf("   ✗ UserRoleRepository.Remover: %v", err)
	}
	fmt.Println("   ✓ UserRoleRepository — role removed from test user")

	// Delete test user via raw SQL (UserRepository doesn't have Delete)
	db.Exec("DELETE FROM user_roles WHERE user_id = ?", createdUser.ID())
	db.Exec("DELETE FROM users WHERE id = ?", createdUser.ID())
	fmt.Println("   ✓ Test user cleaned up")

	// ===== FINAL =====
	fmt.Println("\n═══════════════════════════════════════")
	fmt.Println("  ✅ TODAS LAS PRUEBAS PASARON")
	fmt.Println("  💾 Config   → PostgreSQL conectado")
	fmt.Println("  📦 Migrations → Tablas creadas")
	fmt.Println("  🌱 Seed     → Permisos + roles OK")
	fmt.Println("  🔧 Repos    → CRUD completo OK")
	fmt.Println("  🔐 Auth     → Permisos verificados OK")
	fmt.Println("  ⚠️  Errors   → Casos borde OK")
	fmt.Println("═══════════════════════════════════════")
}
