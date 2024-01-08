import asyncio
import json
import websockets

def get_server_port():
    '''
    Get the port to connect to from the config file
    '''
    with open('../config.json', 'r', encoding='UTF-8') as configFile:
        config = json.load(configFile)
    return config["BotSocketPort"]

class DebugServer:
    def __init__(self) -> None:
        self.clients = []
        self.on_cell_clicked = lambda row, col: None

    async def handler(self, websocket, path):
        print(f"[DEBUG SERVER] Client connected: {websocket.remote_address}")
        self.clients.append(websocket)

        try:
            async for message in websocket:
                print(f"[DEBUG SERVER] Received message from {websocket.remote_address}: {message}")

                message_components = message.split(" ")
                if message_components[0] == "clicked":
                    row = int(message_components[1])
                    col = int(message_components[2])
                    self.on_cell_clicked(row, col)

        finally:
            print(f"[DEBUG SERVER] Client disconnected: {websocket.remote_address}")
            self.clients.remove(websocket)

    async def run(self):
        print("[DEBUG SERVER] Starting server")
        port = get_server_port()
        self.server_socket = await websockets.serve(self.handler, "", port)
        print("[DEBUG SERVER] Server started on port", port)

        await asyncio.Future()

    async def broadcast(self, message):
        for client in self.clients:
            await client.send(message)

    def set_cell_color(self, row, col, color):
        if type(color) == tuple:
            color = f"rgb({color[0]}, {color[1]}, {color[2]})"

        asyncio.create_task(self.broadcast(f"set_cell_color {row} {col} {color}"))

    def reset_cell_colors(self):
        asyncio.create_task(self.broadcast("reset_all_cell_colors"))