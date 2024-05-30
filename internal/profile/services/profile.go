package services

import (
	"context"
	"github.com/krls256/dsd2024additional/internal/profile/entities"
	"github.com/krls256/dsd2024additional/internal/profile/errs"
	pkgEntities "github.com/krls256/dsd2024additional/pkg/entities"
	pkgRepositories "github.com/krls256/dsd2024additional/pkg/repositories/pgsql"
)

func NewProfileService(profileRepository *pkgRepositories.BaseRepository[*entities.Profile]) *ProfileService {
	return &ProfileService{profileRepository: profileRepository}
}

type ProfileService struct {
	profileRepository *pkgRepositories.BaseRepository[*entities.Profile]
}

func (p *ProfileService) Upsert(ctx context.Context, req entities.UpsertProfileRequest) (*entities.Profile, error) {
	profile, ok, err := p.profileRepository.FindBy(ctx, map[string]interface{}{
		"id": req.ID,
	}, pkgEntities.Order{}.SetColumn("created_at").SetDesc().ToSlice())

	if err != nil {
		return nil, err
	}

	if ok {
		profile.SetUpsertProfileRequest(req)
	} else {
		profile = entities.NewProfile(req)
	}

	return profile, p.profileRepository.SaveNoReturn(ctx, profile)
}

func (p *ProfileService) Get(ctx context.Context, req entities.GetProfileRequest) (*entities.Profile, error) {
	profile, ok, err := p.profileRepository.FindBy(ctx, map[string]interface{}{
		"id": req.ID,
	}, pkgEntities.Order{}.SetColumn("created_at").SetDesc().ToSlice())

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errs.ProfileNotExists
	}

	return profile, nil
}
