import asyncio

import websockets

class PacbotServer:
    async def handler(self, websocket):
        while True:
            message = await websocket.recv()
            print(message)

    async def run(self):
        async with websockets.serve(self.handler, "", 1000):
            await asyncio.Future()

if __name__ == "__main__":
    server=PacbotServer()
    asyncio.run(server.run())