from fastapi import (
  APIRouter, 
  Request, 
  HTTPException
)

import random
import requests
import re
from bs4 import BeautifulSoup
from fastapi.encoders import jsonable_encoder
from fastapi.responses import JSONResponse
from pydantic import BaseModel
import random
from cachetools import cached, TTLCache

class RockerResponse(BaseModel):
    id: int
    rocker: str

rock_router = APIRouter()
cache = TTLCache(maxsize=100, ttl=86400)

ROCKERS_URL = "https://parade.com/1020922/jessicasager/best-rock-bands-of-all-time/"


@cached(cache)
@rock_router.get("/assign-rocker")
def new_random_rocker(request: Request):
  rock_request   = requests.get(ROCKERS_URL)

  if (not rock_request.ok):
    print(rock_request)
    raise HTTPException(status_code=500, detail="No Rockers found!")

  soup = BeautifulSoup(rock_request.text)
  rock_names = soup.find_all('h2')

  rand_index  = random.randint(0, len(rock_names) - 1)
  rand_rocker = rock_names[rand_index]
  name = re.split(r'[0-9]{0,3}\.\s', rand_rocker.text)

  if (len(name) < 2):
    raise HTTPException(status_code=500, detail="Bad Name! Try Again!")

  json_compatible_item_data = jsonable_encoder({'id': '1', 'rocker': name[1]})
  return JSONResponse(content=json_compatible_item_data)
