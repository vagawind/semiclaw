package service_test

import (
	"context"
	"testing"

	"github.com/vagawind/semiclaw/internal/application/repository"
	"github.com/vagawind/semiclaw/internal/application/service"
	"github.com/vagawind/semiclaw/internal/types"
	"github.com/vagawind/semiclaw/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateTenantCreatesConcreteDefaultStorageBackend(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "local")
	t.Setenv("SEARXNG_DEFAULT_INSTANCE_URL", "")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&types.Tenant{}, &types.StorageBackend{}, &types.WebSearchProviderEntity{}))
	tenantRepo := repository.NewTenantRepository(db)
	storageRepo := repository.NewStorageBackendRepository(db)
	webSearchRepo := repository.NewWebSearchProviderRepository(db)
	tenantSvc := service.NewTenantService(tenantRepo, storageRepo, webSearchRepo)

	tenant, err := tenantSvc.CreateTenant(context.Background(), &types.Tenant{Name: "workspace"})
	require.NoError(t, err)
	require.NotNil(t, tenant.DefaultStorageBackendID)

	backend, err := storageRepo.GetByID(context.Background(), tenant.ID, *tenant.DefaultStorageBackendID)
	require.NoError(t, err)
	require.NotNil(t, backend)
	assert.Equal(t, "local", backend.Provider)
	assert.Equal(t, types.StorageBackendSourceEnv, backend.Source)
	assert.True(t, backend.LegacyAlias)
}

func TestCreateTenantSeedsDefaultSearxngWhenConfigured(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "local")
	t.Setenv("SEARXNG_DEFAULT_INSTANCE_URL", "http://searxng:8080")
	t.Setenv("SSRF_WHITELIST", "")
	t.Setenv("SSRF_WHITELIST_EXTRA", "searxng")
	utils.ResetSSRFWhitelistForTest()
	t.Cleanup(utils.ResetSSRFWhitelistForTest)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&types.Tenant{}, &types.StorageBackend{}, &types.WebSearchProviderEntity{}))
	tenantRepo := repository.NewTenantRepository(db)
	storageRepo := repository.NewStorageBackendRepository(db)
	webSearchRepo := repository.NewWebSearchProviderRepository(db)
	tenantSvc := service.NewTenantService(tenantRepo, storageRepo, webSearchRepo)

	tenant, err := tenantSvc.CreateTenant(context.Background(), &types.Tenant{Name: "workspace"})
	require.NoError(t, err)

	provider, err := webSearchRepo.GetDefault(context.Background(), tenant.ID)
	require.NoError(t, err)
	require.NotNil(t, provider)
	assert.Equal(t, types.WebSearchProviderTypeSearxng, provider.Provider)
	assert.Equal(t, "http://searxng:8080", provider.Parameters.BaseURL)
	assert.True(t, provider.IsDefault)
}
