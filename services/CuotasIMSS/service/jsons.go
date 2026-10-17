package service

import (
	"encoding/json"
	"os"
	"saavedra/services/CuotasIMSS/types"
)

var cesantiaVejez map[string][]*types.CesantiaVejez
var diasVacaciones []*types.DiasVacaciones
var editables *types.Editables
var factores map[string]float64
var primaRiesgo []*types.Riesgo
var Razones []*types.CuotasBloques

func CargarCuotasBloques() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/cuotas_bloques.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &Razones); err != nil {
		return err
	}
	return nil
}

func CargarCesantiaVejez() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/cesantia_vejez.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &cesantiaVejez); err != nil {
		return err
	}
	return nil
}

func CargarDiasVacaciones() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/dias_vacaciones.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &diasVacaciones); err != nil {
		return err
	}
	return nil
}

func CargarEditables() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/editables.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &editables); err != nil {
		return err
	}
	return nil
}

func CargarFactores() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/factores.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &factores); err != nil {
		return err
	}
	return nil
}

func CargarPrimaRiesgo() error {
	datosBytes, err := os.ReadFile("services/CuotasIMSS/config/prima_riesgo.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal(datosBytes, &primaRiesgo); err != nil {
		return err
	}
	return nil
}

func (s service) EscribirJSONEditables(nuevoJSON *types.Editables) error {
	bytes, err := json.MarshalIndent(nuevoJSON, "", " ")
	if err != nil {
		return err
	}
	if err = os.WriteFile("services/CuotasIMSS/config/editables.json", bytes); err != nil {
		return err
	}
	return nil
}
