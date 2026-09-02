from nicegui import ui

import src.functions.prestamos as fp


class Fiadores:
    def __init__(self):
        self.values = []

    def append(self, entry):
        self.values.append(entry)

    def remove(self, entry):
        self.values.remove(entry)

    def get_to_table(self):
        return [
            {"Fiador": int(item["fiador"].value), "Monto": int(item["monto"].value)}
            for item in self.values
            if item["fiador"].value != 0
        ]

    def get_to_api(self):
        return [
            {"id_user": int(item["fiador"].value), "deuda": int(item["monto"].value)}
            for item in self.values
            if item["fiador"].value != 0
        ]


def render_prestamos(idx: int) -> None:
    user, err = fp.obtener_datos_user(idx)

    if err != "":
        ui.notification(err, color="negative")
        return

    ui.label(f"№ {idx} - {user['nombre'].title()}").classes("text-5xl font-bold")

    with ui.tabs().classes("w-full") as tabs:
        prestamos_tab = ui.tab("Prestamos Hechos")
        solicitar_tab = ui.tab("Solicitar Prestamos")
        capital_tab = ui.tab("Consultar Capital")
        cheques_tab = ui.tab("Cheques del Usuario")
    with ui.tab_panels(tabs, value=prestamos_tab).classes("w-full"):
        with ui.tab_panel(prestamos_tab):
            if user["prestamos"] is not None:
                for prestamo in user["prestamos"]:
                    formato_prestamo(prestamo, idx)
                    ui.separator().classes("my-6")

        with ui.tab_panel(solicitar_tab), ui.stepper().props("vertical").classes("w-full") as stepper:
            with ui.step("Imprimir carta"):
                ui.button(
                    "Abrir Carta",
                    icon="open_in_new",
                    color="purple",
                    on_click=lambda: ui.navigate.to(
                        "http://localhost:8080/prestamos/carta", new_tab=True
                    ),
                )
                with ui.stepper_navigation():
                    ui.button("Siguiente", on_click=stepper.next)

            with ui.step("Ingresar monto solicitado"):
                monto = ui.number("Monto del prestamo:", value=1, min=0).classes(
                    "w-50"
                )
                with ui.stepper_navigation():
                    ui.button("Siguiente", on_click=stepper.next)
                    ui.button("Atras", on_click=stepper.previous).props("flat")

            with ui.step("Añadir fiadores"):
                with ui.row():
                    container = ui.column()
                    fiadores = Fiadores()

                    def add_entry_row():
                        with container, ui.row() as row:
                            fiador = ui.number(
                                label="ID Fiador", value=0
                            ).classes("w-32")
                            monto = ui.number(
                                label="Monto a respaldar", value=0
                            ).classes("w-32")

                            # Guardamos el par y la fila (para poder borrarla luego)
                            entry = {
                                "row": row,
                                "fiador": fiador,
                                "monto": monto,
                            }
                            fiadores.append(entry)

                            # Botón individual para eliminar esta fila específica
                            ui.button(
                                icon="delete",
                                on_click=lambda: remove_pair(entry),
                            ).props("flat color=red")

                        ui.notification(
                            "Los campos que tengan un cero como `ID Fiador` seran ignorados",
                            type="info",
                        )

                    def remove_pair(entry):
                        container.remove(entry["row"])
                        fiadores.remove(entry)

                    # ui.button("Obtener Valores", icon="visibility", on_click=get_all_data).classes("ml-2")

                    ui.button(icon="add", on_click=add_entry_row, color="purple")

                with ui.stepper_navigation():
                    ui.button("Siguiente", on_click=stepper.next)
                    ui.button("Atras", on_click=stepper.previous).props("flat")

            with ui.step("Ingresar motivo"):
                motivo = ui.radio(["PG", "ACUERDO"], value="PG")

                dialog = ui.dialog().props(
                    "backdrop-filter='blur(8px) brightness(40%)'"
                )

                def enviar_solicitud():
                    dialog.close()

                    err = fp.crear_prestamo(
                        idx,
                        {
                            "monto": int(monto.value),
                            "motivo": motivo.value,
                            "fiadores": fiadores.get_to_api(),
                        },
                    )

                    if err != "":
                        ui.notify(err, type="negative")
                    else:
                        ui.navigate.to("/prestamos")
                        ui.notify("Prestamo creado correctamente", type="positive")

                def cancelar_solicitud():
                    ui.notify("Solicitud cancelada", type="warning")
                    dialog.close()

                def preparar_dialogo():
                    dialog.clear()  # Limpia contenido anterior
                    datos_fiadores = fiadores.get_to_table()

                    with dialog, ui.column():
                        ui.label("Confirmar Solicitud:").classes(
                            "text-4xl font-bold text-white"
                        )

                        with ui.column().classes("gap-1"):
                            ui.label(f"Monto: ${int(monto.value):,}").classes(
                                "text-lg text-white"
                            )
                            ui.label(f"Motivo: {motivo.value}").classes(
                                "text-lg text-white"
                            )

                        ui.separator().classes("my-2")
                        ui.label("Fiadores:").classes("font-bold text-white")

                        # Tabla simple con los datos actuales
                        if datos_fiadores:
                            ui.table(rows=datos_fiadores).classes("w-full")
                        else:
                            ui.label("No se agregaron fiadores").classes(
                                "text-grey italic"
                            )

                        with ui.row().classes("w-full justify-end mt-4"):
                            ui.button(
                                "Cancelar", on_click=cancelar_solicitud, color="red"
                            ).props("flat")
                            ui.button(
                                "Confirmar y Enviar",
                                color="green",
                                on_click=enviar_solicitud,
                            )

                    dialog.open()

                with ui.stepper_navigation():
                    ui.button(
                        "Iniciar Solicitud",
                        on_click=preparar_dialogo,
                        color="green",
                    )
                    ui.button("Atras", on_click=stepper.previous).props("flat")

        with ui.tab_panel(capital_tab):
            capital = user["capital_disponible"]

            ui.label(f"Capital pago: {capital['capital_pago']:,}").classes(
                "text-2xl font-bold"
            )
            ui.label(f"Capita retirable: {capital['capital_disponible']:,}").classes(
                "text-xl font-bold"
            )

            ui.separator()

            ui.label("Deudas por fiador:").classes("text-2xl font-bold")
            if capital["deudas_por_fiador"] is not None:
                ui.table(rows=capital["deudas_por_fiador"])
            else:
                ui.label("No hay deudas").classes("text-grey italic")
            ui.label(f"Total: {capital['total_deudas_por_fiador']:,}").classes(
                "text-xl font-bold"
            )

            ui.separator()

            ui.label("Deudas en prestamos:").classes("text-2xl font-bold")
            if capital["deudas_en_prestamos"] is not None:
                ui.table(rows=capital["deudas_en_prestamos"])
            else:
                ui.label("No hay deudas").classes("text-grey italic")
            ui.label(f"Total: {capital['total_deudas_en_prestamos']:,}").classes(
                "text-xl font-bold"
            )

            ui.separator()

            ui.label(f"Total Disponible: {capital['total_disponible']:,}").classes(
                "text-3xl font-bold"
            )

        with ui.tab_panel(cheques_tab):
            if user["cheques_user"] is None:
                return

            cheque_select = ui.select(
                user["cheques_user"],
                value=user["cheques_user"][0] if user["cheques_user"] else None,
                label="Seleccionar cheque a solicitar:",
            ).classes("w-full")

            with ui.row():
                ui.button(
                    "Solicitar",
                    on_click=lambda: ui.navigate.to(
                        f"http://localhost:8080/prestamos/cheque/{cheque_select.value}",
                        new_tab=True
                    ),
                )

            ui.separator()

            with ui.grid(columns=5):
                for cheque in user["cheques"]:
                    with ui.card():
                        ui.label(f"Codigo: {cheque['codigo']}").classes(
                            "text-1xl font-bold"
                        )
                        ui.markdown(f"""
                        + En prestamo: `{cheque["codigo_prestamo"]}`
                        + Monto: `{cheque["monto"]:,}`
                        + Deuda antes: `{cheque["deuda_antes"]:,}`
                        + Deuda despues: `{cheque["deuda_despues"]:,}`
                        + Fecha: `{cheque["fecha"]}`
                        """)


