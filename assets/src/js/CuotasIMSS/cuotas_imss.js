// COMPONENTE LISTA "CUOTAS IMSS"
document.addEventListener('alpine:init', () => {
  Alpine.data("CuotasIMSSListaComponent", () => ({

    page: 1,
    pageSearch: 1,
    records: [],
    loading: false,
    hasNextPage: false,
    pattern: '',
    year: '',
    patternLast: '',
    yearLast: '',
    isSearch: false,
    authError: false,

    async init() { await this.loadRecords() },
    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async goCuotasIMSS() { const res = await OnlyRedirect("/cuotas"); this.authError = res.authError },
    async goConfig() { const res = await OnlyRedirect("/cuotas/config"); this.authError = res.authError },
    async goReport() { const res = await OnlyRedirect("/cuotas/report"); this.authError = res.authError },
    async newRecord() { const res = await OnlyRedirect("/cuotas/new"); this.authError = res.authError },
    async openRecord(id) { await ReadRecord(id, "/cuotas/record") },
    close() { this.authError = false; },
    async clean() {
      this.isSearch = false; this.pageSearch = 1; this.pattern = ''; this.page = 1; return await this.loadRecords()
    },
    block() {
      if (this.isSearch) { return (this.pageSearch === 1 || this.loading); } else { return (this.page === 1 || this.loading); }
    },

    async previousPage() {
      if (this.isSearch) {
        if (this.pageSearch > 1) { this.pageSearch -= 1; return await this.search(); }
      } else { if (this.page > 1) { this.page -= 1; return await this.loadRecords() } }
    },

    async nextPage() {
      if (!this.isSearch) { if (this.hasNextPage) { this.page += 1; return await this.loadRecords() } }
      if (this.isSearch) { if (this.hasNextPage) { this.pageSearch += 1; return await this.search() } }
    },

    currency() {
        if (!this.records) { return; }
        for (var item of this.records) {
            // 1. Validamos los valores originales ANTES de formatear
            let daily = Number(item.daily_payment);
            let sbc = Number(item.sbc);

            // 2. Si no son números válidos (NaN, undefined, null, etc.), les ponemos 0
            item.daily_payment = isNaN(daily) ? FormatterMXN.format(0) : FormatterMXN.format(daily);
            item.sbc = isNaN(sbc) ? FormatterMXN.format(0) : FormatterMXN.format(sbc);
        }
    },

    months() {
      if (!this.records) { return }
      this.records.forEach(item => { const monthName = MonthsFromIntegerToString(item.month); item.month = monthName; })
    },

    async loadRecords() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse(`/cuotas/slice?page=${this.page}`, { method: "GET" })
        this.records = data.records; this.hasNextPage = data.hasNextPage; this.currency(); this.months();
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async search() {
      try {
        this.loading = true; this.isSearch = true;
        if (this.patternLast !== '' && this.patternLast !== this.pattern) { this.pageSearch = 1};
        const data = await FetchDataFromResponse(`/cuotas/search?pattern=${this.pattern}&page=${this.pageSearch}`);
        this.patternLast = this.pattern; this.records = data.records; this.hasNextPage = data.hasNextPage; this.currency(); this.months();
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})

// COMPONENTE NUEVO "CUOTAS IMSS"
document.addEventListener('alpine:init', () => {
  Alpine.data("CuotasIMSSNewComponent", () => ({

    badRequestError: false,
    authError: false,
    loading: false,
    employee: [],
    months: [],
    riesgos: [],
    razones: [],
    selectedEmployee: "",
    selectedMonth: "",
    selectedYear: 0,
    selectedReason: "",
    incapacidad: 0,
    faltas: 0,
    uma: 0,
    minPayment: 0,
    riesgo: 0,
    umaFecha: "Fecha de UMA",
    minPaymentFecha: "Fecha de Salario Minimo",
    empId: null,
    antiguedad: 0,
    factorIntegracion: 0,
    sbc: 0,
    dataEmployee: {},
    dataPatron: {},
    dataObrero: {},
    subtotalPatron: 0,
    subtotalObrero: 0,
    total: 0,

    // Carga data de "seleccion" y el año actual.
    async init() {
      await this.loadMany2One();
      const hoy = new Date();
      this.selectedYear = hoy.getFullYear()
    },

    // Botones de Navegación.
    async goHome() { await GoHome() },
    async goBack() { await GoBack("/cuotas") },
    async logOut() { await LogOut() },
    async goCuotasIMSSNuevo() { const res = await OnlyRedirect("/cuotas/new"); this.authError = res.authError; },

    // Cerrar notificaciones de error.
    closebadRequestError() { this.badRequestError = false },
    closeAuthError() { this.authError = false; },

    // Petición: data de "selección".
    async loadMany2One() {
      try {
        this.loading = true;
        const data = await FetchDataFromResponse("/cuotas/many2one");
        this.employee = data.employee; this.months = data.months; this.uma = data.uma; this.minPayment = data.minPayment;
        this.umaFecha = this.umaFecha + ` ${data.umaFecha}`
        this.minPaymentFecha = this.minPaymentFecha + ` ${data.minPaymentFecha}`
        this.riesgos = data.riesgos
        this.razones = data.razones
      } catch (error) { throw error } finally { this.loading = false }
    },

    // Limpia formulario y cuota calculada: Mantiene "Empleado" y "Mes" seleccionados
    clean() {
      this.dataEmployee = {}; this.dataPatron = {}; this.dataObrero = {}; this.antiguedad = 0; this.factorIntegracion = 0;
      this.sbc = 0; this.incapacidad = 0; this.selectedReason = ""; this.faltas = 0;
    },

    // Limpia razon de incapacidad si no se ha seleccionado una cantidad de dias de incapacidad.
    soloSiIncapacidades(str) { if (this.incapacidad <= 0) { return ""; } else { return str; } },

    // Valida el string "Año" recorta a un largo de 4 caracteres.
    integer(str) {
      const integers = ValidateIntegerStringToInteger(str); if (isNaN(integers)) { return 0; }
      return String(integers).substring(0, 4);
    },

    // Conversor de string a float.
    float(str) {
      const float = ValidateFloatStringToFloat(str);
      if (float.isNaN) { return 0; } else return float;
    },

    // Formato de Moneda de 2 decimales MXN
    currency(item) { return FormatterMXN.format(item); },

    // Limita días de Incapacidad según los días del mes seleccionado.
    validarDiasConMes(str) {
      if (this.selectedMonth !== "") {
        const int = ValidateIntegerStringToInteger(str); if (isNaN(int)) { return 0; }
        const mesDias = MonthsDays[this.selectedMonth]
        if (int > mesDias) { return mesDias; } else { return int}
      } else { return null; }
    },

    // Limita las faltas por cuota a 8 (Maximo Legal)
    limitarFaltas(str) {
      if (this.selectedMonth !== "") {
        const int = ValidateIntegerStringToInteger(str); if (isNaN(int)) { return 0; }
        const maximo = 8
        if (int > maximo) { return maximo; } else { return int}
      } else { return null; }
    },

    // Complementa el "bloqueo" de botones para los campos (Incapacidad Y Razon de Incapacidad)
    match() {
      if (this.incapacidad > 0 && this.selectedReason !== "" || this.incapacidad <= 0 && this.selectedReason === "") {
        return true;
      } else { return false }
    },

    // Bloqueo del boton "Create"
    requiredCreate() {
      return (
        this.loading ||
        this.selectedEmployee === "" || this.selectedMonth === "" || this.selectedYear == 0 || this.riesgo == 0 ||
        Object.keys(this.dataEmployee).length === 0 ||
        Object.keys(this.dataPatron).length === 0 ||
        Object.keys(this.dataObrero).length === 0 ||
        this.match() === false
      );
    },

    // Bloqueo del boton "Calcular"
    requiredCalculate() {
      return (
        this.loading ||
        this.selectedEmployee === "" || this.selectedMonth === "" || this.selectedYear == 0 || this.riesgo == 0 ||
        this.match() === false
      );
    },

    // Guarda el ID de Empleado el objeto solo contiene (Id y Nombre)
    getEmployeeId() {
      if (this.employee) {
        this.employee.forEach(r => {
          if (r.name.toLowerCase() === this.selectedEmployee.toLowerCase()) {
            const intId = ValidateIntegerStringToInteger(r.id); if (intId.isNaN) {
              window.alert("No se cargo el id empleado correctamente"); return;
            } else { this.empId = intId; return }
          }
        })
      }
    },

    // Suma los valores de cada par {clave: valor}
    sum(obj) { var total = 0; Object.getOwnPropertyNames(obj).forEach(function (val, idx, array) { total += obj[val] });
      return total
    },

    // Construye el JSON de petición para el backend.
    construirPeticion() {
      const vals = JSON.stringify({
        empId: this.empId,
        month: this.selectedMonth,
        year: String(this.selectedYear),
        incapacidad: this.incapacidad,
        razon: this.selectedReason,
        faltas: this.faltas,
        riesgo: this.riesgo
      })
      return vals
    },

    // Realiza peticion a Backend con calculos de cuotas & extrae la data.
    async calcularCuota() {
      try {
        this.loading = true
        this.getEmployeeId(); if (this.empId === null) { window.alert("No se cargo el empleado correctamente"); return }
        const vals = this.construirPeticion()
        const data = await FetchDataFromResponse("/cuotas/new", { method: "PUT", body: vals })

        // Manejo de errores y Mensajes
        this.authError = data.authError; this.badRequestError = data.badRequestError;
        if (this.authError || this.badRequestError) { return }

        // Calcular "Totales" y "Subtotales" (Formato de moneda)
        this.subtotalPatron = this.sum(data.patron); this.subtotalObrero = this.sum(data.obrero);
        this.total = this.subtotalPatron + this.subtotalObrero;
        this.subtotalPatron = this.currency(this.subtotalPatron); this.subtotalObrero = this.currency(this.subtotalObrero);
        this.total = this.currency(this.total)

        // 1. Datos principales y de empleados
        data.employee.dailyPayment = this.currency(data.employee.dailyPayment);
        data.sbc = this.currency(data.sbc);

        // 2. Campos del Obrero
        data.obrero.enfermedadMaternidadAdicionalObrero = this.currency(data.obrero.enfermedadMaternidadAdicionalObrero);
        data.obrero.prestacionesDineroObrero = this.currency(data.obrero.prestacionesDineroObrero);
        data.obrero.gastosMedicosPensionadosObrero = this.currency(data.obrero.gastosMedicosPensionadosObrero);
        data.obrero.invalidezYVidaObrero = this.currency(data.obrero.invalidezYVidaObrero);
        data.obrero.cesantiaYVejezObrero = this.currency(data.obrero.cesantiaYVejezObrero);

        // 3. Campos del Patrón
        data.patron.riesgoTrabajoPrima = this.currency(data.patron.riesgoTrabajoPrima);
        data.patron.enfermedadMaternidadCuotaFija = this.currency(data.patron.enfermedadMaternidadCuotaFija);
        data.patron.enfermedadMaternidadAdicional = this.currency(data.patron.enfermedadMaternidadAdicional);
        data.patron.prestacionesDineroPatron = this.currency(data.patron.prestacionesDineroPatron);
        data.patron.gastosMedicosPensionadosPatron = this.currency(data.patron.gastosMedicosPensionadosPatron);
        data.patron.invalidezYVida = this.currency(data.patron.invalidezYVida);
        data.patron.retiro = this.currency(data.patron.retiro);
        data.patron.cesantiaYVejezPatron = this.currency(data.patron.cesantiaYVejezPatron);
        data.patron.guarderias = this.currency(data.patron.guarderias);
        data.patron.aportacionInfonavit = this.currency(data.patron.aportacionInfonavit);

        // Asignar data
        this.dataEmployee = data.employee; this.dataObrero = data.obrero; this.dataPatron = data.patron;
        this.factorIntegracion = data.factorIntegracion; this.sbc = data.sbc; this.antiguedad = data.antiguedad;
      } catch (error) { throw error } finally { this.loading = false; this.empId = null }
    },

    // Crea el registro de Cuotas en la DB
    async createRecord() {
      try {
        this.loading = true;
        this.getEmployeeId(); if (this.empId === null) { window.alert("No se cargo el empleado correctamente"); return; }
        const vals = this.construirPeticion()
        const res = await CreateRecord("/cuotas/new", "/cuotas", { method: "POST", body: vals })
        this.authError = res.authError
        this.badRequestError = res.badRequestError
      } catch (error) { throw error } finally { this.loading = false; }
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data("CuotasIMSSRecordComponent", (id) => ({

    authError: false,
    badRequestError: false,
    name: "",
    dailyPayment: 0,
    year: 0,
    month: 0,
    factorIntegracion: 0,
    sbc: 0,
    esBimestre: false,
    diasDelMes: 0,
    diasPorCotizar: 0,
    incapacidad: 0,
    razon: "",
    faltas: 0,
    riesgoTrabajoPrima: 0,
    enfermedadMaternidadCuotaFija: 0,
    enfermedadMaternidadAdicional: 0,
    prestacionesDineroPatron: 0,
    gastosMedicosPensionadosPatron: 0,
    invalidezYVida: 0,
    retiro: 0,
    cesantiaYVejezPatron: 0,
    guarderias: 0,
    aportacionInfonavit: 0,
    enfermedadMaternidadAdicionalObrero: 0,
    prestacionesDineroObrero: 0,
    gastosMedicosPensionadosObrero: 0,
    invalidezYVidaObrero: 0,
    cesantiaYVejezObrero: 0,
    subtotalObrero: 0,
    subtotalPatron: 0,
    total: 0,

    async init() { await this.LoadRecord() },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/cuotas") },
    async logOut() { await LogOut() },
    async goCuotas() { await OnlyRedirect("/cuotas") },
    async goRecord() { await OnlyRedirect("/cuotas/record") },

    // Formato de Moneda de 2 decimales MXN
    currency(item) { return FormatterMXN.format(item); },

    calcularTotales(data) {
        // Acumulador para cuotas patronales (usando snake_case de data)
        let contadorPatron = 0;
        const cuotasPatron = [
            data.riesgo_trabajo_prima,
            data.enfermedad_maternidad_cuota_fija,
            data.enfermedad_maternidad_adicional,
            data.prestaciones_dinero_patron,
            data.gastos_medicos_pensionados_patron,
            data.invalidez_y_vida,
            data.retiro,
            data.cesantia_y_vejez_patron,
            data.guarderias,
            data.aportacion_infonavit
        ];

        for (let i = 0; i < cuotasPatron.length; i++) {
            contadorPatron += Number(cuotasPatron[i]) || 0;
        }

        // Acumulador para cuotas obreras (usando snake_case de data)
        let contadorObrero = 0;
        const cuotasObrero = [
            data.enfermedad_maternidad_adicional_obrero,
            data.prestaciones_dinero_obrero,
            data.gastos_medicos_pensionados_obrero,
            data.invalidez_y_vida_obrero,
            data.cesantia_y_vejez_obrero
        ];

        for (let i = 0; i < cuotasObrero.length; i++) {
            contadorObrero += Number(cuotasObrero[i]) || 0;
        }

        // Asignaciones con formato de moneda
        this.subtotalPatron = this.currency(contadorPatron);
        this.subtotalObrero = this.currency(contadorObrero);
        this.total = this.currency(contadorPatron + contadorObrero);
    },

    async LoadRecord() {
      try {
        const vals = JSON.stringify({ id: id });
        const data = await FetchDataFromResponse("/cuotas/record", { method: "PUT", body: vals });

        // Primero calculamos pasando el objeto crudo del fetch (Presición Decimal)
        this.calcularTotales(data);

        // Luego asignamos cada propiedad formateada al estado
        this.name = data.name;
        this.dailyPayment = this.currency(data.daily_payment);
        this.year = data.year;
        this.month = data.month;
        this.factorIntegracion = data.factor_integracion;
        this.sbc = this.currency(data.sbc);
        this.esBimestre = data.es_bimestre;
        this.diasDelMes = data.diasDelMes;
        this.diasPorCotizar = data.diasPorCotizar;
        this.incapacidad = data.incapacidad;
        this.razon = data.razon;
        this.faltas = data.faltas;
        this.riesgoTrabajoPrima = this.currency(data.riesgo_trabajo_prima);
        this.enfermedadMaternidadCuotaFija = this.currency(data.enfermedad_maternidad_cuota_fija);
        this.enfermedadMaternidadAdicional = this.currency(data.enfermedad_maternidad_adicional);
        this.prestacionesDineroPatron = this.currency(data.prestaciones_dinero_patron);
        this.gastosMedicosPensionadosPatron = this.currency(data.gastos_medicos_pensionados_patron);
        this.invalidezYVida = this.currency(data.invalidez_y_vida);
        this.retiro = this.currency(data.retiro);
        this.cesantiaYVejezPatron = this.currency(data.cesantia_y_vejez_patron);
        this.guarderias = this.currency(data.guarderias);
        this.aportacionInfonavit = this.currency(data.aportacion_infonavit);
        this.enfermedadMaternidadAdicionalObrero = this.currency(data.enfermedad_maternidad_adicional_obrero);
        this.prestacionesDineroObrero = this.currency(data.prestaciones_dinero_obrero);
        this.gastosMedicosPensionadosObrero = this.currency(data.gastos_medicos_pensionados_obrero);
        this.invalidezYVidaObrero = this.currency(data.invalidez_y_vida_obrero);
        this.cesantiaYVejezObrero = this.currency(data.cesantia_y_vejez_obrero);
      } catch (error) { throw error }
    },

    async deleteRecord() {
      try {
        const vals = JSON.stringify({id: id})
        const data = await DeleteRecord("/cuotas/record", "/cuotas", { method: "DELETE", body: vals })
        this.authError = data.authError
      } catch (error) { throw error}
    },

  }))
})
