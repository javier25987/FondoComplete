import requests


def obtener_calendario() -> tuple[bool, list[dict[str, str]]]:
    url: str = "http://localhost:8080/ajuste/calendario"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return False, response.text

        datos: dict = response.json()
        return True, datos

    except Exception as e:
        return False, f"No es posible conectarse al servidor: {e}"


def crear_listado_de_fechas(primera_fecha_str: str, dobles: list[str]) -> str:
    url: str = "http://localhost:8080/ajuste/calendario"

    try:
        response: requests.Response = requests.post(
            url,
            json={
                "fecha_inicial": primera_fecha_str,
                "fechas_dobles": dobles,
                "hora_de_cierre": "19:00",
            },
        )

        if not response.ok:
            return response.text

        return ""

    except Exception as e:
        return False, f"No es posible conectarse al servidor: {e}"


# def crear_listado_de_fechas(primera_fecha_str: str, dobles: list[str]) -> str:
#     fechas = []

#     dobles = [i.replace("-", "/") for i in dobles]

#     if "/" in primera_fecha_str:
#         primera_fecha_str.replace("/", "-")

#     primera_fecha = datetime.datetime(*map(int, primera_fecha_str.split("-"))) #type:ignore

#     i: int = 0
#     while len(fechas) < 50:
#         new_fecha_date = primera_fecha + datetime.timedelta(days=7 * i)
#         new_fecha_str = new_fecha_date.strftime("%Y/%m/%d")
#         if new_fecha_str in dobles:
#             fechas.append(new_fecha_str)
#         fechas.append(new_fecha_str)
#         i += 1

#     for date in dobles:
#         if date not in fechas:
#             return "n"

#     if len(fechas) != 50:
#         return "n"

#     return "_".join(fechas)

# def format_calendario(calendario: str) -> list[dict[str, str]]:
#     fechas = calendario.split("_")

#     return [
#         {
#             "1 ~ 10": fechas[i],
#             "11 ~ 20": fechas[i+10],
#             "21 ~ 30": fechas[i+20],
#             "31 ~ 40": fechas[i+30],
#             "41 ~ 50": fechas[i+40]
#         } for i in range(10)
#     ]


def obtener_todos_los_ajustes() -> tuple[dict, str]:
    url: str = "http://localhost:8080/ajustes"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        datos: dict = response.json()
        return datos, ""
    except Exception:
        return {}, "No es posible conectarse al servidor"


def crear_calendario(nombre: str, valor_n: float, valor_t: str) -> str:
    url: str = "http://localhost:8080/ajuste"

    try:
        response: requests.Response = requests.put(
            url, json={"ajuste": nombre, "nvalor": int(valor_n), "tvalor": valor_t}
        )

        if not response.ok:
            return response.text

        return ""
    except Exception as e:
        return f"No es posible conectarse al servidor: {e}"


def enviar_rifa_al_servidor(datos) -> str:
    url: str = "http://localhost:8080/ajustes/boletas"

    try:
        response: requests.Response = requests.put(url, json=datos)

        if not response.ok:
            return response.text

        return ""
    except Exception:
        return "No es posible conectarse al servidor"


def obtener_rifas_creadas() -> tuple[dict, str]:
    url: str = "http://localhost:8080/ajustes/boletas"

    try:
        response: requests.Response = requests.get(url)

        if not response.ok:
            return {}, response.text

        datos = response.json()

        return datos, ""

    except Exception:
        return {}, "No es posible conectarse al servidor"
