package repo

import (
	"com.banxiaoxiao.server/repository"
)

var rep repository.Repository

func SetRepository(r repository.Repository) {
	rep = r
}

func GetRepository() repository.Repository {
	return rep
}
