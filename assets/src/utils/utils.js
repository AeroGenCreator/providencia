// Variables temporales
let refreshSubscribers = [];
let isRefreshing = false;

// === CONSTANTES GLOBALES ===

const FormatterMXN = new Intl.NumberFormat('es-MX', {
  style: 'currency',
  currency: 'MXN',
});

// === Funciones Simples ===

function ValidateEmail(email) {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return regex.test(email);
}

function ValidateFloatStringToFloat(str) {
  if (!str) return 0.0;
  const regex = /[^0-9.-]+/g;
  return parseFloat(str.replace(regex, ""))
}

function ValidateIntegerStringToInteger(str) {
  if (!str) return 0;
  const float = ValidateFloatStringToFloat(str)
  return Math.trunc(float)
}

function onTokenRefreshed() {
  // Si existen peticiones en cola, se ejecutn
  refreshSubscribers.forEach((callback) => callback());
  refreshSubscribers = [];
}

// Función que solicita refrescar token.
async function SecureFetching(route, requestContent = {}, customHeaders = {'X-Requested-With': 'jsFrontendComponent'}) {
  console.log(`Attempting secure fetch for ${route}...`)

  const options = {
    ...requestContent,
    credentials: 'include',
    headers: {
      'X-Requested-With': 'jsFrontendComponent',
      ...(requestContent.headers || {}),
      ...customHeaders
    }
  };

  try {
    let response = await fetch(route, options);
    if (response.status === 401) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          refreshSubscribers.push(async () => {
            resolve(await fetch(route, options));
          });
        });
      }

      isRefreshing = true;
      console.log("Attempting refresh...")
      const refreshResponse = await fetch("/session", { method: "POST", credentials: 'include' });
      isRefreshing = false;

      if (refreshResponse.ok) {
        onTokenRefreshed();
        return await fetch(route, options);
      } else {
        return refreshResponse;
      }
    }
    return response;
  } catch (error) {
    console.error("Reusable fetching component error:", error);
    throw error;
  }
}

// === Utilizan Secure Fetching ===

async function GoHome() {
  const res = await this.SecureFetching('/home', { method: 'HEAD' })
  if (res.ok) {
    window.location.href = '/home'
    return
  }
  await this.LogOut()
  alert(res.status === 401 ? 'Sesión expirada' : `Error ${res.status}`)
}

async function GoBack(redirect = "/home") {
  try {
    const res = await SecureFetching("/home", { method: "HEAD" })
    if (!res.ok) {
      throw new Error(res.status)
    }
    window.location.href = redirect
  } catch (error) {
    throw error
  }
}

async function GoNew(redirect) {
  try {
    const res = await SecureFetching("/home", { method: "HEAD" })
    if (!res.ok) {
      throw new Error(res.status)
    }
    window.location.href = redirect
  } catch (error) {
    throw error
  }
}

async function OnlyRedirect(path) {
  try {
    var object = { authError: false };
    const res = await SecureFetching("/home", { method: "HEAD" })
    if (!res.ok) {
      if (res.status === 401) {
        window.alert("Usuario no autorizado / Unauthorized user.")
        return object
      } else {
        throw new Error(res.status)
      }
    }
    const direction = await SecureFetching(path, { method: "GET" })
    if (!direction.ok) {
      if (direction.status === 401) {
        object.authError = true
        return object
      } else {
        throw new Error(direction.status)
      }
    }
    window.location.href = path
    return object
  } catch (error) {
    throw error
  }
}

async function ReadRecord(id, redirect) {
  try {
    var object = { authError: false };
    const res = await SecureFetching("/home", { method: "HEAD" })
    if (!res.ok) {
      if (res.status === 401) {
        object.authError = true;
        return object
      } else {
        throw new Error(res.status)
      }
    }
    window.location.href = `${redirect}?id=${id}`
    return object
  } catch (error) {
    throw error
  }
}

async function CreateRecord(path, redirect, options = {}) {
  try {
    var object = { authError: false };
    const res = await SecureFetching(path, options)
    if (!res.ok) {
      if (res.status === 401) {
        object.authError = true;
        return object
      } else if (res.status === 409) {
        const duplicated = await res.text()
        alert(duplicated)
        return
      } else {
        throw new Error(res.status)
      }
    }
    window.location.href = redirect
    return object
  } catch (error) {
    throw error
  }
}

async function UpdateRecord(path, redirect, options = {}) {
  try {
    var object = {authError: false };
    const res = await SecureFetching(path, options)
    if (!res.ok) {
      if (res.status === 401) {
        object.authError = true;
        return object
      } else {
        throw new Error(await res.text())
      }
    }
    window.location.href = redirect
    return object
  } catch (error) {
    throw error
  }
}

async function DeleteRecord(path, redirect, options = {}) {
  try {
    var object = {authError: false };
    const res = await SecureFetching(path, options)
    if (!res.ok) {
      if (res.status === 401) {
        object.authError = true;
        return object
      } else {
        throw new Error(await res.text())
      }
    }
    window.location.href = redirect
    return object
  } catch (error) {
    throw error
  }
}

async function FetchDataFromResponse(path, options = {}) {
  try {
    const res = await SecureFetching(path, options)
    if (!res.ok) throw new Error(`Error ${res.status}`)
    const data = await res.json()
    return data
  } catch (error) { throw error }
}

// === No utilizan secure Fetching ===

async function Login(values) {
  try {
    const head = { "Content-Type": "application/json" }
    const res = await fetch("/", { method: "POST", body: JSON.stringify(values), credentials: "include", headers: head })
    if (!res.ok) {
      if (res.status === 401) {
        return true
      } else {
        throw new Error(res.status)
      }
    }
    window.location.href = "/home"
  } catch (error) {
    throw error
  }
}

async function LogOut() {
  try {
    const res = await fetch('/session', {
      method: 'PUT',
      credentials: 'include',
      headers: { 'X-Requested-With': 'jsFrontendComponent' },
    })
    if (res.ok) window.location.href = '/'
  } catch (error) {
    console.error('No fue posible cerrar sesión:', error)
  }
}
