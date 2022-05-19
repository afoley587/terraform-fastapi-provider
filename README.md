## LinkedIn Post
So, as some of you may know, I'm currently getting my MS of Computer Science 
from Syracuse University. At this moment, I am taking quite a fun class:
Cryptography.

Last night, we were chatting about block algorithms and looking at
diagrams of Fiestel networks (gasp) and discussing DES and boy, 
I just really wanted to play around with it in python.

But, at that same time, I haven't built a REST API in a while and wanted
to shake off the rust.

But, at that same time, I wanted to build a terraform provider. I mean,
I use terraform hundreds of times a day, but NEVER made a provider? Who am I?

So, I combined all three! Take a look through my blog post here where I take you
through:

* Encrypting with DES in python (I KNOW DES IS OUTDATED, I JUST WANTED TO TOY WITH IT!)
* Build up a FastAPI REST API With Python
* Create a terraform provider to encrypt and decrypt strings for us

## Blog

So, you want to encrypt some data?
So, you want to build your own api?
So, you want to build your own terraform provider for your own api?
What are you? Addicted to DevOps?!

Today, I want to walk you through a fun full stack problem of mine:

  * Create a REST API With FastAPI
  * Play with encryption and decryption
  * Create a terraform provider FOR our API

How cool is that?

So, first some quick notes:

* I will be using DES in this example. YES I KNOW ITS OUTDATED! But, I'm currently getting
  a masters in CS and we learned about DES the other day. So, it seemed relevant and fun to
  me!
* The terraform provider could certainly be improved! There is no client library or
  provider config. I really just wanted to play with the semantics of the terraform
  provider and less with the go semantics. We are also storing both encrypted and 
  plaintext in our statefile, of course thats bad... but I just wanted to show the
  concepts.

So, lets get started!

## Python
So, in order to interact with an API, we need an API. So lets build one.
In this post, we are going to use FastAPI as our API of choice. 

### Layout
Let's first go over the directory layout of the project:
```
root_dir:
  python/
    app.py
    poetry.lock
    pyproject.toml
    run.sh
    routers/
      # routers define the routes of our API
      encryption.py
      __init__.py
    schemas/
      # schemas define the request and response bodies
      # we want to either validate or deliver
      encryption.py
      __init__.py
  golang/
    golang-related-files
  terraform/
    terraform-related-files
```

You'll notice I glazed over the terraform and golang files. They are not important
yet, so lets focus on the python stuff.

### Dependencies
If you're using
TOML format, you'll need these dependencies:
```
[tool.poetry.dependencies]
python = "^3.8"
requests = "^2.24.0"
fastapi = "0.78.0"
uvicorn = "0.17.6"
beautifulsoup4 = "4.11.1"
cachetools = "5.1.0"
pycrypto = "2.6.1"
```

If you're using a requirements.txt format, you'll need these dependecies:
```
```

### app.py
Lets first make our `app.py`, its going to be the entrypoint into our
API:
```
from fastapi import FastAPI

from routers.encryption import encrypt_router

def init_app():
  """Initializes the app object and adds the router to it
  """
  app = FastAPI()
  app.include_router(encrypt_router)
  return app

app = init_app()
```

So first, we import the required libraries. You'll see that we are including our
encryption router from the routers directory using relative imports. More on what those
routers are later!

We then call `init_app` to initialize a FastApi object. We then attach the `encrypt_router`
to our app. When we do that, all of the routes on the `encrypt_router` become available 
on our API.

And that's pretty much it!

### schemas
Schemas are also very simple. They let us nicely and cleanly define the
payloads we should expect from clients and also the payloads we want to deliver
to clients. Lets look at `python/schemas/encryption.py`:

```
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
```

Thats the whole file! So, lets break this down. The request schemas will be 
used as validation on our API endpoints (see routers section). If, for some
example, a client sends us an invalid body, FastAPI will automatically return
a 422 error. If a valid payload is received, it will automatically deserialize
it into the appropriate python object. Let's do an example:

