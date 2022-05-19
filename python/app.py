from fastapi import FastAPI

from routers.encryption import encrypt_router

def init_app():
  """Initializes the app object and adds the router to it
  """
  app = FastAPI()
  app.include_router(encrypt_router)
  return app

app = init_app()