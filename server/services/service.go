package services

import "ecom/server/repos"

type Service struct {
	IService
	Repo repos.IRepo
}

func NewService(repo repos.IRepo) *Service {
	return &Service{Repo: repo}
}