```
## BAD ##
# client request
curl -H "Content-Type: application/json" -d'{ "garbage": "HAHA" }' localhost:8000/some/endpoint

# Server Response
422 Unprocessable Entity

## GOOD ##
# client request - Using DESEncryptedRequest Schema
curl -H "Content-Type: application/json" -d'{ "plaintext": "HAHA" }' localhost:8000/some/endpoint

# Server Response
200
```

FastAPI handles all of that under the hood, which is nice for us. Lets take a look at the
routers to get a better feel.

### routers
Routers are the heart and soul of our python project. They define which routes
correspond to which functions. All of our routes are in `python/routers/encryption.py`. Lets go 
through them! The top of our file is simple, it just does our imports, defines some helpers, and defines this new cool object called an `APIRouter`. This is the same router that we imported
in our `app.py` file above!

```
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
```

Now, we crack into our first router! We create a new route on our 
`encrypt_router` or type `POST` and the URL location `/encrypt/des`.
Notice, we say that our request is of type `DESEncryptedRequest` and 
it's response model will be of t type `DESEncryptedResponse`. This is
using our schemas defined above and shows the neat automated 
serialization that FastAPI supports! 

Then our function uses DES to encrypt the string in the payload, 
calculates its SHA-256 sum, and returns it to the client. We use a SHA-256
as an ID for two reasons:

* It shows that messages didn't change in transit
* It's a deterministic string

```
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
```

Our next router is almost the same thing, but in reverse. If we go step 
by step:

1. We have a new route at `/decrypt/des` of method `POST`
2. It's input body is of type `DESDecryptedRequest`
3. It's response body is of type `DESDecryptedResponse`
4. It pulls the ciphertext from the payload
5. It decrypts the ciphertext
6. It returns the plaintext and the ciphertext SHA-256 sum

```
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
```

### Running

Alas, our API is complete!
We can open one terminal and run it with:
```
#!/bin/bash

export DES_KEY='hello123'

poetry run uvicorn app:app --reload
```

And let's encrypt a string:
```
curl -X POST \
  -H "Content-Type: application/json" \
  -d'{ "plaintext": "foobar" }' \
  http://localhost:8000/encrypt/des
{"id":"a6f5d8295f261b1dbb8631b61dd757045122b25bb29364cd97301086ca5d2e84","ciphertext":"b'+\\x85v\"\\x04\\xfb_\\x9a'"}
```

## GoLang
So, our API is up. Now, lets write the terraform provider!

### main.go
First, we need to create the simple `main.go`.
This file will will serve as the entry point to our provider and follows
a boilerplate, from the most part. `main.go` is used to invoke our provider, specified
by `Provider`.
```
// main.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}
```

### provider.go
`provider.go` is the next file we need to build. This file
will define:

* Provider configurations:
  * usernames
  * passwords
  * api keys
  * etc.
* The resources offered
* The data sources offered

If we break down our file below, we are creating a provider
that offers two resources:
1. `garbage_des_encrypt` which corresponds to the `desEncrypt` object
2. `garbage_des_decrypt` which corresponds to the `desDecrypt` object
```
// provider.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"garbage_des_encrypt": desEncrypt(),
			"garbage_des_decrypt": desDecrypt(),
		},
	}
}
```

### resources
Now, comes the meat and potatoes. The resources! For brevity, 
I am only going to comb through the `des_encrypt.go` file. The
encryption and decryption resources are essentially identical, 
so I'll leave it up to you rip through the decryption.

First, our imports which we will skip
```
package main

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"net/http"
)
```

Then, our resource definition. We need to follow the terraform/hashicorp
schema for resources where the following are defined:

* Create
* Read
* Update
* Delete
* Schema

In our schema, you can see that we are defining two attributes:

* `plaintext` - which will be the plaintext we want to encrypt
* `ciphertext` - The computed, encrypted string from our API
  * Note that `Computed: true` for `ciphertext`, telling terraform that 
    this will be computed during a terraform run

