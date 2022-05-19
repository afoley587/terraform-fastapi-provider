terraform {
  required_providers {
    garbage = {
      version = "~> 0.0.1"
      source  = "terraform-example.com/garbage/garbage"
    }
  }
}

resource "garbage_des_encrypt" "des_encrypt" {
  plaintext = "test"
}

resource "garbage_des_decrypt" "des_decrypt" {
  ciphertext = garbage_des_encrypt.des_encrypt.ciphertext
}

output "ciphertext" {
  value = garbage_des_encrypt.des_encrypt.ciphertext
}

output "ciphertext_sum" {
  value = garbage_des_encrypt.des_encrypt.id
}

output "plaintext" {
  value = garbage_des_decrypt.des_decrypt.plaintext
}

output "plaintext_sum" {
  value = garbage_des_decrypt.des_decrypt.id
}

output "did_properly_encrypt" {
  value = (
    garbage_des_decrypt.des_decrypt.plaintext == garbage_des_encrypt.des_encrypt.plaintext
  )
} 

output "verified_sums" {
  value = (
    garbage_des_decrypt.des_decrypt.id == garbage_des_encrypt.des_encrypt.id
  )
} 