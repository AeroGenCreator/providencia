// servicioComponente se sustituye por el nombre deseado.
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

    async init() { await this.loadRecords() },
    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async goCuotasIMSS() { const res = await OnlyRedirect("/cuotas"); this.authError = res.authError },
    async clean() {
      this.isSearch = false; this.pageSearch = 1; this.pattern = ''; this.page = 1; return await this.loadRecords()
    },
    block() {
      if (this.isSearch) { return (this.pageSearch === 1 || this.loading); } else { return (this.page === 1 || this.loading); }
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

    async loadRecords() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse(`/cuotas/slice?page=${this.page}`, { method: "GET" })
        this.records = data.records; this.hasNextPage = data.hasNextPage; this.currency();
      } catch (error) {
      throw error
      } finally {
      this.loading = false
      }
    },

  }))
})