```
func desEncrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceDesEncryptCreate,
		Read:   resourceDesEncryptRead,
		Update: resourceDesEncryptUpdate,
		Delete: resourceDesEncryptDelete,

		Schema: map[string]*schema.Schema{
			/* Our input to our resource */
			"plaintext": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			/* Our computed attribute of our resource */
			"ciphertext": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
```

We are also going to define two objects that we are going to
use to both marshall and unmarshall data that is going between
us and our server. They should look identical to the schemas defined
in `python/schemas/encryption`! (I hope it's all coming together!)


```
type DesEncryptRequestBody struct {
	Plaintext string `json:"plaintext"`
}

type DesEncryptResponse struct {
	Id         string `json:"id"`
	Ciphertext string `json:"ciphertext"`
}
```

Finally, we have to do the dang thing! Our `resourceDesEncryptCreate` is
going to run on `terraform apply` when a resource is created.

First, we are going to craft our payload to the server by:
1. Pulling the data of our the terraform configuration with `d.Get("plaintext").(string)`
2. Putting this in to the expected schema by creating an object of type `DesEncryptRequestBody`
3. Marshalling this to a JSON object with `json.Marshal(desRequestBody)`

```
func resourceDesEncryptCreate(d *schema.ResourceData, m interface{}) error {
	/* Create the GoLang Object with the ciphertext from our resource */
	log.Printf("[INFO] Starting call to http://127.0.0.1:8000/")
	desRequestBody := DesEncryptRequestBody{
		Plaintext: d.Get("plaintext").(string),
	}

	/* Convert DesEncryptRequestBody to byte using Json.Marshal
	 * Ignoring error.
	 */
	body, _ := json.Marshal(desRequestBody)

	/* Actually send the payload!
	 * Make sure to send the JSON with the proper Content-Type Header
	 * Set the body to our JSON bytes
	 */
```

With our body ready, we can then send that to the server with `http.Post`.
```
	resp, err := http.Post("http://127.0.0.1:8000/encrypt/des", "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Printf("[INFO] Could not initialize 'http.NewRequest' to 'http://localhost:8000/'.")
		return err
	}
```

Once the server responds, we can then:
1. Read the bytes from the server with `ioutil.ReadAll(resp.Body)`
2. Unmarshal the JSON from a string into a usable object of type `DesEncryptResponse` with `json.Unmarshal`

```
	/* Reading the response body. This will be the JSON
	 * representation of the DESEncryptedResponse that we defined
	 * in the python/schemas/encryption.py file!
	 */
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("[INFO] Could not read response body.")
		return err
	}

	/* Deserialize the JSON/Byte response to our expected response */
	var respJsonDecoded DesEncryptResponse
	err = json.Unmarshal([]byte(string(respBody)), &respJsonDecoded)

	if err != nil {
		log.Printf("[INFO] Could not unmarshal JSON response.")
		return err
	}
```

And, alas, we can finally save or terraform state! We do that with the `d.SetId`
and `d.Set` functions. After this, we can look in our state file and see that
our attribute has a SHA-256 sum as its ID and some encrypted string
as its `ciphertext` attribute!
```
	/* Save our outputs, id and plaintext
	 */
	d.SetId(respJsonDecoded.Id)
	d.Set("ciphertext", respJsonDecoded.Ciphertext)
	/* Best practice is just to read it from the API.
	 * of course, this does NOTHING in our case, but
	 * acts as a good sanity check
	 */
	resourceDesEncryptRead(d, m)
	return nil
}
```

The rest of the functions are created to either read, update, or delete the
resource. Other use cases will most likely fill these out, but our server
is stateless, so it might not make much sense in our use case.
```
func resourceDesEncryptRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDesEncryptUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceDesEncryptRead(d, m)
}

func resourceDesEncryptDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
```
### Building
Well, the hard parts are all done. Now we can build our code into a terraform usable plugin:
```
cd golang
go mod init 'terraform-example'
go fmt
go build -o terraform-provider-garbage
# Note for Mac Users, I had to use 
# export CGO_CPPFLAGS="-Wno-error -Wno-nullability-completeness -Wno-expansion-to-defined -Wbuiltin-requires-header"
# go build -o terraform-provider-garbage
```

Once its built, you will have to copy it to the following location: PLUGIN_DIR/HOSTNAME/PROVIDER/PROVIDER/VERSION/PLATFORM/`
Where:

* PLUGIN_DIR is the terraform plugins dir, typically at `$HOME/.terraform.d/plugins`
* HOSTNAME is the terraform registry hostname, such as `terraform-example.com`
* PROVIDER is the provider name, such as `garbage`
* VERSION is the version number, such as `0.0.1`
* PLATFORM is the platform of your machine, such as `darwin_amd64`

So, you can run 
```
mkdir -p ${PLUGIN_DIR}/${HOSTNAME}/${PROVIDER}/${PROVIDER}/${VERSION}/${PLATFORM}/
cp terraform-provider-garbage ${PLUGIN_DIR}/${HOSTNAME}/${PROVIDER}/${PROVIDER}/${VERSION}/${PLATFORM}/
```
And you're good to go! A Makefile is included for your ease of use in the Git Repo.

## Terraform
Now is the easy part - running it. Lets create a new terraform file and get rocking!

### main.tf
At the top of the `main.tf`, let's set use our new provider:

```
terraform {
  required_providers {
    garbage = {
      version = "~> 0.0.1"
      source  = "terraform-example.com/garbage/garbage"
    }
  }
}
```

Let's add a few resources:
```
resource "garbage_des_encrypt" "des_encrypt" {
  plaintext = "test"
}

resource "garbage_des_decrypt" "des_decrypt" {
  ciphertext = garbage_des_encrypt.des_encrypt.ciphertext
}
```

Notice that we are creating some ciphertext from plaintext and then doing the reverse operation.
If all goes well, both should have the same:
* plaintext
* ciphertext
* id

Finally, we show our outputs:
```
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
```

### Running
First, make sure your python API server is running, and then you can run 
`terraform init` and `terraform apply`!
```
prompt> terraform init

Initializing the backend...

Initializing provider plugins...
- Finding terraform-example.com/garbage/garbage versions matching "~> 0.0.1"...
- Installing terraform-example.com/garbage/garbage v0.0.1...
- Installed terraform-example.com/garbage/garbage v0.0.1 (unauthenticated)

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

prompt> terraform apply -auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # garbage_des_decrypt.des_decrypt will be created
  + resource "garbage_des_decrypt" "des_decrypt" {
      + ciphertext = (known after apply)
      + id         = (known after apply)
      + plaintext  = (known after apply)
    }

  # garbage_des_encrypt.des_encrypt will be created
  + resource "garbage_des_encrypt" "des_encrypt" {
      + ciphertext = (known after apply)
      + id         = (known after apply)
      + plaintext  = "test"
    }

Plan: 2 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + ciphertext           = (known after apply)
  + ciphertext_sum       = (known after apply)
  + did_properly_encrypt = (known after apply)
  + plaintext            = (known after apply)
  + plaintext_sum        = (known after apply)
  + verified_sums        = (known after apply)
garbage_des_encrypt.des_encrypt: Creating...
garbage_des_encrypt.des_encrypt: Creation complete after 0s [id=3cd9d7aefaa1f16c5333f804c725f5235f724aee4fa59f16d2a14600479b2a84]
garbage_des_decrypt.des_decrypt: Creating...
garbage_des_decrypt.des_decrypt: Creation complete after 0s [id=3cd9d7aefaa1f16c5333f804c725f5235f724aee4fa59f16d2a14600479b2a84]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

ciphertext = "b'\\xe5\\x92BO\\x1f\\x97\\xbd2'"
ciphertext_sum = "3cd9d7aefaa1f16c5333f804c725f5235f724aee4fa59f16d2a14600479b2a84"
did_properly_encrypt = true
plaintext = "test"
plaintext_sum = "3cd9d7aefaa1f16c5333f804c725f5235f724aee4fa59f16d2a14600479b2a84"
verified_sums = true
```