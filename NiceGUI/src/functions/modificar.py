import requests


def obtener_usuarios() -> tuple[list[dict], str]:
    url: str = "http://localhost:8080/modificar"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return [], response.text

        data: list[dict] = response.json()

        return data, ""
    except Exception as e:
        return [], f"No se pudo conectar al servidor: {e}"


def guardar_usuario(data: dict) -> str:
    url: str = "http://localhost:8080/modificar"

    try:
        response: requests.Response = requests.post(url, json=data)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar al servidor: {err}"

def modificar_fechas_de_un_prestamo(data: dict) -> str:
    url: str = "http://localhost:8080/modificar/fechas"

    try:
        response: requests.Response = requests.put(url, json=data)

        if not response.ok:
            return response.text

        return ""
    except Exception as err:
        return f"No se pudo conectar al servidor: {err}"
