package mockimpl

import "go-jwt-interface/api"

type UserProviderImpl struct {
	UserList []api.MinimalUser
}

func (provider *UserProviderImpl) Init(){
	provider.UserList = append(provider.UserList,api.GenerateMinimalUser(151235,"test","test"))
	provider.UserList = append(provider.UserList,api.GenerateMinimalUser(919123,"admin","admin"))
	provider.UserList = append(provider.UserList,api.GenerateMinimalUser(712736213,"foo","bar"))
}

func (provider UserProviderImpl) ProvideValidUsers() []api.MinimalUser{
	return provider.UserList
}