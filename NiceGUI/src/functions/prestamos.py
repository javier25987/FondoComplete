import requests


def obtener_datos_user(idx: int) -> tuple[dict, str]:
    url: str = f"http://localhost:8080/prestamos/{idx}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        datos: dict = response.json()

        return datos, ""
    except Exception as err:
        return {}, f"No se pudo conectar con el servidor: {err}"


def crear_cheque(codigo: int):
    url: str = f"http://localhost:8080/prestamos/cheque/{codigo}"

    try:
        response: requests.Response = requests.post(url)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar con el servidor: {err}"


def abonar_a_prestamo(idx: int, data: dict) -> str:
    url: str = f"http://localhost:8080/prestamos/{idx}"

    try:
        response: requests.Response = requests.put(url, json=data)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar con el servidor: {err}"


def crear_prestamo(idx: int, data: dict) -> str:
    url: str = f"http://localhost:8080/prestamos/{idx}"

    try:
        response: requests.Response = requests.post(url, json=data)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar con el servidor: {err}"
