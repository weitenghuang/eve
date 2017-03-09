package rohr

type Authorizable interface {
	AuthorizedRead(user *User) bool
	AuthorizedWrite(user *User) bool
	AuthorizedExecute(user *User) bool
	BindAuthorization(auth Authorization)
}

func (auth *Authorization) DefaultAuthorizedRead(user *User) bool {
	if user.Id == AGENT_USER {
		return true
	}
	return (auth.Owner == user.Id)
}

func (auth *Authorization) DefaultAuthorizedWrite(user *User) bool {
	if user.Id == AGENT_USER {
		return true
	}
	return (auth.Owner == user.Id)
}

func (auth *Authorization) DefaultAuthorizedExecute(user *User) bool {
	if user.Id == AGENT_USER {
		return true
	}
	return (auth.Owner == user.Id)
}

func (infra *Infrastructure) AuthorizedRead(user *User) bool {
	return infra.Authorization.DefaultAuthorizedRead(user)
}

func (infra *Infrastructure) AuthorizedWrite(user *User) bool {
	return infra.Authorization.DefaultAuthorizedWrite(user)
}

func (infra *Infrastructure) AuthorizedExecute(user *User) bool {
	return infra.Authorization.DefaultAuthorizedExecute(user)
}

func (infra *Infrastructure) BindAuthorization(auth Authorization) {
	infra.Authorization = auth
}

func (quoin *Quoin) AuthorizedRead(user *User) bool {
	return quoin.Authorization.DefaultAuthorizedRead(user)
}

func (quoin *Quoin) AuthorizedWrite(user *User) bool {
	return quoin.Authorization.DefaultAuthorizedWrite(user)
}

func (quoin *Quoin) AuthorizedExecute(user *User) bool {
	return quoin.Authorization.DefaultAuthorizedExecute(user)
}

func (quoin *Quoin) BindAuthorization(auth Authorization) {
	quoin.Authorization = auth
}

func (quoinArchive *QuoinArchive) AuthorizedRead(user *User) bool {
	return quoinArchive.Authorization.DefaultAuthorizedRead(user)
}

func (quoinArchive *QuoinArchive) AuthorizedWrite(user *User) bool {
	return quoinArchive.Authorization.DefaultAuthorizedWrite(user)
}

func (quoinArchive *QuoinArchive) AuthorizedExecute(user *User) bool {
	return quoinArchive.Authorization.DefaultAuthorizedExecute(user)
}

func (quoinArchive *QuoinArchive) BindAuthorization(auth Authorization) {
	quoinArchive.Authorization = auth
}

func (provider *Provider) AuthorizedRead(user *User) bool {
	return true
}

func (provider *Provider) AuthorizedWrite(user *User) bool {
	return true
}

func (provider *Provider) AuthorizedExecute(user *User) bool {
	return true
}

func (provider *Provider) BindAuthorization(auth Authorization) {
	provider.Authorization = auth
}
