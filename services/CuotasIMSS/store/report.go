package store

func ReportQuery(year int) {
	q := `
	SELECT
    year,
    month,
    name AS empleado,
    -- 1. Subtotal de cuotas del Obrero (sumando todas las que tienen sufijo _obrero)
    (
        SUM(enfermedad_maternidad_adicional_obrero) +
        SUM(prestaciones_dinero_obrero) +
        SUM(gastos_medicos_pensionados_obrero) +
        SUM(invalidez_y_vida_obrero) +
        SUM(cesantia_y_vejez_obrero)
    ) AS subtotal_cuota_obrero,
    -- 2. Subtotal de cuotas del Patrón (sumando las demás prestaciones)
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
    -- 3. Total general (la suma de ambos subtotales)
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
}
