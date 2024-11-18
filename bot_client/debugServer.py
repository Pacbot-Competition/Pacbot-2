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

def format_color(color):
    if type(color) == tuple:
        return f"rgb({color[0]},{color[1]},{color[2]})"
    else:
        return color

class DebugServer:
    def __init__(self) -> None:
        self.clients = []
        self.on_cell_clicked = lambda row, col: None
        self.is_resetting = False

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
        asyncio.create_task(self.broadcast(f"set_cell_color {row} {col} {format_color(color)}"))

    def reset_cell_colors(self):
        asyncio.create_task(self.broadcast("reset_all_cell_colors"))

    def set_cell_color_multiple(self, positions, color):
        asyncio.create_task(self.broadcast(f"set_cell_colors {' '.join(map(lambda pos: f'{pos[0]} {pos[1]}', positions))} {format_color(color)}"))

    def set_cell_colors(self, new_cell_colors):
        asyncio.create_task(self.broadcast(f"set_cell_colors {' '.join(map(lambda ncc: f'{ncc[0][0]} {ncc[0][1]} {format_color(ncc[1])}', new_cell_colors))}"))

    def set_path(self, path):
        asyncio.create_task(self.broadcast(f"set_path {' '.join(map(lambda pos: f'{pos[0]} {pos[1]}', path))}"))
        
    async def pause_game(self):
        await asyncio.create_task(self.broadcast("pause_game"))
        
    async def reset_game(self):
        if not self.is_resetting:
            self.is_resetting = True
            await asyncio.create_task(self.broadcast("reset_game"))
            self.is_resetting = False
