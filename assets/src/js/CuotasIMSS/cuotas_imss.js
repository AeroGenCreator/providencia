// COMPONENTE LISTA "CUOTAS IMSS"
document.addEventListener('alpine:init', () => {
  Alpine.data("CuotasIMSSListaComponent", () => ({

    page: 1,
    pageSearch: 1,
    records: [],
    loading: false,
    hasNextPage: false,
    pattern: '',
    patternLast: '',
    isSearch: false,
    authError: false,

    async init() { await this.loadRecords() },
    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async goCuotasIMSS() { const res = await OnlyRedirect("/cuotas"); this.authError = res.authError },
    async newRecord() { const res = await OnlyRedirect("/cuotas/new"); this.authError = res.authError },
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
      if (!this.records) { return }
      for (const item of this.records) {
        const moneyDailyPayment = FormatterMXN.format(item.daily_payment)
        const moneySBC = FormatterMXN.format(item.sbc)
        item.daily_payment = moneyDailyPayment
        item.sbc = moneySBC
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

    authError: false,
    loading: false,
    employee: [],
    months: [],
    selectedEmployee: "",
    selectedMonth: "",
    selectedYear: "",

    async init() { await this.loadMany2One() },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/cuotas") },
    async logOut() { await LogOut() },
    async goCuotasIMSSNuevo() { const res = await OnlyRedirect("/cuotas/new"); this.authError = res.authError; },

    async loadMany2One() {
      try {
        const data = await FetchDataFromResponse("/cuotas/many2one");
        this.employee = data.employee; this.months = data.months
      } catch (error) { throw error } finally { this.loading = false }
    },

    integer(str) {
      const integers = ValidateIntegerStringToInteger(str)
      if (isNaN(integers)) { return '' }
      return String(integers).substring(0, 4);
    },

    requiredCreate() { },
    requiredCalculate() { },

  }))
})
