from nicegui import ui


def render_analizar(idx: int) -> None:
    with ui.card():
        ui.label("Consultar Boletas:").classes("text-3xl font-blod")
        ui.button("Consultar", on_click=lambda: ui.navigate.to(
            "http://localhost:8080/analizar/boletas", new_tab=True
        ))