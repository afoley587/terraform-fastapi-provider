from pydantic import BaseModel

class DESEncryptedRequest(BaseModel):
    plaintext: str

class DESEncryptedResponse(BaseModel):
    id: str
    ciphertext: str

class DESDecryptedRequest(BaseModel):
  ciphertext: str

class DESDecryptedResponse(BaseModel):
  id: str
  plaintext: str