def formato_prestamo(prestamo, idx):
    with ui.row():
        with ui.column():
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
                        ui.label(f"Deuda Act: {prestamo['deuda_actual']:,}").classes(
                            "text-2xl font-bold"
                        )

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

                if prestamo["fiadores"] is not None:
                    with ui.card():
                        ui.label("Fiadores:").classes("text-2xl font-bold")
                        ui.table(rows=prestamo["fiadores"])

            with ui.row():
                with ui.card():
                    ui.label("Notas de prestamo:").classes("text-2xl font-bold")
                    with ui.list().props("dense separator"):
                        for nota in prestamo["anotaciones"]:
                            ui.item(nota)

                if prestamo["deuda_actual"] > 0:
                    with ui.card():
                        ui.label("Abonar al prestamo:").classes("text-2xl font-bold")
                        tipo_de_abono = ui.radio(["DEUDA", "INTERESES"], value="DEUDA")

                        monto = ui.number("Monto a abonar:", value=1, min=0).classes(
                            "w-full"
                        )

                        dialog = ui.dialog().props(
                            "backdrop-filter='blur(8px) brightness(40%)'"
                        )

                        def enviar_solicitud():
                            dialog.close()

                            err = fp.abonar_a_prestamo(
                                idx,
                                {
                                    "monto": int(monto.value),
                                    "codigo": prestamo["codigo"],
                                    "tipo": tipo_de_abono.value,
                                },
                            )

                            if err != "":
                                ui.notify(err, type="negative")
                            else:
                                ui.navigate.to("/prestamos")
                                ui.notify(
                                    "Abono relizado correctamente", type="positive"
                                )

                        def cancelar_solicitud():
                            ui.notify("Solicitud cancelada", type="warning")
                            dialog.close()

                        def preparar_dialog():
                            dialog.clear()

                            with dialog, ui.column():
                                ui.label("Confirmar Abono:").classes(
                                    "text-4xl font-bold text-white"
                                )

                                ui.label(
                                    f"Abonar a prestamo: {prestamo['codigo']}"
                                ).classes("text-xl font-bold text-white")
                                ui.label(
                                    f"Monto a abonar: {int(monto.value):,}"
                                ).classes("text-xl font-bold text-white")
                                ui.label(
                                    f"Abono a: {tipo_de_abono.value}"
                                ).classes("text-xl font-bold text-white")

                                with ui.row().classes("w-full justify-end mt-4"):
                                    ui.button(
                                        "Cancelar",
                                        on_click=cancelar_solicitud,
                                        color="red",
                                    ).props("flat")
                                    ui.button(
                                        "Confirmar y Enviar",
                                        color="green",
                                        on_click=enviar_solicitud,
                                    )

                            dialog.open()

                        ui.button(
                            "Iniciar Abono", on_click=preparar_dialog, color="green"
                        ).classes("w-full")

        if prestamo["deuda_actual"] > 0:
            with ui.card():
                ui.label("Fechas de pago:").classes("text-2xl font-bold")
                with ui.list().props("dense separator"):
                    for fecha in prestamo["fechas_de_pago"]:
                        ui.item(fecha)
