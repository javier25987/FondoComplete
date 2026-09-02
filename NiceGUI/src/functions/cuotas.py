import requests


def obtener_datos_user(idx: int) -> tuple[dict, str]:
    url: str = f"http://localhost:8080/cuotas/{idx}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        datos: dict = response.json()

        return datos, ""
    except Exception as err:
        return {}, f"No se pudo conectar con el servidor: {err}"


def pagar_semanas(idx: int, data: dict) -> tuple[int, str]:
    url: str = f"http://localhost:8080/cuotas/{idx}"

    try:
        response: requests.Response = requests.put(url, json=data)

        if not response.ok:
            return -1, response.text

        dato: dict = response.json()

        return dato["codigo_cheque"], ""
    except Exception as err:
        return -1, f"No se pudo conectar con el servidor: {err}"


def crear_cheque(codigo: int) -> str:
    url: str = f"http://localhost:8080/cuotas/cheque/{codigo}"

    try:
        response: requests.Response = requests.post(url)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar con el servidor: {err}"
