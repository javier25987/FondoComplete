import src.pages.ajustes as pa
import src.pages.ver_usuarios as pvu
import src.pages.prestamos as pp
import src.pages.modificar as fm
import src.pages.cuotas as pc
import src.pages.deudas as pd
import src.pages.insertar_boletas as pib
import src.pages.rifas as pr

from nicegui import app, ui

# Funciones de Estado y Autenticación


class BindUser:
    def __init__(self):
        self.number = app.storage.user.get("user", 1)

    def increment(self):
        self.number += 1

    def decrement(self):
        if self.number > 1:
            self.number -= 1


def log_in() -> None:
    app.storage.user["is_admin"] = True
    ui.notify("Sesión iniciada como Administrador", color="positive")
    slide_bar.refresh()


def log_out() -> None:
    app.storage.user["is_admin"] = False
    ui.notify("Sesión cerrada", color="negative")
    slide_bar.refresh()


def actualizar_user(nuevo_valor):
    app.storage.user["user"] = int(nuevo_valor)
    ui.notify(f"Cargando datos del usuario: {nuevo_valor}", color="info")
    # Refrescamos el contenido principal de la página
    contenido_principal.refresh()


# --- Componentes de la Interfaz ---


@ui.refreshable
def slide_bar():
    prop: str = "flat align=left"
    clas: str = "w-full justify-start"

    es_admin = app.storage.user.get("is_admin", False)

    with ui.column().classes("w-full gap-2 p-4"):
        ui.label("Páginas Generales").classes("text-base font-bold text-gray-400 mb-2")

        ui.button("Menu", icon="home", on_click=lambda: ui.navigate.to("/")).props(
            prop
        ).classes(clas)

        ui.button(
            "Cuotas", icon="savings", on_click=lambda: ui.navigate.to("/cuotas")
        ).props(prop).classes(clas)

        ui.button(
            "Prestamos",
            icon="attach_money",
            on_click=lambda: ui.navigate.to("/prestamos"),
        ).props(prop).classes(clas)

        ui.button(
            "Analizar", icon="analytics", on_click=lambda: ui.navigate.to("/analizar")
        ).props(prop).classes(clas)

        ui.button(
            "Transferencias",
            icon="account_balance",
            on_click=lambda: ui.navigate.to("/transferencias"),
        ).props(prop).classes(clas)

        ui.button(
            "Rifas", icon="note", on_click=lambda: ui.navigate.to("/rifas")
        ).props(prop).classes(clas)

        ui.button(
            "Deudas", icon="payments", on_click=lambda: ui.navigate.to("/deudas")
        ).props(prop).classes(clas)

        # ui.button(
        #     "Anotaciones", icon="book", on_click=lambda: ui.navigate.to("/anotaciones")
        # ).props(prop).classes(clas)

        ui.button(
            "Ver Usuarios",
            icon="person_search",
            on_click=lambda: ui.navigate.to("/verusuarios"),
        ).props(prop).classes(clas)

        # ui.button(
        #     "Registros",
        #     icon="library_books",
        #     on_click=lambda: ui.navigate.to("/registros"),
        # ).props(prop).classes(clas)

        ui.label("Paginas Administrativas").classes(
            "text-base font-bold text-gray-400 mt-4 mb-2"
        )

        if es_admin:
            ui.button(
                "Ingresar Boletas",
                icon="library_add",
                on_click=lambda: ui.navigate.to("/insertarboletas"),
            ).props(prop).classes(clas)

            ui.button(
                "Modificar",
                icon="change_circle",
                on_click=lambda: ui.navigate.to("/modificar"),
            ).props(prop).classes(clas)

            ui.button(
                "Ajustes", icon="settings", on_click=lambda: ui.navigate.to("/ajustes")
            ).props(prop).classes(clas)

            ui.button("LogOut", icon="logout", on_click=log_out).props("flat").classes(
                "w-full justify-start text-red"
            )
        else:
            ui.button("Login Admin", icon="login", on_click=log_in).props(
                "flat"
            ).classes("w-full justify-start text-primary")


