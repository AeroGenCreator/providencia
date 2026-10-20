package store

import "saavedra/services/CuotasIMSS/types"

func (s store) ReportQuery(year int) ([]*types.ReportRow, error) {

	rows, err := s.db.Query(SQLstring, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*types.ReportRow
	for rows.Next() {
		var r types.ReportRow
		if err = rows.Scan(&r.Year, &r.Month, &r.Name, &r.SubtotalObrero, &r.SubtotalPatron, &r.Total); err != nil {
			return nil, err
		}
		records = append(records, &r)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return records, nil
}

const SQLstring string = `
SELECT year, month, name,
(
    SUM(enfermedad_maternidad_adicional_obrero) +
    SUM(prestaciones_dinero_obrero) +
    SUM(gastos_medicos_pensionados_obrero) +
    SUM(invalidez_y_vida_obrero) +
    SUM(cesantia_y_vejez_obrero)
) AS subtotal_cuota_obrero,
(
    SUM(riesgo_trabajo_prima) +
    SUM(enfermedad_maternidad_cuota_fija) +
    SUM(enfermedad_maternidad_adicional) +
    SUM(prestaciones_dinero_patron) +
    SUM(gastos_medicos_pensionados_patron) +
    SUM(invalidez_y_vida) +
    SUM(retiro) +
    SUM(cesantia_y_vejez_patron) +
    SUM(guarderias) +
    SUM(aportacion_infonavit)
) AS subtotal_cuota_patron,
(
    SUM(enfermedad_maternidad_adicional_obrero) +
    SUM(prestaciones_dinero_obrero) +
    SUM(gastos_medicos_pensionados_obrero) +
    SUM(invalidez_y_vida_obrero) +
    SUM(cesantia_y_vejez_obrero) +
    SUM(riesgo_trabajo_prima) +
    SUM(enfermedad_maternidad_cuota_fija) +
    SUM(enfermedad_maternidad_adicional) +
    SUM(prestaciones_dinero_patron) +
    SUM(gastos_medicos_pensionados_patron) +
    SUM(invalidez_y_vida) +
    SUM(retiro) +
    SUM(cesantia_y_vejez_patron) +
    SUM(guarderias) +
    SUM(aportacion_infonavit)
) AS total
FROM
  	cuotas_imss
WHERE
  	year = ?
GROUP BY
   year,
   month,
   name
ORDER BY
   year,
   month,
   name;`
