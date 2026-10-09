document.addEventListener('alpine:init', () => {
  Alpine.data("employeesListComponent", () => ({

    page: 1,
    pageSearch: 1,
    records: [],
    loading: false,
    hasNextPage: false,
    pattern: '',
    patternLast: '',
    isSearch: false,

    async init() { this.loadRecords() },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async openRecord(id) { await ReadRecord(id, "/employee/record") },
    async newRecord() { await OnlyRedirect("/employee/new") },
    async goEmployee() { await OnlyRedirect("/employee") },

    currency() {
      if (!this.records) { return }
      for (const item of this.records) {
        const money = FormatterMXN.format(item.dailyPayment)
        item.dailyPayment = money
      }
    },

    async clean() {
      this.isSearch = false; this.pageSearch = 1; this.pattern = ''; this.page = 1; return await this.loadRecords()
    },

    async search() {
      try {
        this.loading = true
        this.isSearch = true
        if (this.patternLast !== '' && this.patternLast !== this.pattern) { this.pageSearch = 1};
        const data = await FetchDataFromResponse(
          `/employee/search?pattern=${this.pattern}&page=${this.pageSearch}`
        )
        this.patternLast = this.pattern; this.records = data.records; this.hasNextPage = data.hasNextPage;
        this.currency();
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async loadRecords() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse(
          `/employee/slice?page=${this.page}`
        )
        this.records = data.records;
        this.hasNextPage = data.hasNextPage;
        this.currency();
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    block() {
      if (!this.isSearch) { return (this.page === 1 || this.loading) };
      if (this.isSearch) { return (this.pageSearch === 1 || this.loading) };
    },

    async previousPage() {
      if (!this.isSearch) {
        if (this.page > 1) {
          this.page -= 1;
          return await this.loadRecords()
        };
      }
      if (this.pageSearch) {
        if (this.pageSearch > 1) {
          this.pageSearch -= 1;
          return await this.search()
        }
      }
    },

    async nextPage() {
      if (!this.isSearch) { if (this.hasNextPage) { this.page += 1; return await this.loadRecords() } }
      if (this.isSearch) { if (this.hasNextPage) { this.pageSearch += 1; return await this.search() } }
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data("employeesNewComponent", () => ({

    authError: false,
    emailError: false,
    name: '',
    hireDate: '',
    dailyPayment: '',
    phone: '',
    email: '',
    nss: '',
    curp: '',
    loading: false,

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/employee") },
    async logOut() { await LogOut() },
    async goEmployees() { await OnlyRedirect("/employee") },
    close() { this.authError = false; this.emailError = false },

    required() {
      return (this.name === '' || this.hireDate === '' || this.dailyPayment === '' || this.loading)
    },

    currency(str) {
      const float = ValidateFloatStringToFloat(str)
      if (isNaN(float)) { return '' }
      const money = FormatterMXN.format(float);
      if (money === '$0.00') { return ''}
      return money
    },

    integers(str) {
      const integer = ValidateIntegerStringToInteger(str)
      if (isNaN(integer)) { return '' }
      return String(integer)
    },

    async createRecord() {
      try {
        this.loading = true
        if (!ValidateEmail(this.email)) { this.emailError = true; return }
        const payment = ValidateFloatStringToFloat(this.dailyPayment)
        const number = ValidateIntegerStringToInteger(this.phone)
        var values = JSON.stringify({
          name: this.name,
          hireDate: this.hireDate,
          dailyPayment: payment,
          phone: number,
          email: this.email,
          nss: this.nss,
          curp: this.curp
        })
        const res = await CreateRecord("/employee/new", "/employee", { method: "POST", body: values })
        this.authError = res.authError
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data("employeesRecordComponent", (id) => ({

    loading: false,
    authError: false,
    emailError: false,
    name: '',
    hireDate: '',
    dailyPayment: '',
    phone: '',
    email: '',
    nss: '',
    curp: '',

    async init() { await this.loadRecord(id) },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/employee") },
    async logOut() { await LogOut() },
    async goEmployees() { await OnlyRedirect("/employee") },

    close() { this.authError = false; this.emailError = false; },
    required() {
      return (this.name === '' || this.hireDate === '' || this.dailyPayment === '' || this.loading)
    },

    async loadRecord(id) {
      try {
        this.loading = true
        const data = await FetchDataFromResponse(
          "/employee/record", { method: "POST", body: JSON.stringify({ id: id }) }
        )
        this.name = data.name
        this.hireDate = data.hireDate
        this.dailyPayment = FormatterMXN.format(data.dailyPayment)
        this.phone = String(data.phone)
        this.email = data.email
        this.nss = data.nss
        this.curp = data.curp
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    currency(str) {
      const float = ValidateFloatStringToFloat(str)
      if (isNaN(float)) { return '' }
      const money = FormatterMXN.format(float);
      if (money === '$0.00') { return ''}
      return money
    },

    integers(str) {
      const integer = ValidateIntegerStringToInteger(str)
      if (isNaN(integer)) { return '' }
      return String(integer)
    },

    async updateRecord() {
      try {
        this.loading = true
        if (!ValidateEmail(this.email)) { this.emailError = true; return }
        const payment = ValidateFloatStringToFloat(this.dailyPayment)
        const number = ValidateIntegerStringToInteger(this.phone)
        var values = JSON.stringify({
          id: id,
          name: this.name,
          hireDate: this.hireDate,
          dailyPayment: payment,
          phone: number,
          email: this.email,
          nss: this.nss,
          curp: this.curp
        })
        const res = await UpdateRecord("/employee/record", "/employee", { method: "PUT", body: values })
        this.authError = res.authError
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteRecord() {
      try {
        this.loading = true
        const values = JSON.stringify({id: id})
        const res = await DeleteRecord("/employee/record", "/employee", { method: "DELETE", body: values })
        this.authError = res.authError
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})
