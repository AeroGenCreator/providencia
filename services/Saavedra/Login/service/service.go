package service

import (
	"saavedra/services/Saavedra/Login/store"
	loginTypes "saavedra/services/Saavedra/Login/types"
	"saavedra/services/Saavedra/Session/types"
	userTypes "saavedra/services/Saavedra/Users/types"
	"saavedra/utils"
)

type Service interface {
	Login(user *userTypes.UserStr) (*loginTypes.Credentials, error)
}

type service struct {
	store store.Store
}

func New(store store.Store) Service {
	return service{store: store}
}

func (s service) Login(user *userTypes.UserStr) (*loginTypes.Credentials, error) {

	// 1. Obtención del usuario desde la base de datos.
	dbUser, err := s.store.ReadUser(user)
	if err != nil {
		return nil, err
	}

	// 2. Validación de contraseña.
	valid := utils.CheckPasswordHash(user.Password, dbUser.Password)

	// 3. Generar JWT Token.
	var credentials loginTypes.Credentials
	if !valid {
		credentials.Token = ""
		return &credentials, nil
	}
	token, err := utils.GenerateToken(dbUser.Id, dbUser.Name, dbUser.Party)
	if err != nil {
		return nil, err
	}
	// 4. Sí y solo sí se generarón credenciasles entonces se registra el token como sessión de usuario.
	session := types.Session{UserId: dbUser.Id, Token: token}
	if err = s.store.CreateSession(&session); err != nil {
		return nil, err
	}
	credentials.Token = token
	return &credentials, nil
}
