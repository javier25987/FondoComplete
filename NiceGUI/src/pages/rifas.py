from nicegui import ui
import src.functions.rifas as fr

def render_rifas(idx: int) -> None:
    user, err = fr.obtener_datos_user(idx)
    if err != "":
        ui.notification(err, type="negative")
        return

    if not user["boletas_disponibles"]:
        user["boletas_disponibles"] = []

    ui.label(
        f"№ {idx} - {user["nombre"].title()}"
    ).classes("text-5xl font-bold")

    with ui.tabs().classes("w-full") as tabs:
        rifas_tab = ui.tab("Boletas del usuario")
        cheques_tab = ui.tab("Cheques de pagos")
    with ui.tab_panels(tabs, value=rifas_tab).classes("w-full"):
        with ui.tab_panel(rifas_tab):
            with ui.row().classes("w-full"):
                with ui.card().classes("flex-1 h-64"):
                    ui.label("Entregar boletas:").classes("text-3xl font-bold")

                    boletas_a_entregar = ui.select(
                        user["boletas_disponibles"], 
                        multiple=True, 
                        label="Boletas a entregar:"
                    ).classes('w-full').props('use-chips')

                    dialog_1 = ui.dialog().props(
                        "backdrop-filter='blur(8px) brightness(40%)'"
                    )

                    def entregar_boletas():
                        dialog_1.close()

                        boletas = [int(x) for x in boletas_a_entregar.value]

                        err = fr.entregar_boletas(idx, {"boletas": boletas})

                        if err != "":
                            ui.notification(err, type="negative")
                            return

                        ui.navigate.to("/rifas")
                        ui.notify(
                            "Boletas entregadas correctamente",
                            type="positive",
                        )

                    def cancelar_solicitud():
                        ui.notify("Solicitud cancelada", type="warning")
                        dialog_1.close()

                    def preparar_dialogo_1():
                        dialog_1.clear()  # Limpia contenido anterior

                        with dialog_1, ui.column():
                            ui.label(f"№ {idx} - {user["nombre"].title()}").classes(
                                "text-4xl font-bold text-white"
                            )
                            ui.separator().classes("my-2")

                            ui.label(
                                f"Boletas a entregar: {boletas_a_entregar.value}"
                            ).classes("text-lg text-white")

                            ui.separator().classes("my-2")

                            with ui.row().classes("w-full justify-end mt-4"):
                                ui.button(
                                    "Cancelar", on_click=cancelar_solicitud, color="red"
                                ).props("flat")
                                ui.button(
                                    "Confirmar",
                                    color="green",
                                    on_click=entregar_boletas,
                                )

                        dialog_1.open()

                    ui.button("Entregar", on_click=preparar_dialogo_1, color="green")

                with ui.card().classes("flex-1 h-64"):
                    ui.label("Pagar boletas:").classes("text-3xl font-bold")

                    ui.label(f"Deuda por boletas: {user["deuda_por_boletas"]:,}").classes("text-xl")

                    monto_a_pagar = ui.number("Monto a pagar:").classes("w-full")
                    
                    dialog_2 = ui.dialog().props(
                        "backdrop-filter='blur(8px) brightness(40%)'"
                    )

                    def pagar_deuda():
                        dialog_2.close()

                        if monto_a_pagar.value is None:
                            ui.notification("Escriba al menos un valor", type="warning")
                            return

                        err = fr.pagar_deuda_por_boletas(idx, {
                            "monto": int(monto_a_pagar.value)
                        })

                        if err != "":
                            ui.notification(err, type="negative")
                            return

                        ui.navigate.to("/rifas")
                        ui.notify("Pago relizado correctamente", type="positive")


                    def cancelar_solicitud():
                        ui.notify("Solicitud cancelada", type="warning")
                        dialog_2.close()

                    def preparar_dialogo_2():
                        dialog_2.clear()  # Limpia contenido anterior

                        with dialog_2, ui.column():
                            ui.label(f"№ {idx} - {user["nombre"].title()}").classes(
                                "text-4xl font-bold text-white"
                            )
                            ui.separator().classes("my-2")

                            ui.label(f"Monto a pagar: {int(monto_a_pagar.value):,}").classes("text-lg text-white")

                            ui.separator().classes("my-2")

                            with ui.row().classes("w-full justify-end mt-4"):
                                ui.button(
                                    "Cancelar", on_click=cancelar_solicitud, color="red"
                                ).props("flat")
                                ui.button(
                                    "Confirmar",
                                    color="green",
                                    on_click=pagar_deuda,
                                )

                        dialog_2.open()

                    ui.button("Pagar", on_click=preparar_dialogo_2, color="green")

            with ui.card().classes("w-full"):
                ui.label("Boletas entregadas:").classes("text-3xl font-bold")

                with ui.grid(columns=4).classes("gap-x-30 text-2xl"):
                    
                    if user["boletas_entregadas"]:
                        for boleta in user["boletas_entregadas"]:
                            ui.label(boleta)
