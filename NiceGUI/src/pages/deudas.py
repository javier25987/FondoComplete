import src.functions.deudas as fd

from nicegui import ui


def render_deudas(idx: int) -> None:
    user, err = fd.obtener_deudas(idx)
    if err != "":
        ui.notify(err, type="negative")
        return

    ui.label(f"№ {idx} - {user['nombre'].title()}").classes("text-5xl font-bold")

    with ui.card().classes("w-full"):
        ui.label("Crear una nueva deuda:").classes("text-3xl font-bold")

        motivo_deuda = ui.input("Motivo de la deuda:").classes("w-full")

        with ui.row().classes("w-full items-center"):
            monto_deuda = ui.number("Monto de la deuda:", value=0, min=0).classes(
                "w-80"
            )

            tipo_deuda = ui.select(
                ["MULTA POR CABEZA", "MULTA EXTRA"],
                label="Tipo de la deuda:",
                value="MULTA EXTRA",
            ).classes("w-80")

            ui.space()

            dialog = ui.dialog().props("backdrop-filter='blur(8px) brightness(40%)'")

            def enviar_solicitud():
                dialog.close()

                err = fd.crear_deuda(
                    idx,
                    {
                        "motivo": motivo_deuda.value,
                        "monto": int(monto_deuda.value),
                        "tipo": tipo_deuda.value,
                    },
                )

                if err != "":
                    ui.notify(err, type="negative")
                else:
                    ui.navigate.to("/deudas")

            def cancelar_solicitud():
                ui.notify("Solicitud cancelada", type="warning")
                dialog.close()

            def preparar_dialog():
                dialog.clear()

                with dialog, ui.column():
                    ui.label("Crear deuda:").classes("text-4xl font-bold text-white")

                    ui.label(f"Motivo: {motivo_deuda.value}").classes(
                        "text-xl font-bold text-white"
                    )

                    ui.label(f"Monto: {int(monto_deuda.value):,}").classes(
                        "text-xl font-bold text-white"
                    )

                    ui.label(f"Tipo: {tipo_deuda.value}").classes(
                        "text-xl font-bold text-white"
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

            ui.button("Crear", icon="create", on_click=preparar_dialog)

    if not user["deudas"]:
        return

    with ui.card().classes("w-full"):
        ui.label("Deudas del usuario:").classes("text-3xl font-bold")

    with ui.grid(columns=3).classes("w-full"):
        for deuda in user["deudas"]:
            formato_deuda(deuda, idx)


def formato_deuda(deuda, idx) -> None:
    with ui.card():
        ui.markdown(f"**Codigo:** `{deuda['codigo']}`").classes("text-xl")
        ui.label("Motivo:").classes("text-base font-bold")
        ui.label(deuda["motivo"]).classes("text-base font-bold text-grey")
        ui.markdown(f"**Monto:** `{deuda['monto']:,}`").classes("text-base")
        ui.markdown(f"**Tipo:** `{deuda['tipo']}`").classes("text-base")
        ui.markdown(f"**Fecha de creacion:** `{deuda['fecha_de_creacion']}`").classes(
            "text-base"
        )

        dialog = ui.dialog().props("backdrop-filter='blur(8px) brightness(40%)'")

        def enviar_solicitud():
            dialog.close()

            err = fd.pagar_deuda(idx, deuda["codigo"])

            if err != "":
                ui.notify(err, type="negative")
            else:
                ui.navigate.to("/deudas")

        def cancelar_solicitud():
            ui.notify("Solicitud cancelada", type="warning")
            dialog.close()

        def preparar_dialog():
            dialog.clear()

            with dialog, ui.column():
                ui.label("Confirmar pago:").classes("text-4xl font-bold text-white")

                ui.label(f"Deuda a pagar: {deuda['codigo']}").classes(
                    "text-xl font-bold text-white"
                )

                ui.label(f"Monto a pagar: {deuda['monto']:,}").classes(
                    "text-xl font-bold text-white"
                )

                with ui.row().classes("w-full justify-end mt-4"):
                    ui.button(
                        "Cancelar", on_click=cancelar_solicitud, color="red"
                    ).props("flat")
                    ui.button(
                        "Confirmar y Enviar", color="green", on_click=enviar_solicitud
                    )

            dialog.open()

        if deuda["estado"]:
            ui.button("Pagar Deuda", icon="payments", on_click=preparar_dialog)
        else:
            ui.markdown(f"**Fecha de pago:** `{deuda['fecha_de_pago']}`").classes(
                "text-base"
            )
