document.addEventListener('alpine:init', () => {
  Alpine.data('loginComponent', () => ({
    email: '',
    password: '',
    status: false,

    unlock() { return (this.email === '' || this.password === '') },
    close(){this.status = false},
    goDocs() { window.location.href = "/docs" },
    goObtain() { window.location.href =  "/obtain" },
    goPhil() { window.location.href = "/philosophy" },
    goTutorial() { window.location.href = "/tutorial" },

    async login() {
      const credentials = { "email": this.email, "password": this.password }
      const failure = Login(credentials)
      if (failure) { this.status = failure }
      return
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data('homeComponent', () => ({

    authError: false,

    async goEmployees() {
      var res = await OnlyRedirect("/employee");
      this.authError = res.authError
    },
    async goUsers() {
      var res = await OnlyRedirect("/users");
      this.authError = res.authError
    },

    async goHome() { await GoHome() },
    async goBack() { await GoBack() },
    async logOut() { await LogOut() },
    async close() { this.authError = false },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data('usersListComponent', () => ({

    records: [],
    page: 1,
    searchPage: 1,
    hasNextPage: false,
    pattern: '',
    lastPattern: '',
    loading: false,
    isSearch: false,

    async init() { await this.loadRecords() },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async goUsers() { await OnlyRedirect("/users") },
    async newRecord() { await OnlyRedirect("/users/new") },
    async openRecord(id) { await ReadRecord(id, "/users/record") },
    clean() { return (this.pattern !== '') },

    block() {
      if (!this.isSearch) { return (this.page === 1 || this.loading) };
      if (this.isSearch) { return (this.searchPage === 1 || this.loading) };
    },

    async previousPage() {
      if (!this.isSearch) {
        if (this.page > 1) {
          this.page -= 1;
          return await this.loadRecords()
        };
      }
      if (this.searchPage) {
        if (this.searchPage > 1) {
          this.searchPage -= 1;
          return await this.searchRecord()
        }
      }
    },

    async nextPage() {
      if (!this.isSearch) { if (this.hasNextPage) { this.page += 1; return await this.loadRecords() } }
      if (this.isSearch) { if (this.hasNextPage) { this.searchPage += 1; return await this.searchRecord() } }
    },

    async loadRecords() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse(`/users/slice?page=${this.page}`)
        this.records = data.records
        this.hasNextPage = data.hasNextPage
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async cleanBar() {
      this.isSearch = false;
      this.searchPage = 1;
      this.pattern = '';
      this.page = 1;
      await this.loadRecords()
    },

    async searchRecord() {
      try {
        this.loading = true
        this.isSearch = true
        if (this.lastPattern !== '' && this.lastPattern !== this.pattern) { this.searchPage = 1};
        const data = await FetchDataFromResponse(`/users/search?pattern=${this.pattern}&page=${this.searchPage}`)
        this.lastPattern = this.pattern
        this.records = data.records
        this.hasNextPage = data.hasNextPage
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data('usersNewComponent', () => ({

    parties: [],
    loading: false,
    name: '',
    email: '',
    password: '',
    repeatPassword: '',
    party: '',
    message: false,

    async init() { await this.loadMany2One() },

    close() { this.message = false },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/users") },
    async logOut() { await LogOut() },
    async goUsers() { await OnlyRedirect("/users") },
    required() {
      return (this.name === '' || this.email === '' || this.password === '' || this.repeatPassword === '' || this.party === '')
    },

    async loadMany2One() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse("/users/many2one")
        this.parties = data.parties
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async createRecord() {
      if (this.password !== this.repeatPassword) {
        this.message = true
        return
      }
      if (!ValidateEmail(this.email)) {
        this.message = true
        return
      }
      this.message = false
      const values = { "name": this.name, "email": this.email, "password": this.repeatPassword, "party": this.party }
      await CreateRecord("/users/new", "/users", {method: "POST", body: JSON.stringify(values)})
    },

  }))
})

document.addEventListener('alpine:init', () => {
  Alpine.data('usersRecordComponent', (id) => ({

    parties: [],
    loading: false,
    name: '',
    email: '',
    password: '',
    repeatPassword: '',
    party: '',
    message: false,
    invalid: false,

    async init() {
      await this.loadMany2One()
      await this.loadRecord(id)
    },

    close() {
      this.message = false
      this.invalid = false
    },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/users") },
    async logOut() { await LogOut() },
    async goUsers() { await OnlyRedirect("/users") },

    required() {
      return (this.name === '' || this.email === '' || this.party === '')
    },

    async loadMany2One() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse("/users/many2one")
        this.parties = data.parties
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async loadRecord(id) {
      try {
        this.loading = true
        const record = JSON.stringify({"id": id})
        const data = await FetchDataFromResponse("/users/record", { method: "POST", body: record})
        this.name = data.name
        this.email = data.email
        this.party = data.party
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateRecord() {
      try {
        this.loading = true
        this.message = false
        var checkPass = ''
        if (this.password && this.repeatPassword) {
          if (this.password !== this.repeatPassword) {
            return this.message = true
          } else {
            checkPass = this.password
          }
        }
        if (!ValidateEmail(this.email)) {
          return this.message = true
        }
        const record = JSON.stringify(
          { "id": id, "name": this.name, "email": this.email, "password": checkPass, "party": this.party }
        )
        await UpdateRecord("/users/record", "/users", {method: "PUT", body: record})
      } catch (error) {
        throw error
      } finally {
       this.loading = false
      }
    },

    async deleteRecord() {
      try {
        const record = JSON.stringify({ "id": id })
        const res = await DeleteRecord("/users/record", "/users", { method: "DELETE", body: record })
        this.invalid = res.authError
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})