def layout(titulo_pagina):  # barra superior y barra lateral
    with (
        ui.header(elevated=True)
        .style("background-color: #3874c8")
        .classes("items-center")
    ):
        ui.button(on_click=lambda: left_drawer.toggle(), icon="menu").props(
            "flat color=white"
        )
        ui.label(f"Fondo San Javier - {titulo_pagina}").classes("font-bold text-lg")

    with ui.left_drawer(fixed=False).style("background-color: #f8f9fa") as left_drawer:
        slide_bar()

        ui.separator().classes("my-2")

        number_user = BindUser()

        # Ingresar nuevo usuario
        num_input = (
            ui.number(
                min=1,
                label="Puesto del usuario:",
            )
            .classes("text-xl w-full")
            .bind_value(number_user, "number")
        )

        with ui.row().classes("w-full"):
            ui.button(icon="add", on_click=lambda: number_user.increment()).classes(
                "flex-1"
            )
            ui.button(icon="remove", on_click=lambda: number_user.decrement()).classes(
                "flex-1"
            )

        ui.button(
            "Buscar",
            color="green",
            on_click=lambda: actualizar_user(num_input.value),
            icon="search",
        ).classes("w-full")


@ui.refreshable
def contenido_principal(contenido):
    idx = app.storage.user.get("user", 1)
    contenido(idx)


# Contenido de la pagina menu


def render_menu(idx):
    with ui.column().classes("w-full items-center mt-10"):
        ui.icon("rocket", size="lg", color="primary")
        ui.label("Bienvenido").classes("text-h3")


def render_analisis(idx):
    pass


# Rutas a cada pagina


@ui.page("/")
def pagina_menu():
    layout("Menu")
    # Pasamos la función que dibuja el contenido
    contenido_principal(render_menu)


@ui.page("/cuotas")
def pagina_cuotas():
    layout("Cuotas")
    contenido_principal(pc.render_cuotas)


@ui.page("/prestamos")
def pagina_prestamos():
    layout("Prestamos")
    contenido_principal(pp.render_prestamos)


@ui.page("/analizar")
def pagina_analizar():
    layout("Analizar")
    contenido_principal(render_analisis)


@ui.page("/transferencias")
def pagina_transferencias():
    layout("Transferencias")
    contenido_principal(render_analisis)


@ui.page("/rifas")
def pagina_rifas():
    layout("Rifas")
    contenido_principal(pr.render_rifas)


@ui.page("/deudas")
def pagina_deudas():
    layout("Deudas")
    contenido_principal(pd.render_deudas)


# @ui.page("/anotaciones")
# def pagina_anotaciones():
#     layout("Anotaciones")
#     contenido_principal(render_analisis)


@ui.page("/verusuarios")
def pagina_versocios():
    layout("Ver Usuarios")
    contenido_principal(pvu.render_ver_usuarios)


# @ui.page("/registros")
# def pagina_registos():
#     layout("Registros")
#     contenido_principal(render_analisis)


@ui.page("/insertarboletas")
def pagina_insertar_boletas():
    if not app.storage.user.get("is_admin", False):
        ui.navigate.to("/")
        return

    layout("Ingresar Boletas")

    pib.render_insertar_boletas()


@ui.page("/modificar")
def pagina_modificar():
    if not app.storage.user.get("is_admin", False):
        ui.navigate.to("/")
        return

    layout("Modificar Usuarios")

    fm.render_modificar()


@ui.page("/ajustes")
def pagina_ajustes():
    if not app.storage.user.get("is_admin", False):
        ui.navigate.to("/")
        return

    layout("Ajustes")

    pa.render_ajutes()


# Ejecución
ui.run(title="Fondo San Javier", port=8088, storage_secret="1026551455_2026")
