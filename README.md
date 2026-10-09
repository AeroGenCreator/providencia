![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=black)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=white&labelColor=black)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white&labelColor=black)
![CSS](https://img.shields.io/badge/CSS-CSS?style=for-the-badge&logo=CSS&logoColor=white&labelColor=black&color=purple)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white&labelColor=black)
![Alpine.js](https://img.shields.io/badge/Alpine.js-8BC0D0?style=for-the-badge&logo=alpine.js&logoColor=white&labelColor=black)
![JWT](https://img.shields.io/badge/jwt-jwt?style=for-the-badge&logo=jsonwebtokens&logoColor=white&labelColor=black&color=yellow)
![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white&labelColor=black)
![FastAPI](https://img.shields.io/badge/FastAPI-009688?style=for-the-badge&logo=fastapi&logoColor=white&labelColor=black)
![PyPI](https://img.shields.io/badge/PyPI-3775A9?style=for-the-badge&logo=pypi&logoColor=white&labelColor=black)
![Pytest](https://img.shields.io/badge/Pytest-0A9EDC?style=for-the-badge&logo=pytest&logoColor=white&labelColor=black)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black&labelColor=white)
![macOS](https://img.shields.io/badge/macos-macos?style=for-the-badge&logo=apple&logoColor=white&labelColor=black&color=pink)
![Bulma](https://img.shields.io/badge/bulma-bulma?style=for-the-badge&logo=bulma&logoColor=white&labelColor=black&color=brown)

# 🖌️ Saavedra

Saavedra es una caja de herramientas unida a través de un contrato de flujos.

`Una arquitectura pensada para el código explicito`.

- Saavedra fue pensado para el desarrollo de servicios empresariales; asimismo, se tomó inspiración de los módulos de los ERPs modernos.

Sin embargo, cada servicio se construye desde cero, lo cual facilita la personalización por cliente.

## 🚀 Inicio Rápido

En lugar de clonar, se recomienda descargar el `.zip`. De esta manera, cada nuevo proyecto podrá ser respaldado en un repositorio distinto sin necesidad de ser rechazado por el repositorio original.

Una vez tengas una copia de Saavedra, cualquier código extra podrá ser agregado o eliminado sin problema.

Sin embargo, es recomendable seguir los siguientes pasos:

- [Tutorial Básico - Construye un CRUD]()
- [Filosofía - ¿Por qué Saavedra?]()
- [Documentación - Fragmentos y Piezas]()

Es importante resaltar que el tutorial es corto debido a que Saavedra no es un `motor` de desarrollo; en su lugar, es una serie de contratos de flujo y fragmentos de código reutilizables fáciles de implementar.

_A continuación, el esquema de capas usado por Saavedra:_

```txt
Frontend -> Endpoint -> Servicio -> Store -> Database
```

## 🏛️ Acerca de

### 📦 Backend

El backend de Saavedra se inspira en `Repository Pattern`. Esto es fundamental para construir CRUDs explícitos y desacoplados sin depender de Frameworks u ORMs.

Una explicación más profunda de este patrón de desarrollo se encuentra en: [Saavedra Filosofía]()

### 🎨 Frontend

El frontend hace peticiones gracias a `Alpine.js`. Las peticiones viajan de manera segura gracias a funciones `JavaScript` reutilizables.

- Secure Fetching JWT
- Redirección Automática
- Peticiones CRUD preconstruidas.

### 📖 Principios

- Cada servicio resuelve una necesidad de negocio.
- La data auditable se almacena en la base de datos, mientras que la de configuración se guarda en JSON.
- Codificación explícita.
- Asistencia de IA permitida gracias a la separación de módulos, servicios y capas.

#### 🐳 Bienvenido

Accede al [Tutorial]() para desarrollar un flujo usando Go & JavaScript.

```Go
type Bienvenido struct {
	Tutorial  string `json:"tutorial"`
	Proyectos []int  `json:"proyectos"`
}

func Construir() {
	servicios := []int{1, 2, 3, 4, 5}
	proyectos := Bienvenido{
		Proyectos: servicios,
	}
}
```
