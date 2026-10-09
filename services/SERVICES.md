## VISTA LISTA

```html
<!-- PLANTILLA: VISTA LISTA -->
<body x-data="Componente">
<main class="hero">
    <section class="container">
        <div class="columns">
            <div class="column is-1">
                <div class="buttons is-left">
                    <button class="button is-primary is-small is-fullwidth is-rounded" @click="newRecord">
                        <span class="icon"><i class="fa-solid fa-circle-plus"></i></span>
                        <span>Nuevo</span>
                    </button>
                </div>
            </div>
            <div class="column is-3">
                <div class="field">
                    <div class="control is-small has-icons-left">
                        <input class="input is-small is-rounded" type="text" placeholder="Buscador por nombre" x-model="pattern">
                        <span class="icon is-small is-left"><i class="fa-solid fa-magnifying-glass"></i></span>
                    </div>
                </div>
            </div>
            <div class="column is-4">
                <div class="buttons is-left">
                    <button
                        class="button is-info is-small is-rounded"
                        @click="search"
                        :disabled="(pattern === '')"
                    >
                        <span class="icon"><i class="fa-solid fa-magnifying-glass"></i></span>
                        <span>Buscar</span>
                    </button>
                    <button class="button is-danger is-small is-rounded" @click="clean" x-show="(pattern !== '')">
                        <span class="icon"><i class="fa-solid fa-spray-can-sparkles"></i></span>
                        <span>Limpiar</span>
                    </button>
                </div>
            </div>
            <div class="column is-4">
                <div class="buttons is-right">
                    <button
                        class="button is-small is-warning"
                        @click="previousPage"
                        :disabled="block"
                    >
                        <span class="icon"><i class="fa-solid fa-minus"></i></span>
                    </button>
                    <span class="button is-static is-small is-text" x-text="`Página ${page}`"></span>
                    <button class="button is-small is-warning" @click="nextPage" :disabled="(!hasNextPage || loading)">
                        <span class="icon"><i class="fa-solid fa-plus"></i></span>
                    </button>
                </div>
            </div>
        </div>
        <table class="table is-striped is-narrow is-hoverable is-fullwidth">
            <thead>
                <tr>
                    <th>Id</th>
                </tr>
            </thead>
            <tbody>
                <template x-for="" :key="">
                    <tr @click="openRecord()" class="focus">
                        <td x-text=""></td>
                    </tr>
                </template>
                <tr x-show="(records === [])">
                    <td colspan="8">...</td>
                </tr>
            </tbody>
        </table>
    </section>
</main>
```

```js
// COMPONENTE PARA UNA VISTA LISTA BASADO EN LA PLANTILLA ANTERIOR
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
    async goBack() { await GoBack("") },
    async logOut() { await LogOut() },
    async openRecord(id) { await ReadRecord(id, "") },
    async newRecord() { await OnlyRedirect("") },
    async goEmployee() { await OnlyRedirect("") },

    currency() {
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
          `?pattern=${this.pattern}&page=${this.pageSearch}`
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
          `?page=${this.page}`
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
```

## VISTA FORMULARIO NUEVO

```html
<main class="hero">
<section class="container">
    <div class="box">
        <div class="columns">
            <div class="column is-12">
                <div class="buttons">
                    <button class="button is-info is-small" @click="createRecord" :disabled="required">
                        <span class="icon"><i class="fa-solid fa-floppy-disk"></i></span>
                        <span>Crear Registro</span>
                    </button>
                </div>
            </div>
        </div>
        <div class="columns">
            <div class="column is-12" x-show="authError">
                <div class="notification is-danger is-light">
                    <button class="delete" @click="close"></button>
                    <p class="has-text-left is-size-7">
                        <span class="icon-text">
                            <span class="icon"><i class="fa-solid fa-triangle-exclamation"></i></span>
                            <span>...</span>
                        </span>
                    </p>
                </div>
            </div>
        </div>
        <div class="columns">
            <!-- INPUT -->
            <div class="column is-4">
                <label class="label is-small">...</label>
                <div class="control has-icons-left">
                    <input class="input is-small" type="" placeholder="" x-model=""/>
                    <span class="icon is-small is-left"><!-- ICON --></span>
                </div>
            </div>
            <!--CAMPO MONEDA O FLOAT -->
            <div class="column is-4">
                <label class="label is-small">...</label>
                <div class="control has-icons-left">
                    <input 
                        class="input is-small"
                        type="text"
                        placeholder=""
                        x-model="var"
                        @blur="var = currency($event.target.value)"
                    />
                    <span class="icon is-small is-left"><!-- ICON --></span>
                </div>
            </div>
            <!-- CAMPO INTEGER -->
            <div class="column is-4">
                <label class="label is-small">...</label>
                <div class="control has-icons-left">
                    <input
                        class="input is-small"
                        type=""
                        placeholder=""
                        x-model="var"
                        @blur="var = integer($event.target.value)"
                    />
                    <span class="icon is-small is-left"><!-- ICON --></span>
                </div>
            </div>
            <!-- SELECTION -->
            <div class="column is-2">
                <label class="label is-small">...</label>
                <div class="control has-icons-left">
                    <input class="input is-small" list="list-id" id="input-id" placeholder="" x-model="">
                    <span class="icon is-small is-left"><!-- ICON --></span>
                    <datalist id="list-id">
                        <template x-for="" :key="">
                            <option :value="" x-text=""></option>
                        </template>
                    </datalist>
                </div>
            </div>    
        </div>

    </div>
</section>
</main>
```

```js
// COMPONENTE PARA VISTA CREACION REGISTRO FORMULARIO
document.addEventListener('alpine:init', () => {
  Alpine.data("employeesNewComponent", () => ({

    authError: false,
    emailError: false,
    var: '',
    loading: false,

    async init() { await this.loadMany2One() },

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/employee") },
    async logOut() { await LogOut() },
    async goEmployees() { await OnlyRedirect("/employee") },
    close() { this.authError = false; this.emailError = false },

    required() {
      return (this.loading)
    },

    currency(str) {
      const float = ValidateFloatStringToFloat(str)
      if (isNaN(float)) { return '' }
      return FormatterMXN.format(float)
    },

    integers(str) {
      const integer = ValidateIntegerStringToInteger(str)
      if (isNaN(integer)) { return '' }
      return String(integer)
    },

    async loadMany2One() {
      try {
        this.loading = true
        const data = await FetchDataFromResponse("")
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async createRecord() {
      try {
        this.loading = true
        var values = JSON.stringify({field: this.var})
        const res = await CreateRecord("", "", { method: "POST", body: values })
        this.authError = res.authError
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

  }))
})
```
