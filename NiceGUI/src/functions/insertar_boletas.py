import requests


def obtener_boletas_disponibles() -> list[dict[int, str]]:
    url: str = "http://localhost:8080/insertar"
    
    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return [], response.text

        data: list[dict] = response.json()

        return data, ""
    except Exception as e:
        return [], f"No se pudo conectar al servidor: {e}"


def enviar_boletas_ingresadas(boletas) -> str:
    url: str = "http://localhost:8080/insertar"
    
    try:
        response: requests.Response = requests.post(url, json=boletas)

        if not response.ok:
            return response.text

        return ""
    except Exception as e:
        return f"No se pudo conectar al servidor: {e}"


def eliminar_boleta(codigo) -> str:
    url: str = f"http://localhost:8080/insertar?codigo={codigo}"
    
    try:
        response: requests.Response = requests.delete(url)

        if not response.ok:
            return response.text

        return ""
    except Exception as e:
        return f"No se pudo conectar al servidor: {e}"