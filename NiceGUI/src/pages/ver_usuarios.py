import src.functions.ver_usuarios as fvu

from nicegui import ui


def render_ver_usuarios(idx: int) -> None:
    with ui.tabs().classes("w-full") as tabs:
        buscar_tab = ui.tab("Buscar Usuarios")
        acuerdos_tab = ui.tab("Ver Acuerdo")
        boletas_tab = ui.tab("Buscar Boletas")
        prestamos_tab = ui.tab("Buscar Prestamo")
    with ui.tab_panels(tabs, value=buscar_tab).classes("w-full"):
        with ui.tab_panel(buscar_tab):
            with ui.card().classes("w-full"):
                ui.label("Buscar por nombre o fragmento:").classes("text-3xl font-bold")

                nombre_buscar = ui.input(
                    "Fragmento a buscar",
                    placeholder="vacio para ver todos los usuarios",
                    value="",
                ).classes("w-full")

                def refrescar_usuarios() -> None:
                    users, err = fvu.obtener_usuarios(nombre_buscar.value)
                    if err != "":
                        ui.notify(err, type="negative")
                        return

                    tabla_usuarios.rows[:] = users
                    tabla_usuarios.update()

                ui.button("Buscar", icon="search", on_click=refrescar_usuarios)

            tabla_usuarios = ui.table(
                columns=[
                    {"name": "id", "label": "ID", "field": "id"},
                    {"name": "nombre", "label": "Nombre", "field": "nombre"},
                    {"name": "telefono", "label": "Telefono", "field": "telefono"},
                    {"name": "puestos", "label": "Puestos", "field": "puestos"},
                    {"name": "estado", "label": "Estado", "field": "estado"},
                ],
                rows=[],
            ).classes("w-full")

            refrescar_usuarios()

        with ui.tab_panel(acuerdos_tab):
            with ui.card().classes("w-full bg-blue-300"):
                ui.label(
                    "Todos los usuarios acontinuacion necesitan firmar acuerdo"
                ).classes("text-base")

            acuerdos, err = fvu.obtener_acuerdos()
            if err == "":
                ui.table(rows=acuerdos).classes("w-full")
            else:
                with ui.card().classes("w-full bg-red-300"):
                    ui.label(err).classes("text-base")

        with ui.tab_panel(boletas_tab):
            pass

        @ui.refreshable
        def prestamo_mostar(id, nombre, prestamo) -> None:
            with ui.row():
                with ui.column():
                    with ui.card().classes("w-full"):
                        ui.label(
                            f"ID Usuario: {id} - Nombre: {nombre.title()}"
                        ).classes("text-3xl font-bold")

                    with ui.row():
                        with ui.column():
                            with ui.card().classes("w-full"):
                                ui.label(f"Codigo: {prestamo['codigo']}").classes(
                                    "text-3xl font-bold"
                                )

                            with ui.card():
                                ui.label(f"Estado: {prestamo['estado']}").classes(
                                    "text-2xl font-bold"
                                )
                                ui.label(f"Motivo: {prestamo['motivo']}").classes(
                                    "text-2xl font-bold"
                                )
                                ui.label(
                                    f"Deuda Act: {prestamo['deuda_actual']:,}"
                                ).classes("text-2xl font-bold")

                        with ui.card():
                            ui.label("Intereses:").classes("text-2xl font-bold")
                            ui.markdown(f"+ **Interes**: {prestamo['interes']}%")
                            ui.markdown(
                                f"+ **Intereses Vencidos**: {prestamo['interes_vencido']:,}"
                            )
                            ui.markdown(
                                f"+ **Intereses Generados**: {prestamo['interes_generado']:,}"
                            )

                        with ui.card():
                            ui.label("Deudas:").classes("text-2xl font-bold")
                            ui.markdown(f"+ **Deuda**: {prestamo['deuda']:,}")
                            ui.markdown(f"+ **Monto**: {prestamo['monto']:,}")
                            ui.markdown(
                                f"+ **Pago**: {int(100 * (1 - prestamo['deuda'] / prestamo['monto']))}%"
                            )

                        with ui.card():
                            ui.label("Fiadores:").classes("text-2xl font-bold")
                            if prestamo["fiadores"] is not None:
                                ui.table(rows=prestamo["fiadores"])

                    with ui.row():
                        with ui.card():
                            ui.label("Notas de prestamo:").classes("text-2xl font-bold")
                            with ui.list().props("dense separator"):
                                for nota in prestamo["anotaciones"]:
                                    ui.item(nota)

                with ui.card():
                    ui.label("Fechas de pago:").classes("text-2xl font-bold")
                    with ui.list().props("dense separator"):
                        for fecha in prestamo["fechas_de_pago"]:
                            ui.item(fecha)

        with ui.tab_panel(prestamos_tab):
            with ui.card().classes("w-full"):
                ui.label("Buscar prestamo:").classes("text-3xl font-bold")

                with ui.row().classes("items-center"):
                    codigo_prestamo = ui.number(
                        "Codigo del prestamo a buscar:", value=2, min=1
                    ).classes("w-60")

                    def actulizar_prestamo() -> None:
                        prestamo, err = fvu.obtener_prestamo(int(codigo_prestamo.value))
                        if err != "":
                            ui.notify(err, type="negative")
                            return

                        prestamo_mostar.refresh(
                            prestamo["id"], prestamo["nombre"], prestamo["prestamo"]
                        )

                    ui.button("Buscar", icon="search", on_click=actulizar_prestamo)

            prestamo, _ = fvu.obtener_prestamo(2)

            prestamo_mostar(prestamo["id"], prestamo["nombre"], prestamo["prestamo"])
