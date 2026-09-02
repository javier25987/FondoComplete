import requests


def obtener_datos_user(idx: int) -> tuple[dict, str]:
    url: str = f"http://localhost:8080/rifas/{idx}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        datos: dict = response.json()

        return datos, ""
    except Exception as err:
        return {}, f"No se pudo conectar con el servidor: {err}"


def entregar_boletas(idx: int, datos: dict) -> str:
    url: str = f"http://localhost:8080/rifas/{idx}"

    try:
        response: requests.Response = requests.post(url, json=datos)

        if not response.ok:
            return response.text

        return ""
        
    except Exception as err:
        return {}, f"No se pudo conectar con el servidor: {err}"


def pagar_deuda_por_boletas(idx: int, datos: dict) -> str:
    url: str = f"http://localhost:8080/rifas/{idx}"

    try:
        response: requests.Response = requests.put(url, json=datos)

        if not response.ok:
            return response.text

        return ""
        
    except Exception as err:
        return {}, f"No se pudo conectar con el servidor: {err}"