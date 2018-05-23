package mockimpl

import "go-jwt-interface/api"

type UserProviderImpl struct {
	UserList []api.MinimalUser
}

func (provider *UserProviderImpl) Init(){
	provider.UserList = append(provider.UserList,api.MinimalUser{151235,"test","test",[]string{"admin","member"}})
	provider.UserList = append(provider.UserList,api.MinimalUser{919123,"admin","admin",[]string{"sysadmin","member"}})
	provider.UserList = append(provider.UserList,api.MinimalUser{712736213,"foo","bar",[]string{"member"}})
}

func (provider UserProviderImpl) ProvideValidUsers() []api.MinimalUser{
	return provider.UserList
}