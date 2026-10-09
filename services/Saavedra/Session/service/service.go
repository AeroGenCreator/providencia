package service

import (
	"saavedra/services/Saavedra/Session/store"
	"saavedra/utils"
	"strconv"
)

type Service interface {
	RefreshToken(id, name, party, token string) (string, error)
	ExpellUser(id string) error
}

type service struct {
	store store.Store
}

func New(store store.Store) Service {
	return service{store: store}
}

func (s service) RefreshToken(id, name, party, token string) (string, error) {

	// 1. Conversión del id a integer
	intId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	// 2. Query a tabla 'session'.
	session, err := s.store.ReadSession(intId)
	if err != nil {
		return "", err
	}
	// 3. Valida la vigencia del request token con el token de la session en la bd.
	if token != session.Token {
		return "", nil
	}
	// 4. Generar un nuevo token para la sessión del usuario.
	newToken, err := utils.GenerateToken(intId, name, party)
	// 5. Guardar nuevo token para el usuario en la session en la bd.
	if err := s.store.UpdateSession(intId, newToken); err != nil {
		return "", err
	}
	return newToken, nil
}

// Si se llama a expulsion, se cambia el valor del token antiguo por string vacio.
func (s service) ExpellUser(id string) error {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err = s.store.UpdateSession(intId, ""); err != nil {
		return err
	}
	return nil
}
