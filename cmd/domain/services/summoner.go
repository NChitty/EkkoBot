package services

import (
	"context"
	"fmt"

	"github.com/NChitty/lol-discord-bot/cmd/domain/models"
)

type RiotAdapter interface {
	GetRankedStats(ctx context.Context, summoner models.Summoner) (models.SummonerStats, error)
	GetSummoner(ctx context.Context, name string, tag string) (models.Summoner, error)
}

const SUMMONER_NOT_FOUND_ERRORF string = "Could not find sumoner %s#%s"

type SummonerRepository interface {
	GetSummoner(ctx context.Context, name string, tag string) (models.Summoner, error)
	SaveSummoner(ctx context.Context, summoner models.SummonerStats) (models.SummonerStats, error)
}

type SummonerService struct {
	riotAdapter        RiotAdapter
	summonerRepository SummonerRepository
}

func NewSummonerService(riotAdapter RiotAdapter, summonerRepository SummonerRepository) *SummonerService {
	return &SummonerService{riotAdapter, summonerRepository}
}

func (s *SummonerService) GetSummonerStats(ctx context.Context, name string, tag string) (models.SummonerStats, error) {
	summoner, err := s.summonerRepository.GetSummoner(ctx, name, tag)
	if err != nil && err.Error() == fmt.Sprintf(SUMMONER_NOT_FOUND_ERRORF, name, tag) {
		return s.CreateSummoner(ctx, name, tag)
	}

	stats, err := s.riotAdapter.GetRankedStats(ctx, summoner)
	if err != nil {
		return models.SummonerStats{}, err
	}

	return s.summonerRepository.SaveSummoner(ctx, stats)
}

func (s *SummonerService) CreateSummoner(ctx context.Context, name string, tag string) (models.SummonerStats, error) {
	if summoner, err := s.riotAdapter.GetSummoner(ctx, name, tag); err != nil {
		return models.SummonerStats{}, err
	} else {
		stats, err := s.riotAdapter.GetRankedStats(ctx, summoner)
		if err != nil {
			return models.SummonerStats{}, err
		}
		return s.summonerRepository.SaveSummoner(ctx, stats)
	}
}
