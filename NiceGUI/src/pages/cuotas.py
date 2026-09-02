import src.functions.cuotas as fc

from nicegui import ui


def render_cuotas(idx: int) -> None:
    user, err = fc.obtener_datos_user(idx)

    if err != "":
        ui.notification(err, color="negative")
        return

    ui.label(
        f"№ {idx} - {user['nombre'].title()} : {user['puestos']} puesto(s)"
    ).classes("text-5xl font-bold")

    ui.label(f"Numero de telefono: {user['telefono']}").classes("text-3xl font-bold")

    with ui.tabs().classes("w-full") as tabs:
        cuotas_tab = ui.tab("Pagos Semanales")
        cheques_tab = ui.tab("Cheques del Usuario")
    with ui.tab_panels(tabs, value=cuotas_tab).classes("w-full"):
        with ui.tab_panel(cuotas_tab):
            with ui.row():
                ui.table(rows=user["tabla"], pagination=50)

                with ui.card():
                    ui.label("Proceso de pago").classes(
                        "text-2xl font-bold text-blue-400"
                    )

                    cuotas_a_pagar = ui.select(
                        list(range(user["cuotas_a_pagar"] + 1)),
                        value=0,
                        label="Cuotas a Pagar:",
                    ).classes("w-full")

                    multas_a_pagar = ui.select(
                        list(range(user["multas_a_pagar"] + 1)),
                        value=0,
                        label="Multas a Pagar:",
                    ).classes("w-full")

                    metodo_de_pago = ui.radio(
                        ["EFECTIVO", "TRANSFERENCIA"], value="EFECTIVO"
                    )

                    dialog = ui.dialog().props(
                        "backdrop-filter='blur(8px) brightness(40%)'"
                    )

                    def enviar_solicitud():
                        dialog.close()

                        codigo, err = fc.pagar_semanas(
                            idx,
                            {
                                "cuotas": cuotas_a_pagar.value,
                                "multas": multas_a_pagar.value,
                                "metodo": metodo_de_pago.value,
                            },
                        )

                        if err != "":
                            ui.notify(err, type="negative")
                        else:
                            ui.navigate.to("/cuotas")
                            ui.notify(
                                f"Pago relizado correctamente, cheque de pago {codigo}",
                                type="positive",
                            )

                    def cancelar_solicitud():
                        ui.notify("Solicitud cancelada", type="warning")
                        dialog.close()

                    def preparar_dialogo():
                        dialog.clear()  # Limpia contenido anterior

                        with dialog, ui.column():
                            ui.label(f"№ {idx} - {user['nombre']}").classes(
                                "text-4xl font-bold text-white"
                            )
                            ui.label(f"Puestos: {user['puestos']}").classes(
                                "text-lg text-white"
                            )
                            ui.separator().classes("my-2")

                            with ui.column().classes("gap-1"):
                                ui.label(
                                    f"Cuotas a pagar: {cuotas_a_pagar.value}"
                                ).classes("text-lg text-white")
                                ui.label(
                                    f"Valor de cuota por puesto: {
                                        user['valor_cuota']:,}"
                                ).classes("text-lg text-white")
                                total_cuotas: int = (
                                    cuotas_a_pagar.value
                                    * user["valor_cuota"]
                                    * user["puestos"]
                                )
                                ui.label(f"Total en cuotas: {total_cuotas:,}").classes(
                                    "text-xl text-white"
                                )

                                ui.separator().classes("my-2")

                                ui.label(
                                    f"Multas a pagar: {multas_a_pagar.value}"
                                ).classes("text-lg text-white")
                                ui.label(
                                    f"Valor de multa por puesto: {
                                        user['valor_multa']:,}"
                                ).classes("text-lg text-white")
                                total_multas = (
                                    multas_a_pagar.value
                                    * user["valor_multa"]
                                    * user["puestos"]
                                )
                                ui.label(f"Total en multas: {total_multas:,}").classes(
                                    "text-xl text-white"
                                )

                                ui.separator().classes("my-2")

                                ui.label(
                                    f"Metodo de pago: {metodo_de_pago.value}"
                                ).classes("text-lg text-white")
                                ui.label(
                                    f"Total: {total_cuotas + total_multas:,}"
                                ).classes("text-2xl font-bold text-white")

                            ui.separator().classes("my-2")

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

                    ui.button(
                        "Iniciar Pago", on_click=preparar_dialogo, color="green"
                    ).classes("w-full")

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
                        f"http://localhost:8080/cuotas/cheque/{cheque_select.value}",
                        new_tab=True,
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
                        + Cuotas pagas: `{cheque["cuotas_pagas"]}`
                        + Multas pagas: `{cheque["multas_pagas"]}`
                        + Metodo de pago: `{cheque["metodo_pago"]}`
                        + Fecha: `{cheque["fecha_y_hora"]}`
                        """)
