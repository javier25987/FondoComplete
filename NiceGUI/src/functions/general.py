import subprocess
import webbrowser
import os


def open_pdf_in_chrome(path):
    chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
    file_path = os.path.abspath(path)

    if os.path.exists(chrome_path):
        subprocess.Popen([chrome_path, file_path])
    else:
        print("No se encontró Chrome en la ruta especificada.")


def open_pdf(path: str) -> None:
    use_chorome: bool = False

    if use_chorome:
        open_pdf_in_chrome(path)
        return

    webbrowser.open_new(path)
