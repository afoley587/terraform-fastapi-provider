import ast
from Crypto.Cipher import DES
from fastapi import APIRouter
import hashlib
import os

from schemas.encryption import (
  DESEncryptedRequest,
  DESEncryptedResponse,
  DESDecryptedRequest,
  DESDecryptedResponse
)

# To be used for encryption, must be 8 bytes long
KEY            = bytes(os.environ['DES_KEY'], 'utf-8')
encrypt_router = APIRouter()

def pad(text):
  """Pads byte string so it is a multiple of 8 bytes long
  """
  n = len(text) % 8
  return text + (b' ' * (8 - n))

@encrypt_router.post("/encrypt/des")
async def enc_des(
  request: DESEncryptedRequest, 
  response_model=DESEncryptedResponse
):
  """Encrypts a plaintext string denoted in the payload
  """
  plaintext      = bytes(request.plaintext, 'utf-8')
  des            = DES.new(KEY, DES.MODE_ECB)
  padded_text    = pad(plaintext)
  encrypted_text = des.encrypt(padded_text)
  response       = DESEncryptedResponse(
    id=hashlib.sha256(encrypted_text).hexdigest(), 
    ciphertext=str(encrypted_text)
  )
  return response

@encrypt_router.post("/decrypt/des")
async def dec_des(
  request: DESDecryptedRequest, 
  response_model=DESDecryptedResponse
):
  """Decrypts a ciphertext string denoted in the payload
  """
  des            = DES.new(KEY, DES.MODE_ECB)
  encrypted_text = ast.literal_eval(request.ciphertext)
  decrypted_text = des.decrypt(encrypted_text)
  response       = DESDecryptedResponse(
    id=hashlib.sha256(encrypted_text).hexdigest(), 
    plaintext=decrypted_text.decode('utf-8').strip()
  )
  return response