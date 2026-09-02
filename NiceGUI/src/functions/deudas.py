import requests


def obtener_deudas(idx: int) -> tuple[dict, str]:
    url: str = f"http://localhost:8080/deudas/{idx}"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        data: dict = response.json()

        return data, ""
    except Exception as err:
        return {}, f"No se puco conectar al servidor: {err}"


def pagar_deuda(idx: int, codigo: int) -> str:
    url: str = f"http://localhost:8080/deudas/{idx}?codigo={codigo}"

    try:
        response = requests.put(url)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se puede conectar al servidor {err}"


def crear_deuda(idx: int, data: dict) -> str:
    url: str = f"http://localhost:8080/deudas/{idx}"

    try:
        response = requests.post(url, json=data)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se puede conectar al servidor {err}"
