package common

import (
	v2 "learn_words/datasources/v2"
	"learn_words/datasources/v2/models"
	"slices"
)

type DictionaryService struct {
	ds v2.RWDataSourceV2
}

func (s *DictionaryService) GetGroupNames() ([]string, error) {
	groups, err := s.ds.ReadAllGroups()
	if err != nil {
		return nil, err
	}
	return *groups.AsList().Names(), nil
}

func (s *DictionaryService) GetRandomWords(count int, groupIds []int64) (models.WordList, error) {
	groups, err := s.ds.ReadAllGroups()
	if err != nil {
		return nil, err
	}
	var wordIds []int64
	for _, id := range groupIds {
		wordIds = append(wordIds, (*groups)[id].Words...)
	}
	slices.Sort(wordIds)
	slices.Compact(wordIds)
	wordIds = slices.Clip(wordIds)

	words, err := s.ds.ReadWords(wordIds)
	if err != nil {
		return nil, err
	}
	result := words.Shuffle(count)
	return result, nil
}

func (s *DictionaryService) IncrementStatValue(wordId int64, key int) error {
	stat, err := s.ds.LoadStat(wordId)
	if err != nil {
		return err
	}
	stat.Statistic[key] += 1
	err = s.ds.UpdateStat(stat)
	if err != nil {
		return err
	}
	return nil
}

func NewDictionaryService(ds v2.RWDataSourceV2) *DictionaryService {
	return &DictionaryService{ds: ds}
}
