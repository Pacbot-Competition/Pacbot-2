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

class PacbotServer:
    async def handler(self, websocket):
        while True:
            await websocket.send("changeColor 1 1 #ff0000") # For testing
            print("Waiting for message")
            message = await websocket.recv()
            print(message)

    async def run(self):
        print("Starting server")
        port = get_server_port()
        async with websockets.serve(self.handler, "", port):
            print("Server started on port", port)
            await asyncio.Future()

if __name__ == "__main__":
    server=PacbotServer()
    asyncio.run(server.run())