import requests


def obtener_usuarios(llave: str) -> tuple[list, str]:
    url: str = f"http://localhost:8080/ver/usuarios?nombre={llave}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return [], response.text

        data: list = response.json()

        return data, ""
    except Exception as err:
        return [], f"No se puco conectar al servidor: {err}"


def obtener_acuerdos() -> tuple[list, str]:
    url: str = "http://localhost:8080/ver/acuerdo"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return [], response.text

        data: list = response.json()

        return data, ""
    except Exception as err:
        return [], f"No se puco conectar al servidor: {err}"


def obtener_prestamo(codigo: int) -> tuple[dict, str]:
    url: str = f"http://localhost:8080/ver/prestamo/{codigo}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        data: dict = response.json()

        return data, ""
    except Exception as err:
        return {}, f"No se puco conectar al servidor: {err}"
