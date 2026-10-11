// servicioComponente se sustituye por el nombre deseado.
document.addEventListener('alpine:init', () => {
  Alpine.data("CuotasIMSSListaComponent", () => ({

    records: [],

    async goHome() { await GoHome() },
    async goBack() { await GoBack("/home") },
    async logOut() { await LogOut() },
    async goCuotasIMSS() {
        const res = await OnlyRedirect("/cuotas")
        this.authError = res.authError
    },
  }))
})
